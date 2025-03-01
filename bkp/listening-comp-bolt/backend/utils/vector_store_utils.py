from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams
from langchain_openai import OpenAIEmbeddings
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_community.vectorstores import Qdrant
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Initialize OpenAI embeddings
embeddings = OpenAIEmbeddings(
    model="text-embedding-ada-002",
    deployment="text-embedding-ada-002",
    openai_api_key=os.getenv("OPENAI_API_KEY"),
    openai_api_version=os.getenv("OPENAI_API_VERSION"),
    openai_api_base=os.getenv("OPENAI_API_BASE"),
    openai_api_type=os.getenv("OPENAI_API_TYPE"),
)

# Initialize Qdrant client
qdrant_client = QdrantClient(url=os.getenv("QDRANT_URL", "http://localhost:6333"))
collection_name = os.getenv("QDRANT_COLLECTION_NAME", "italian_learning")

def initialize_vector_store():
    """
    Initialize or get the Qdrant vector store
    """
    # Check if collection exists
    collections = qdrant_client.get_collections().collections
    collection_names = [collection.name for collection in collections]
    
    if collection_name not in collection_names:
        # Create new collection
        qdrant_client.create_collection(
            collection_name=collection_name,
            vectors_config=VectorParams(size=1536, distance=Distance.COSINE),
        )
    
    # Return Qdrant vector store
    return Qdrant(
        client=qdrant_client,
        collection_name=collection_name,
        embeddings=embeddings,
    )

def add_text_to_vector_store(text, metadata=None):
    """
    Add text to vector store
    """
    # Initialize vector store
    vector_store = initialize_vector_store()
    
    # Split text into chunks
    text_splitter = RecursiveCharacterTextSplitter(
        chunk_size=1000,
        chunk_overlap=200,
        length_function=len,
    )
    chunks = text_splitter.split_text(text)
    
    # Create metadata for each chunk
    if metadata:
        metadatas = [metadata for _ in chunks]
    else:
        metadatas = None
    
    # Add chunks to vector store
    vector_store.add_texts(texts=chunks, metadatas=metadatas)
    
    return len(chunks)

def search_vector_store(query, filter=None, k=5):
    """
    Search vector store
    """
    # Initialize vector store
    vector_store = initialize_vector_store()
    
    # Search vector store
    results = vector_store.similarity_search(
        query=query,
        k=k,
        filter=filter,
    )
    
    return results