import gradio as gr
import os
import re
from youtube_transcript_api import YouTubeTranscriptApi, TranscriptsDisabled

def extract_video_id(url):
    """
    Extract the YouTube video ID from a URL
    """
    # Regular expression to match YouTube video IDs
    youtube_regex = r'(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})'
    match = re.search(youtube_regex, url)
    
    if match:
        return match.group(1)
    return None

def download_youtube_transcript(url, language="it", status_box=None):
    """
    Download transcript from YouTube for a given video URL
    """
    try:
        # Update status if status_box is provided
        if status_box:
            status_box.value = "Downloading transcript... Please wait."
            status_box.elem_classes = "status-msg status-info"
            
        video_id = extract_video_id(url)
        if not video_id:
            error_msg = "Invalid YouTube URL. Could not extract video ID."
            if status_box:
                status_box.value = error_msg
                status_box.elem_classes = "status-msg status-error"
            return {"success": False, "message": error_msg}
        
        # Get available transcripts
        transcript_list = YouTubeTranscriptApi.list_transcripts(video_id)
        
        # Try to get the transcript in the specified language
        try:
            transcript = transcript_list.find_transcript([language])
        except:
            # If not available, try to get the auto-generated one and translate it
            try:
                transcript = transcript_list.find_transcript(['en'])
                transcript = transcript.translate(language)
            except:
                # If still not available, get any available transcript
                transcript = transcript_list.find_transcript([])
        
        # Get the transcript data
        transcript_data = transcript.fetch()
        
        # Format the transcript
        formatted_transcript = ""
        for entry in transcript_data:
            start_time = entry['start']
            text = entry['text']
            formatted_transcript += f"[{format_time(start_time)}] {text}\n"
        
        # Save the transcript to a file
        os.makedirs("output/transcript-yt", exist_ok=True)
        output_path = f"output/transcript-yt/{video_id}_transcript.txt"
        with open(output_path, "w", encoding="utf-8") as f:
            f.write(formatted_transcript)
        
        success_msg = f"Transcript downloaded successfully and saved to {output_path}"
        if status_box:
            status_box.value = "✅ " + success_msg
            status_box.elem_classes = "status-msg status-success"
            
        return {
            "success": True, 
            "message": success_msg,
            "transcript": formatted_transcript,
            "video_id": video_id
        }
    
    except TranscriptsDisabled:
        error_msg = "Transcripts are disabled for this video."
        if status_box:
            status_box.value = "❌ " + error_msg
            status_box.elem_classes = "status-msg status-error"
        return {"success": False, "message": error_msg}
    except Exception as e:
        error_msg = f"Error: {str(e)}"
        if status_box:
            status_box.value = "❌ " + error_msg
            status_box.elem_classes = "status-msg status-error"
        return {"success": False, "message": error_msg}

def format_time(seconds):
    """
    Format time in seconds to MM:SS format
    """
    minutes = int(seconds // 60)
    seconds = int(seconds % 60)
    return f"{minutes:02d}:{seconds:02d}"

def create_youtube_transcript_interface(parent):
    """
    Creates the YouTube transcript interface
    """
    with parent:
        with gr.Group(elem_classes="content-section"):
            gr.Markdown("### YouTube Transcript Downloader")
            gr.Markdown("Enter a YouTube URL to download its transcript. Preferably Italian language videos.")
            
            with gr.Row():
                url_input = gr.Textbox(
                    label="YouTube URL",
                    placeholder="https://www.youtube.com/watch?v=...",
                )
            
            with gr.Row():
                language_input = gr.Dropdown(
                    label="Language",
                    choices=["it", "en", "fr", "es", "de"],
                    value="it",
                    info="Select the language of the transcript"
                )
                download_btn = gr.Button("Download Transcript", elem_classes="sidebar-btn")
            
            with gr.Row():
                output_message = gr.Textbox(label="Status", elem_classes="status-msg")
            
            with gr.Row():
                transcript_output = gr.TextArea(
                    label="Transcript",
                    placeholder="Transcript will appear here...",
                    lines=15,
                    max_lines=30,
                )
        
        # Set up event handlers
        download_btn.click(
            download_youtube_transcript,
            inputs=[url_input, language_input],
            outputs=[transcript_output],
        )
        
        return transcript_output