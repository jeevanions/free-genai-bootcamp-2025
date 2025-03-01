import gradio as gr
import os
import re
import subprocess
import whisper
from datetime import timedelta

def extract_video_id(url):
    """
    Extract the YouTube video ID from a URL
    """
    youtube_regex = r'(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})'
    match = re.search(youtube_regex, url)
    
    if match:
        return match.group(1)
    return None

def download_youtube_video(url):
    """
    Download YouTube video using yt-dlp
    """
    video_id = extract_video_id(url)
    if not video_id:
        return {"success": False, "message": "Invalid YouTube URL"}
    
    output_path = f"output/transcript-vid/{video_id}.mp4"
    
    try:
        # Create command to download video
        command = [
            "yt-dlp",
            "-f", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best",
            "-o", output_path,
            url
        ]
        
        # Execute the command
        subprocess.run(command, check=True)
        
        return {"success": True, "message": "Video downloaded successfully", "path": output_path, "video_id": video_id}
    except subprocess.CalledProcessError as e:
        return {"success": False, "message": f"Error downloading video: {str(e)}"}
    except Exception as e:
        return {"success": False, "message": f"Unexpected error: {str(e)}"}

def transcribe_with_whisper(video_path, language="it"):
    """
    Transcribe video using Whisper
    """
    try:
        # Load Whisper model
        model = whisper.load_model("base")
        
        # Transcribe the audio
        result = model.transcribe(video_path, language=language)
        
        # Format the transcript with timestamps
        formatted_transcript = ""
        for segment in result["segments"]:
            start_time = str(timedelta(seconds=int(segment["start"])))
            text = segment["text"]
            formatted_transcript += f"[{start_time}] {text}\n"
        
        # Save the transcript to a file
        video_id = os.path.basename(video_path).split('.')[0]
        output_path = f"output/transcript-vid/{video_id}_whisper.txt"
        with open(output_path, "w", encoding="utf-8") as f:
            f.write(formatted_transcript)
        
        return {
            "success": True,
            "message": f"Transcription completed and saved to {output_path}",
            "transcript": formatted_transcript
        }
    except Exception as e:
        return {"success": False, "message": f"Error during transcription: {str(e)}"}

def process_youtube_video(url, language="it"):
    """
    Download YouTube video and transcribe it with Whisper
    """
    # First download the video
    download_result = download_youtube_video(url)
    
    if not download_result["success"]:
        return download_result["message"], ""
    
    # Then transcribe it
    transcribe_result = transcribe_with_whisper(download_result["path"], language)
    
    if not transcribe_result["success"]:
        return transcribe_result["message"], ""
    
    return transcribe_result["message"], transcribe_result["transcript"]

def create_whisper_transcript_interface(parent):
    """
    Creates the Whisper transcript interface
    """
    with parent:
        gr.Markdown("### Generate Transcript using Whisper")
        gr.Markdown("Enter a YouTube URL to download the video and generate a transcript using Whisper.")
        
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
                info="Select the language of the video"
            )
            process_btn = gr.Button("Process Video")
        
        with gr.Row():
            status_output = gr.Textbox(label="Status")
        
        with gr.Row():
            transcript_output = gr.TextArea(
                label="Whisper Transcript",
                placeholder="Transcript will appear here...",
                lines=15,
                max_lines=30,
            )
        
        # Set up event handlers
        process_btn.click(
            process_youtube_video,
            inputs=[url_input, language_input],
            outputs=[status_output, transcript_output],
        )
        
        return transcript_output