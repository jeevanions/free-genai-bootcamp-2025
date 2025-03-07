import json
import os
from dotenv import load_dotenv
import requests  # Add this import for direct debugging calls

from comps import MegaServiceEndpoint, MicroService, ServiceOrchestrator, ServiceRoleType, ServiceType
from comps.cores.mega.utils import handle_message
from comps.cores.proto.api_protocol import (
    ChatCompletionRequest,
    ChatCompletionResponse,
    ChatCompletionResponseChoice,
    ChatMessage,
    UsageInfo,
)
from comps.cores.proto.docarray import LLMParams, RerankerParms, RetrieverParms
from fastapi import Request
from fastapi.responses import StreamingResponse
from langchain_core.prompts import PromptTemplate


class ChatTemplate:
    @staticmethod
    def generate_rag_prompt(question, documents):
        context_str = "\n".join(documents)
        template = """
### You are a helpful, respectful and honest assistant to help the user with questions. \
Please refer to the search results obtained from the local knowledge base. \
But be careful to not incorporate the information that you think is not relevant to the question. \
If you don't know the answer to a question, please don't share false information. \n
### Search results: {context} \n
### Question: {question} \n
### Answer:
"""
        return template.format(context=context_str, question=question)

load_dotenv(".env")


MEGA_SERVICE_PORT = int(os.getenv("MEGA_SERVICE_PORT", 8888))
GUARDRAIL_SERVICE_HOST_IP = os.getenv("GUARDRAIL_SERVICE_HOST_IP", "0.0.0.0")
GUARDRAIL_SERVICE_PORT = int(os.getenv("GUARDRAIL_SERVICE_PORT", 80))

EMBEDDING_SERVER_HOST_IP = os.getenv("EMBEDDING_SERVICE_HOST_IP", "0.0.0.0")
EMBEDDING_SERVER_PORT = int(os.getenv("EMBEDDING_SERVICE_PORT", 8007))

RETRIEVER_SERVICE_HOST_IP = os.getenv("RETRIEVER_SERVICE_HOST_IP", "0.0.0.0")
RETRIEVER_SERVICE_PORT = int(os.getenv("RETRIEVER_SERVICE_PORT", 8006))

RERANK_SERVER_HOST_IP = os.getenv("RERANKER_SERVICE_HOST_IP", "0.0.0.0")
RERANK_SERVER_PORT = int(os.getenv("RERANKER_SERVICE_PORT", 8005))

LLM_SERVER_HOST_IP = os.getenv("LLM_SERVICE_HOST_IP", "0.0.0.0")
LLM_SERVER_PORT = int(os.getenv("LLM_SERVICE_PORT", 8008))

LLM_MODEL = os.getenv("LLM_MODEL_ID", "llama3")

QDRANT_HOST = os.getenv("QDRANT_HOST", "localhost")
QDRANT_PORT = int(os.getenv("QDRANT_PORT", 6333))
QDRANT_EMBED_DIMENSION = os.getenv("QDRANT_EMBED_DIMENSION", 768)
QDRANT_INDEX_NAME = os.getenv("QDRANT_INDEX_NAME", "rag_qdrant")
print(f"QDRANT_HOST: {QDRANT_HOST}, QDRANT_PORT: {QDRANT_PORT}, QDRANT_EMBED_DIMENSION: {QDRANT_EMBED_DIMENSION}, QDRANT_INDEX_NAME: {QDRANT_INDEX_NAME}")


def align_inputs(self, inputs, cur_node, runtime_graph, llm_parameters_dict, **kwargs):
    if self.services[cur_node].service_type == ServiceType.EMBEDDING:
        inputs["input"] = inputs["text"]
        del inputs["text"]
    elif self.services[cur_node].service_type == ServiceType.RETRIEVER:
        # prepare the retriever params
        retriever_parameters = kwargs.get("retriever_parameters", None)
        if (retriever_parameters):
            inputs.update(retriever_parameters.model_dump())
    elif self.services[cur_node].service_type == ServiceType.LLM:
        # Debug the inputs coming in
        print(f"LLM inputs before transformation: {inputs}")
        
        # convert TGI/vLLM to unified OpenAI /v1/chat/completions format
        next_inputs = {}
        next_inputs["model"] = LLM_MODEL
        next_inputs["messages"] = [{"role": "user", "content": inputs["inputs"]}]
        next_inputs["max_tokens"] = llm_parameters_dict["max_tokens"]
        next_inputs["top_p"] = llm_parameters_dict["top_p"]
        next_inputs["stream"] = inputs["stream"]
        
        # Only include parameters that are actually supported by the LLM API
        if "frequency_penalty" in inputs:
            next_inputs["frequency_penalty"] = inputs["frequency_penalty"]
        if "temperature" in inputs:
            next_inputs["temperature"] = inputs["temperature"]
        
        # Log the final request being sent
        print(f"LLM request being sent: {json.dumps(next_inputs, indent=2)}")
        
        # For debugging - try a direct request to verify the LLM service is working
        # try:
        #     direct_url = f"http://{LLM_SERVER_HOST_IP}:{LLM_SERVER_PORT}/v1/chat/completions"
        #     print(f"Making direct test request to: {direct_url}")
        #     direct_response = requests.post(
        #         direct_url,
        #         json=next_inputs,
        #         timeout=5
        #     )
        #     print(f"Direct LLM test response status: {direct_response.status_code}")
        #     if direct_response.status_code == 200:
        #         print(f"Direct response sample: {direct_response.text[:100]}...")
        #     else:
        #         print(f"Error response: {direct_response.text}")
        # except Exception as e:
        #     print(f"Direct test request failed: {str(e)}")
        
        inputs = next_inputs
    return inputs


