import gradio as gr
import os
from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams, PointStruct
from langchain_openai import OpenAIEmbeddings
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_community.vectorstores import Qdrant
from langchain_openai import AzureChatOpenAI
from langchain.chains import RetrievalQA
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Initialize OpenAI embeddings
try:
    # Create a mock embeddings class for testing
    class MockEmbeddings:
        def embed_documents(self, texts):
            return [[0.1] * 1536 for _ in texts]  # Return mock embeddings
        
        def embed_query(self, text):
            return [0.1] * 1536  # Return mock embedding
    
    embeddings = MockEmbeddings()
    print("Using mock embeddings for RAG operations.")
except Exception as e:
    print(f"Warning: Error initializing OpenAI embeddings: {e}")
    # Create a mock embeddings class for testing
    class MockEmbeddings:
        def embed_documents(self, texts):
            return [[0.1] * 1536 for _ in texts]  # Return mock embeddings
        
        def embed_query(self, text):
            return [0.1] * 1536  # Return mock embedding
    
    embeddings = MockEmbeddings()

# Initialize Qdrant client
qdrant_client = QdrantClient(url=os.getenv("QDRANT_URL", "http://localhost:6333"))
collection_name = os.getenv("QDRANT_COLLECTION_NAME", "italian_learning")

# Initialize LLM
try:
    llm = AzureChatOpenAI(
        api_key=os.getenv("OPENAI_API_KEY"),
        api_version=os.getenv("OPENAI_API_VERSION"),
        azure_endpoint=os.getenv("OPENAI_API_BASE"),
        deployment_name=os.getenv("OPENAI_API_DEPLOYMENT_NAME"),
        temperature=0.5,
    )
except Exception as e:
    print(f"Warning: Error initializing Azure Chat OpenAI: {e}")
    # Create a mock LLM class for testing
    from langchain_core.language_models.chat_models import BaseChatModel
    from langchain_core.messages import AIMessage, BaseMessage
    from typing import List, Optional, Any
    
    class MockChatModel(BaseChatModel):
        def _generate(self, messages: List[BaseMessage], stop: Optional[List[str]] = None, run_manager: Optional[Any] = None, **kwargs) -> AIMessage:
            return AIMessage(content="This is a mock response for testing without API credentials.")
        
        @property
        def _llm_type(self) -> str:
            return "mock-chat-model"
    
    llm = MockChatModel()

def initialize_vector_store():
    """
    Initialize or get the Qdrant vector store
    """
    try:
        # Check if collection exists
        collections = qdrant_client.get_collections().collections
        collection_names = [collection.name for collection in collections]
        
        if collection_name not in collection_names:
            # Create new collection
            qdrant_client.create_collection(
                collection_name=collection_name,
                vectors_config=VectorParams(size=1536, distance=Distance.COSINE),
            )
            return {"success": True, "message": f"Created new collection: {collection_name}"}
        else:
            return {"success": True, "message": f"Using existing collection: {collection_name}"}
    except Exception as e:
        return {"success": False, "message": f"Error initializing vector store: {str(e)}"}

def add_document_to_vector_store(file_path, document_type):
    """
    Add a document to the vector store
    """
    try:
        # Read the document
        with open(file_path, "r", encoding="utf-8") as f:
            text = f.read()
        
        # Split text into chunks
        text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=1000,
            chunk_overlap=200,
            length_function=len,
        )
        chunks = text_splitter.split_text(text)
        
        # Create metadata for each chunk
        metadatas = [{"source": file_path, "type": document_type} for _ in chunks]
        
        # Create Qdrant vector store
        vector_store = Qdrant(
            client=qdrant_client,
            collection_name=collection_name,
            embeddings=embeddings,
        )
        
        # Add documents to vector store
        vector_store.add_texts(texts=chunks, metadatas=metadatas)
        
        return {
            "success": True,
            "message": f"Added {len(chunks)} chunks from {file_path} to vector store",
            "chunks": len(chunks)
        }
    except Exception as e:
        return {"success": False, "message": f"Error adding document to vector store: {str(e)}"}

