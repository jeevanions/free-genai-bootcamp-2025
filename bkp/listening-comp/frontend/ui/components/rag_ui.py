"""
RAG UI component for knowledge base queries.
"""
import gradio as gr
from backend.services.knowledge_base import KnowledgeBaseService
from backend.services.transcription import TranscriptionService

def create_rag_ui():
    """Create the RAG UI component."""
    kb_service = KnowledgeBaseService()
    transcription_service = TranscriptionService()
    
    with gr.Column():
        gr.Markdown("## Knowledge Base")
        gr.Markdown("Add transcripts to the knowledge base and query for information.")
        
        with gr.Tab("Add to Knowledge Base"):
            youtube_url = gr.Textbox(
                placeholder="Enter YouTube URL (e.g., https://www.youtube.com/watch?v=...)",
                label="YouTube URL"
            )
            source_type = gr.Radio(
                choices=["YouTube Transcript", "Whisper Transcript", "OCR Text"],
                value="YouTube Transcript",
                label="Source Type"
            )
            add_btn = gr.Button("Add to Knowledge Base")
            add_status = gr.Markdown("")
            
            def add_to_kb(url, source):
                """Add transcript to knowledge base."""
                if not url:
                    return "Please enter a YouTube URL"
                
                video_id = transcription_service.youtube._get_video_id(url)
                
                # Get transcript based on source type
                if source == "YouTube Transcript":
                    result = transcription_service.get_youtube_transcript(url)
                elif source == "Whisper Transcript":
                    result = transcription_service.generate_whisper_transcript(url)
                else:  # OCR Text
                    result = transcription_service.extract_ocr_text(url)
                
                if not result["success"]:
                    return f"❌ {result['message']}"
                
                # Add to knowledge base
                metadata = {
                    "video_id": video_id,
                    "source_type": source,
                    "url": url
                }
                
                kb_result = kb_service.add_transcript_to_kb(
                    result["transcript_file"],
                    metadata
                )
                
                if kb_result["success"]:
                    return f"✅ {kb_result['message']}"
                else:
                    return f"❌ {kb_result['message']}"
        
        with gr.Tab("Query Knowledge Base"):
            query = gr.Textbox(
                placeholder="Enter your query...",
                label="Query"
            )
            top_k = gr.Slider(
                minimum=1,
                maximum=10,
                value=3,
                step=1,
                label="Number of Results"
            )
            query_btn = gr.Button("Search")
            results = gr.JSON(label="Results")
            query_status = gr.Markdown("")
            
            def query_kb(query_text, k):
                """Query the knowledge base."""
                if not query_text:
                    return None, "Please enter a query"
                
                result = kb_service.query_knowledge_base(query_text, int(k))
                
                if result["success"]:
                    return result["results"], f"✅ {result['message']}"
                else:
                    return None, f"❌ {result['message']}"
        
        # Connect UI components
        add_btn.click(add_to_kb, [youtube_url, source_type], [add_status])
        query_btn.click(query_kb, [query, top_k], [results, query_status])
    
    return gr.Column()