def align_outputs(self, data, cur_node, inputs, runtime_graph, llm_parameters_dict, **kwargs):
    next_data = {}
    if self.services[cur_node].service_type == ServiceType.EMBEDDING:
        # print(f"Embedding service response: {data}")
        if isinstance(data, dict) and "error" in data:
            # Handle error case - pass the error downstream
            print(f"Embedding service error: {data['error']}")
            # Create a dummy embedding to continue the flow
            next_data = {"text": inputs["input"], "embedding": [0.0] * 1024, "error": data["error"]}
        else:
            # Handle normal case
            assert isinstance(data.get("data"), list), f"Expected list from embedding service, got {type(data)}: {data}"
            next_data = {"text": inputs["input"], "embedding": data.get("data")[0].get("embedding")}
    elif self.services[cur_node].service_type == ServiceType.RETRIEVER:
        # Print detailed information about retriever response structure for debugging
        print(f"Retriever response type: {type(data)}")
        print(f"Retriever response structure: {dir(data) if hasattr(data, '__dir__') else 'No dir method'}")
        print(f"Retriever response representation: {repr(data)}")
        
        # Safely check if we got an error response
        if isinstance(data, dict) and "error" in data:
            print(f"Retriever service error: {data['error']}")
            next_data["inputs"] = f"The system encountered an error retrieving relevant information. Here's the original question: {inputs.get('text', 'No question provided')}"
            
            # Skip reranker if it's in the pipeline
            if len(runtime_graph.downstream(cur_node)) > 0 and runtime_graph.downstream(cur_node)[0].startswith("rerank"):
                for ds in runtime_graph.downstream(cur_node):
                    for nds in runtime_graph.downstream(ds):
                        runtime_graph.add_edge(cur_node, nds)
                    runtime_graph.delete_node_if_exists(ds)
                
            return next_data
        
        # Safely extract retrieved docs
        retrieved_docs = []
        initial_query = ""
        
        # Handle different response structures
        if hasattr(data, "retrieved_docs"):
            retrieved_docs = data.retrieved_docs
            initial_query = data.initial_query if hasattr(data, "initial_query") else inputs.get("text", "")
        elif isinstance(data, dict):
            retrieved_docs = data.get("retrieved_docs", [])
            initial_query = data.get("initial_query", inputs.get("text", ""))
        
        print(f"Retrieved docs count: {len(retrieved_docs)}")
        
        # Extract document text safely
        docs = []
        for doc in retrieved_docs:
            if isinstance(doc, str):
                docs.append(doc)
            elif isinstance(doc, dict) and "text" in doc:
                docs.append(doc["text"])
            elif hasattr(doc, "text"):
                docs.append(doc.text)
            elif hasattr(doc, "page_content"):
                docs.append(doc.page_content)
        
        print(f"Extracted document texts: {len(docs)}")
        
        # Handle empty results
        if not docs:
            print("No documents returned from retriever")
            next_data["inputs"] = f"I don't have specific information about that in my knowledge base. Here's what I know about: {inputs.get('text', 'your question')}"
            
            # Skip reranker if it's in the pipeline
            if len(runtime_graph.downstream(cur_node)) > 0 and runtime_graph.downstream(cur_node)[0].startswith("rerank"):
                for ds in runtime_graph.downstream(cur_node):
                    for nds in runtime_graph.downstream(ds):
                        runtime_graph.add_edge(cur_node, nds)
                    runtime_graph.delete_node_if_exists(ds)
                
            return next_data
        
        with_rerank = False
        if len(runtime_graph.downstream(cur_node)) > 0:
            with_rerank = runtime_graph.downstream(cur_node)[0].startswith("rerank")
            
        if with_rerank and docs:
            # Prepare for reranking
            next_data["query"] = initial_query
            next_data["texts"] = docs
        else:
            # Skip reranking, prepare prompt for LLM
            # Handle template generation with docs
            prompt = initial_query
            chat_template = llm_parameters_dict["chat_template"]
            if chat_template:
                prompt_template = PromptTemplate.from_template(chat_template)
                input_variables = prompt_template.input_variables
                if sorted(input_variables) == ["context", "question"]:
                    prompt = prompt_template.format(question=initial_query, context="\n".join(docs))
                elif input_variables == ["question"]:
                    prompt = prompt_template.format(question=initial_query)
                else:
                    print(f"{prompt_template} not used, we only support 2 input variables ['question', 'context']")
                    prompt = ChatTemplate.generate_rag_prompt(initial_query, docs)
            else:
                prompt = ChatTemplate.generate_rag_prompt(initial_query, docs)

            next_data["inputs"] = prompt

    elif self.services[cur_node].service_type == ServiceType.RERANK:
        # rerank the inputs with the scores
        reranker_parameters = kwargs.get("reranker_parameters", None)
        top_n = reranker_parameters.top_n if reranker_parameters else 1
        docs = inputs["texts"]
        reranked_docs = []
        for best_response in data[:top_n]:
            reranked_docs.append(docs[best_response["index"]])

        # handle template
        # if user provides template, then format the prompt with it
        # otherwise, use the default template
        prompt = inputs["query"]
        chat_template = llm_parameters_dict["chat_template"]
        if chat_template:
            prompt_template = PromptTemplate.from_template(chat_template)
            input_variables = prompt_template.input_variables
            if sorted(input_variables) == ["context", "question"]:
                prompt = prompt_template.format(question=prompt, context="\n".join(reranked_docs))
            elif input_variables == ["question"]:
                prompt = prompt_template.format(question=prompt)
            else:
                print(f"{prompt_template} not used, we only support 2 input variables ['question', 'context']")
                prompt = ChatTemplate.generate_rag_prompt(prompt, reranked_docs)
        else:
            prompt = ChatTemplate.generate_rag_prompt(prompt, reranked_docs)

        next_data["inputs"] = prompt

    elif self.services[cur_node].service_type == ServiceType.LLM and not llm_parameters_dict["stream"]:
        next_data["text"] = data["choices"][0]["message"]["content"]
    else:
        next_data = data

    return next_data