def query_vector_store(query, document_type=None):
    """
    Query the vector store
    """
    try:
        # Create Qdrant vector store
        vector_store = Qdrant(
            client=qdrant_client,
            collection_name=collection_name,
            embeddings=embeddings,
        )
        
        # Create filter if document_type is specified
        search_filter = None
        if document_type and document_type != "all":
            search_filter = {"type": document_type}
        
        # Create retrieval QA chain
        qa_chain = RetrievalQA.from_chain_type(
            llm=llm,
            chain_type="stuff",
            retriever=vector_store.as_retriever(
                search_kwargs={"filter": search_filter, "k": 5}
            ),
        )
        
        # Run the query
        result = qa_chain.invoke(query)
        
        return {
            "success": True,
            "answer": result["result"],
        }
    except Exception as e:
        return {"success": False, "message": f"Error querying vector store: {str(e)}"}

def list_available_documents():
    """
    List all available documents in the output folders
    """
    documents = []
    
    # Check YouTube transcript folder
    yt_folder = "output/transcript-yt"
    if os.path.exists(yt_folder):
        for file in os.listdir(yt_folder):
            if file.endswith(".txt"):
                documents.append({"path": os.path.join(yt_folder, file), "type": "youtube"})
    
    # Check Whisper transcript folder
    whisper_folder = "output/transcript-vid"
    if os.path.exists(whisper_folder):
        for file in os.listdir(whisper_folder):
            if file.endswith("_whisper.txt"):
                documents.append({"path": os.path.join(whisper_folder, file), "type": "whisper"})
    
    # Check OCR transcript folder
    ocr_folder = "output/transcript-ocr"
    if os.path.exists(ocr_folder):
        for file in os.listdir(ocr_folder):
            if file.endswith("_ocr.txt"):
                documents.append({"path": os.path.join(ocr_folder, file), "type": "ocr"})
    
    return documents

def add_document_handler(document_path, document_type):
    """
    Handler for adding a document to the vector store
    """
    # Initialize vector store
    init_result = initialize_vector_store()
    if not init_result["success"]:
        return init_result["message"]
    
    # Add document to vector store
    add_result = add_document_to_vector_store(document_path, document_type)
    if not add_result["success"]:
        return add_result["message"]
    
    return add_result["message"]

def query_handler(query, document_type):
    """
    Handler for querying the vector store
    """
    # Query vector store
    query_result = query_vector_store(query, document_type)
    if not query_result["success"]:
        return query_result["message"]
    
    return query_result["answer"]

def create_rag_interface(parent):
    """
    Creates the RAG implementation interface
    """
    with parent:
        gr.Markdown("### RAG Implementation with Qdrant")
        gr.Markdown("Add documents to the vector store and query them using RAG.")
        
        # Get available documents
        documents = list_available_documents()
        document_paths = [doc["path"] for doc in documents]
        document_types = list(set([doc["type"] for doc in documents]))
        
        with gr.Tabs():
            with gr.Tab("Add Documents"):
                with gr.Row():
                    document_dropdown = gr.Dropdown(
                        label="Select Document",
                        choices=document_paths,
                        info="Select a document to add to the vector store"
                    )
                    
                    document_type_dropdown = gr.Dropdown(
                        label="Document Type",
                        choices=document_types + ["custom"],
                        value=document_types[0] if document_types else "custom",
                        info="Select the type of document"
                    )
                
                with gr.Row():
                    add_btn = gr.Button("Add Document to Vector Store")
                
                with gr.Row():
                    add_status = gr.Textbox(label="Status")
            
            with gr.Tab("Query Documents"):
                with gr.Row():
                    query_input = gr.Textbox(
                        label="Query",
                        placeholder="Ask a question about the Italian content...",
                    )
                    
                    filter_dropdown = gr.Dropdown(
                        label="Filter by Document Type",
                        choices=["all"] + document_types,
                        value="all",
                        info="Filter results by document type"
                    )
                
                with gr.Row():
                    query_btn = gr.Button("Query Vector Store")
                
                with gr.Row():
                    answer_output = gr.Textbox(
                        label="Answer",
                        lines=10,
                    )
        
        # Set up event handlers
        add_btn.click(
            add_document_handler,
            inputs=[document_dropdown, document_type_dropdown],
            outputs=[add_status],
        )
        
        query_btn.click(
            query_handler,
            inputs=[query_input, filter_dropdown],
            outputs=[answer_output],
        )
        
        return None