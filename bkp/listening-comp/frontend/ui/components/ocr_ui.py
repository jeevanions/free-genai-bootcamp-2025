"""
OCR UI component for extracting text from video frames.
"""
import gradio as gr
from backend.services.transcription import TranscriptionService

def create_ocr_ui():
    """Create the OCR UI component."""
    transcription_service = TranscriptionService()
    
    with gr.Column():
        gr.Markdown("## OCR Text Extraction")
        gr.Markdown("Extract text from video frames using OCR.")
        
        youtube_url = gr.Textbox(
            placeholder="Enter YouTube URL (e.g., https://www.youtube.com/watch?v=...)",
            label="YouTube URL"
        )
        interval = gr.Slider(
            minimum=1,
            maximum=30,
            value=5,
            step=1,
            label="Frame Interval (seconds)"
        )
        extract_btn = gr.Button("Extract Text")
        output = gr.Textbox(label="Extracted Text", lines=10)
        status = gr.Markdown("")
        
        def extract_text(url, interval_value):
            """Extract text from video frames."""
            if not url:
                return "", "Please enter a YouTube URL"
            
            status_md = "⏳ Downloading video and extracting frames... (this may take a few minutes)"
            yield "", status_md
            
            result = transcription_service.extract_ocr_text(url, int(interval_value))
            
            if result["success"]:
                return result.get("text", ""), f"✅ {result['message']}"
            else:
                return "", f"❌ {result['message']}"
        
        # Connect UI components
        extract_btn.click(extract_text, [youtube_url, interval], [output, status])
    
    return gr.Column()