def align_generator(self, gen, **kwargs):
    print("Starting align_generator with stream response")
    for line in gen:
        line = line.decode("utf-8")
        print(f"Raw stream line: {line}")
        
        start = line.find("{")
        end = line.rfind("}") + 1

        if start == -1 or end == 0:
            print("No valid JSON found in line, skipping")
            continue

        json_str = line[start:end]
        try:
            # sometimes yield empty chunk, do a fallback here
            json_data = json.loads(json_str)
            if (
                json_data["choices"][0]["finish_reason"] != "eos_token"
                and "content" in json_data["choices"][0]["delta"]
            ):
                content = json_data["choices"][0]["delta"]["content"]
                print(f"Yielding content: {content}")
                # Format as proper SSE with JSON payload
                response = json.dumps({"content": content})
                yield f"data: {response}\n\n"
        except Exception as e:
            print(f"Error processing JSON: {str(e)}")
            # Properly format as SSE
            yield f"data: {json.dumps({'content': json_str})}\n\n"
    print("Stream complete, sending DONE marker")
    yield "data: [DONE]\n\n"


class ChatWithPdfService:
    def __init__(self, host="0.0.0.0", port=8000):
        self.host = host
        self.port = port
        ServiceOrchestrator.align_inputs = align_inputs
        ServiceOrchestrator.align_outputs = align_outputs
        ServiceOrchestrator.align_generator = align_generator
        self.megaservice = ServiceOrchestrator()
        self.endpoint = str(MegaServiceEndpoint.CHAT_QNA)

    def add_remote_service(self):

        embedding = MicroService(
            name="embedding",
            host=EMBEDDING_SERVER_HOST_IP,
            port=EMBEDDING_SERVER_PORT,
            endpoint="/v1/embeddings",
            use_remote_service=True,
            service_type=ServiceType.EMBEDDING,
        )

        retriever = MicroService(
            name="retriever",
            host=RETRIEVER_SERVICE_HOST_IP,
            port=RETRIEVER_SERVICE_PORT,
            endpoint="/v1/retrieval",
            use_remote_service=True,
            service_type=ServiceType.RETRIEVER,
        )

        rerank = MicroService(
            name="rerank",
            host=RERANK_SERVER_HOST_IP,
            port=RERANK_SERVER_PORT,
            endpoint="/v1/reranking",
            use_remote_service=True,
            service_type=ServiceType.RERANK,
        )

        # Print the LLM service configuration for debugging
        print(f"LLM service configured at: http://{LLM_SERVER_HOST_IP}:{LLM_SERVER_PORT}/v1/chat/completions")
        
        llm = MicroService(
            name="llm",
            host=LLM_SERVER_HOST_IP,
            port=LLM_SERVER_PORT,
            endpoint="/v1/chat/completions",
            use_remote_service=True,
            service_type=ServiceType.LLM,
        )
        self.megaservice.add(embedding).add(retriever).add(rerank).add(llm)
        self.megaservice.flow_to(embedding, retriever)
        self.megaservice.flow_to(retriever, rerank)
        self.megaservice.flow_to(rerank, llm)

    async def handle_request(self, request: Request):
        data = await request.json()
        print(f"Received request: {data}")
        
        # Set stream explicitly to True
        stream_opt = data.get("stream", True)
        
        # Fix message format if it's not properly structured
        if "messages" in data and isinstance(data["messages"], str):
            print("Converting string message to proper format")
            data["messages"] = [{"role": "user", "content": data["messages"]}]
        
        # Ensure data is valid before proceeding
        if not data.get("messages"):
            return {"error": "No messages provided in request"}
        
        chat_request = ChatCompletionRequest.model_validate(data)
        prompt = handle_message(chat_request.messages)
        print(f"Processed prompt: {prompt}")
        
        # Ensure streaming is explicitly enabled
        parameters = LLMParams(
            max_tokens=chat_request.max_tokens if chat_request.max_tokens else 1024,
            top_k=chat_request.top_k if chat_request.top_k else 10,
            top_p=chat_request.top_p if chat_request.top_p else 0.95,
            temperature=chat_request.temperature if chat_request.temperature else 0.01,
            frequency_penalty=chat_request.frequency_penalty if chat_request.frequency_penalty else 0.0,
            presence_penalty=chat_request.presence_penalty if chat_request.presence_penalty else 0.0,
            repetition_penalty=chat_request.repetition_penalty if chat_request.repetition_penalty else 1.03,
            stream=True,  # Force stream to True for testing
            chat_template=chat_request.chat_template if chat_request.chat_template else None,
        )
        retriever_parameters = RetrieverParms(
            search_type=chat_request.search_type if chat_request.search_type else "similarity",
            k=chat_request.k if chat_request.k else 4,
            distance_threshold=chat_request.distance_threshold if chat_request.distance_threshold else None,
            fetch_k=chat_request.fetch_k if chat_request.fetch_k else 20,
            lambda_mult=chat_request.lambda_mult if chat_request.lambda_mult else 0.5,
            score_threshold=chat_request.score_threshold if chat_request.score_threshold else 0.2,
        )
        reranker_parameters = RerankerParms(
            top_n=chat_request.top_n if chat_request.top_n else 1,
        )
        result_dict, runtime_graph = await self.megaservice.schedule(
            initial_inputs={"text": prompt, "stream": True},  # Explicitly set stream here as well
            llm_parameters=parameters,
            retriever_parameters=retriever_parameters,
            reranker_parameters=reranker_parameters,
        )
        
        for node, response in result_dict.items():
            print(f"Node {node} response type: {type(response)}")
            if isinstance(response, StreamingResponse):
                print("Returning StreamingResponse")
                return response
        last_node = runtime_graph.all_leaves()[-1]
        response = result_dict[last_node]["text"]
        choices = []
        usage = UsageInfo()
        choices.append(
            ChatCompletionResponseChoice(
                index=0,
                message=ChatMessage(role="assistant", content=response),
                finish_reason="stop",
            )
        )
        return ChatCompletionResponse(model="chatqna", choices=choices, usage=usage)

    def start(self):

        self.service = MicroService(
            self.__class__.__name__,
            service_role=ServiceRoleType.MEGASERVICE,
            host=self.host,
            port=self.port,
            endpoint=self.endpoint,
            input_datatype=ChatCompletionRequest,
            output_datatype=ChatCompletionResponse,
        )

        self.service.add_route(self.endpoint, self.handle_request, methods=["POST"])

        self.service.start()


# if __name__ == "__main__":

# load_dotenv(".env")

chat_with_pdf = ChatWithPdfService(host=os.getenv("HOST_IP"), port=MEGA_SERVICE_PORT)
chat_with_pdf.add_remote_service()

chat_with_pdf.start()
