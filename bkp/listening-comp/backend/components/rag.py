"""
RAG (Retrieval-Augmented Generation) component for knowledge base management.
"""
import os
import time
from typing import List, Dict, Any, Optional
from pathlib import Path
import json
import requests
from qdrant_client import QdrantClient
from qdrant_client.http import models
from qdrant_client.http.exceptions import ResponseHandlingException
from dotenv import load_dotenv
import openai

# Load environment variables
load_dotenv()

class RAGComponent:
    """Component for managing knowledge base and RAG functionality."""
    
    def __init__(self):
        """Initialize the RAG component with Qdrant and embedding configuration."""
        # Configure Qdrant
        self.qdrant_url = os.getenv("QDRANT_URL", "http://localhost:6333")
        self.collection_name = os.getenv("QDRANT_COLLECTION_NAME", "listening_comp")
        self.qdrant_available = False
        
        # Configure Azure OpenAI for embeddings
        openai.api_type = "azure"
        openai.api_key = os.getenv("AZURE_OPENAI_API_KEY")
        openai.api_base = os.getenv("AZURE_OPENAI_ENDPOINT")
        openai.api_version = "2023-05-15"  # Update this to the latest version if needed
        self.deployment_name = os.getenv("AZURE_OPENAI_DEPLOYMENT_NAME")
        
        # Initialize Qdrant client
        try:
            # Check if Qdrant is available
            self._check_qdrant_availability()
            
            # Initialize client
            self.client = QdrantClient(url=self.qdrant_url)
            
            # Create collection if it doesn't exist
            self._create_collection_if_not_exists()
            
            self.qdrant_available = True
            print(f"Successfully connected to Qdrant at {self.qdrant_url}")
        except Exception as e:
            print(f"Warning: Could not connect to Qdrant: {e}. RAG functionality will be limited.")
            self.client = None
    
    def _check_qdrant_availability(self):
        """Check if Qdrant is available by making a simple HTTP request."""
        try:
            # Try the health endpoint first (newer versions)
            response = requests.get(f"{self.qdrant_url}/health", timeout=5)
            if response.status_code == 200:
                return
                
            # Fall back to collections endpoint (always available if Qdrant is running)
            response = requests.get(f"{self.qdrant_url}/collections", timeout=5)
            if response.status_code != 200:
                raise ConnectionError(f"Qdrant returned status code {response.status_code}")
        except requests.RequestException as e:
            raise ConnectionError(f"Could not connect to Qdrant: {e}")
    
    def _create_collection_if_not_exists(self):
        """Create Qdrant collection if it doesn't exist."""
        if not self.qdrant_available or self.client is None:
            return
            
        collections = self.client.get_collections().collections
        collection_names = [collection.name for collection in collections]
        
        if self.collection_name not in collection_names:
            self.client.create_collection(
                collection_name=self.collection_name,
                vectors_config=models.VectorParams(
                    size=1536,  # OpenAI embedding dimension
                    distance=models.Distance.COSINE
                )
            )
    
    def _get_embedding(self, text: str) -> List[float]:
        """
        Get embedding for text using Azure OpenAI.
        
        Args:
            text: Text to embed
            
        Returns:
            Embedding vector
        """
        response = openai.Embedding.create(
            input=text,
            engine=self.deployment_name  # Use the same deployment or a dedicated embedding model
        )
        return response['data'][0]['embedding']
    
    def add_document(self, document: str, metadata: Dict[str, Any]) -> str:
        """
        Add a document to the knowledge base.
        
        Args:
            document: Document text
            metadata: Document metadata (e.g., source, timestamp)
            
        Returns:
            ID of the added document
        """
        if not self.qdrant_available or self.client is None:
            return f"qdrant-not-available-{metadata.get('video_id', 'doc')}"
            
        # Generate a unique ID
        doc_id = f"{metadata.get('video_id', 'doc')}_{metadata.get('type', 'unknown')}_{metadata.get('timestamp', '')}"
        
        try:
            # Get embedding
            embedding = self._get_embedding(document)
            
            # Add to Qdrant
            self.client.upsert(
                collection_name=self.collection_name,
                points=[
                    models.PointStruct(
                        id=doc_id,
                        vector=embedding,
                        payload={
                            "text": document,
                            **metadata
                        }
                    )
                ]
            )
            
            return doc_id
        except Exception as e:
            print(f"Error adding document to Qdrant: {e}")
            return f"error-{doc_id}"
    
    def query(self, query_text: str, top_k: int = 3) -> List[Dict[str, Any]]:
        """
        Query the knowledge base.
        
        Args:
            query_text: Query text
            top_k: Number of results to return
            
        Returns:
            List of relevant documents with metadata
        """
        if not self.qdrant_available or self.client is None:
            return [{
                "text": "Qdrant is not available. Please start the Qdrant service using Docker.",
                "score": 0.0,
                "metadata": {"error": "qdrant_not_available"}
            }]
            
        try:
            # Get query embedding
            query_embedding = self._get_embedding(query_text)
            
            # Search in Qdrant
            search_result = self.client.search(
                collection_name=self.collection_name,
                query_vector=query_embedding,
                limit=top_k
            )
            
            # Format results
            results = []
            for hit in search_result:
                payload = hit.payload
                results.append({
                    "text": payload.get("text", ""),
                    "score": hit.score,
                    "metadata": {k: v for k, v in payload.items() if k != "text"}
                })
            
            return results
        except Exception as e:
            print(f"Error querying Qdrant: {e}")
            return [{
                "text": f"Error querying knowledge base: {str(e)}",
                "score": 0.0,
                "metadata": {"error": "query_error"}
            }]
    
    def add_transcript(self, transcript_path: Path, metadata: Dict[str, Any]) -> List[str]:
        """
        Add a transcript to the knowledge base by chunking it.
        
        Args:
            transcript_path: Path to transcript file
            metadata: Metadata for the transcript
            
        Returns:
            List of document IDs added
        """
        if not self.qdrant_available or self.client is None:
            return [f"qdrant-not-available-{metadata.get('video_id', 'doc')}"]
            
        try:
            # Read transcript
            with open(transcript_path, "r", encoding="utf-8") as f:
                transcript = f.read()
            
            # Chunk the transcript (simple chunking by paragraphs)
            chunks = [chunk.strip() for chunk in transcript.split("\n\n") if chunk.strip()]
            
            # Add each chunk
            doc_ids = []
            for i, chunk in enumerate(chunks):
                chunk_metadata = {
                    **metadata,
                    "chunk_index": i,
                    "total_chunks": len(chunks)
                }
                doc_id = self.add_document(chunk, chunk_metadata)
                doc_ids.append(doc_id)
            
            return doc_ids
        except Exception as e:
            print(f"Error adding transcript to knowledge base: {e}")
            return [f"error-{metadata.get('video_id', 'doc')}"]
