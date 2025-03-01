"""
Whisper transcription UI component.
"""
import gradio as gr
from backend.services.transcription import TranscriptionService

def create_whisper_ui():
    """Create the Whisper transcription UI component."""
    transcription_service = TranscriptionService()
    
    with gr.Column():
        gr.Markdown("## Whisper Transcription")
        gr.Markdown("Generate transcripts from YouTube videos using Whisper.")
        
        youtube_url = gr.Textbox(
            placeholder="Enter YouTube URL (e.g., https://www.youtube.com/watch?v=...)",
            label="YouTube URL"
        )
        transcribe_btn = gr.Button("Generate Transcript")
        output = gr.Textbox(label="Transcript", lines=10)
        status = gr.Markdown("")
        
        def generate_transcript(url):
            """Generate transcript using Whisper."""
            if not url:
                return "", "Please enter a YouTube URL"
            
            status_md = "⏳ Downloading video and generating transcript... (this may take a few minutes)"
            yield "", status_md
            
            result = transcription_service.generate_whisper_transcript(url)
            
            if result["success"]:
                return result.get("transcript", ""), f"✅ {result['message']}"
            else:
                return "", f"❌ {result['message']}"
        
        # Connect UI components
        transcribe_btn.click(generate_transcript, [youtube_url], [output, status])
    
    return gr.Column()
