"""
YouTube transcript UI component.
"""
import gradio as gr
from backend.services.transcription import TranscriptionService

def create_youtube_ui():
    """Create the YouTube transcript UI component."""
    transcription_service = TranscriptionService()
    
    with gr.Column():
        gr.Markdown("## YouTube Transcript Downloader")
        gr.Markdown("Download transcripts directly from YouTube videos.")
        
        youtube_url = gr.Textbox(
            placeholder="Enter YouTube URL (e.g., https://www.youtube.com/watch?v=...)",
            label="YouTube URL"
        )
        download_btn = gr.Button("Download Transcript")
        output = gr.Textbox(label="Transcript", lines=10)
        status = gr.Markdown("")
        
        def download_transcript(url):
            """Download transcript from YouTube."""
            if not url:
                return "", "Please enter a YouTube URL"
            
            result = transcription_service.get_youtube_transcript(url)
            
            if result["success"]:
                return result.get("transcript", ""), f"✅ {result['message']}"
            else:
                return "", f"❌ {result['message']}"
        
        # Connect UI components
        download_btn.click(download_transcript, [youtube_url], [output, status])
    
    return gr.Column()
