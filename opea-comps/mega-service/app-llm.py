import os

from comps import (
    MegaServiceEndpoint,
    MicroService,
    ServiceOrchestrator,
    ServiceRoleType,
    ServiceType,
)
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

MEGA_SERVICE_PORT = os.getenv("MEGA_SERVICE_PORT", 8888)

# EMBEDDING_SERVICE_HOST_IP = os.getenv("EMBEDDING_SERVICE_HOST_IP", "0.0.0.0")
# EMBEDDING_SERVICE_PORT = os.getenv("EMBEDDING_SERVICE_PORT", 6000)
LLM_SERVICE_HOST_IP = os.getenv("LLM_SERVICE_HOST_IP", "0.0.0.0")
LLM_SERVICE_PORT = os.getenv("LLM_SERVICE_PORT", 9000)  ## LLM is running on 8008


class ChatWithPdfService:
    def __init__(self, host="0.0.0.0", port=8000):
        self.host = host
        self.port = port
        self.megaservice = ServiceOrchestrator()
        self.endpoint = str(MegaServiceEndpoint.CHAT_QNA)

    def add_remote_service(self):
        # embedding = MicroService(
        #     name="embedding",
        #     host=EMBEDDING_SERVICE_HOST_IP,
        #     port=EMBEDDING_SERVICE_PORT,
        #     endpoint="/v1/embeddings",
        #     use_remote_service=True,
        #     service_type=ServiceType.EMBEDDING,
        # )
        llm = MicroService(
            name="llm",
            host=LLM_SERVICE_HOST_IP,
            port=LLM_SERVICE_PORT,
            endpoint="/v1/chat/completions",
            use_remote_service=True,
            service_type=ServiceType.LLM,
        )
        self.megaservice.add(llm)

    async def handle_request(self, request: Request):
        try:
            data = await request.json()

            # Transform messages if it's a string
            if isinstance(data.get("messages"), str):
                data["messages"] = [{"role": "user", "content": data["messages"]}]

            stream_opt = data.get("stream", True)
            chat_request = ChatCompletionRequest.model_validate(data)
            prompt = handle_message(chat_request.messages) # Not used for now

            parameters = LLMParams(
                model="llama3", # This parameter is needed for LLM service  to know which model to use.
                max_tokens=chat_request.max_tokens if chat_request.max_tokens else 1024,
                top_k=chat_request.top_k if chat_request.top_k else 10,
                top_p=chat_request.top_p if chat_request.top_p else 0.95,
                temperature=chat_request.temperature
                if chat_request.temperature
                else 0.01,
                frequency_penalty=chat_request.frequency_penalty
                if chat_request.frequency_penalty
                else 0.0,
                presence_penalty=chat_request.presence_penalty
                if chat_request.presence_penalty
                else 0.0,
                repetition_penalty=chat_request.repetition_penalty
                if chat_request.repetition_penalty
                else 1.03,
                stream=stream_opt,
                chat_template=chat_request.chat_template
                if chat_request.chat_template
                else None,
            )

            result_dict, runtime_graph = await self.megaservice.schedule(
                initial_inputs=data,
                llm_parameters=parameters,
            )

            # Handle streaming response
            for node, response in result_dict.items():
                if isinstance(response, StreamingResponse):
                    return response

            # Handle non-streaming response
            last_node = runtime_graph.all_leaves()[-1]
            if last_node not in result_dict:
                raise KeyError(f"Last node {last_node} not found in result_dict")

            response_data = result_dict[last_node]
            if isinstance(response_data, dict):
                response_text = (
                    response_data.get("text")
                    or response_data.get("content")
                    or str(response_data)
                )
            else:
                response_text = str(response_data)

            choices = []
            usage = UsageInfo()
            choices.append(
                ChatCompletionResponseChoice(
                    index=0,
                    message=ChatMessage(role="assistant", content=response_text),
                    finish_reason="stop",
                )
            )
            return ChatCompletionResponse(model="chatqna", choices=choices, usage=usage)

        except Exception as e:
            import traceback

            print(f"Error processing request: {str(e)}")
            print(f"Traceback: {traceback.format_exc()}")
            raise

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


if __name__ == "__main__":
    chatqna = ChatWithPdfService(port=MEGA_SERVICE_PORT)
    chatqna.add_remote_service()
    chatqna.start()
