"""
Knowledge base service for managing RAG functionality.
"""
from pathlib import Path
from typing import Dict, Any, List, Optional
import os

from backend.components.rag import RAGComponent

class KnowledgeBaseService:
    """Service for managing knowledge base operations."""
    
    def __init__(self):
        """Initialize the knowledge base service."""
        self.rag = RAGComponent()
    
    def add_transcript_to_kb(self, transcript_path: str, metadata: Dict[str, Any]) -> Dict[str, Any]:
        """
        Add a transcript to the knowledge base.
        
        Args:
            transcript_path: Path to transcript file
            metadata: Metadata for the transcript
            
        Returns:
            Dictionary with status and document IDs
        """
        try:
            doc_ids = self.rag.add_transcript(Path(transcript_path), metadata)
            
            return {
                "success": True,
                "document_count": len(doc_ids),
                "document_ids": doc_ids,
                "message": f"Successfully added {len(doc_ids)} document chunks to knowledge base"
            }
        except Exception as e:
            return {
                "success": False,
                "message": f"Failed to add transcript to knowledge base: {str(e)}"
            }
    
    def query_knowledge_base(self, query_text: str, top_k: int = 3) -> Dict[str, Any]:
        """
        Query the knowledge base.
        
        Args:
            query_text: Query text
            top_k: Number of results to return
            
        Returns:
            Dictionary with query results
        """
        try:
            results = self.rag.query(query_text, top_k)
            
            return {
                "success": True,
                "query": query_text,
                "results": results,
                "result_count": len(results),
                "message": f"Found {len(results)} relevant documents"
            }
        except Exception as e:
            return {
                "success": False,
                "query": query_text,
                "message": f"Failed to query knowledge base: {str(e)}"
            }
