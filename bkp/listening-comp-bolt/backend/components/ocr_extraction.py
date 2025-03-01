import gradio as gr
import os
import re
import cv2
import pytesseract
import subprocess
from PIL import Image
import numpy as np

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
    Download YouTube video using yt-dlp if not already downloaded
    """
    video_id = extract_video_id(url)
    if not video_id:
        return {"success": False, "message": "Invalid YouTube URL"}
    
    output_path = f"output/transcript-vid/{video_id}.mp4"
    
    # Check if video already exists
    if os.path.exists(output_path):
        return {"success": True, "message": "Video already downloaded", "path": output_path, "video_id": video_id}
    
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

def extract_frames(video_path, frame_interval=1):
    """
    Extract frames from video at specified intervals
    """
    try:
        # Open the video file
        video = cv2.VideoCapture(video_path)
        
        # Get video properties
        fps = video.get(cv2.CAP_PROP_FPS)
        total_frames = int(video.get(cv2.CAP_PROP_FRAME_COUNT))
        
        # Calculate frame extraction interval
        frame_step = int(fps * frame_interval)
        
        frames = []
        frame_count = 0
        
        while True:
            ret, frame = video.read()
            
            if not ret:
                break
            
            # Extract frame at specified interval
            if frame_count % frame_step == 0:
                frames.append(frame)
            
            frame_count += 1
        
        video.release()
        
        return {"success": True, "message": f"Extracted {len(frames)} frames", "frames": frames}
    except Exception as e:
        return {"success": False, "message": f"Error extracting frames: {str(e)}"}

def perform_ocr(frames, language="ita"):
    """
    Perform OCR on extracted frames
    """
    try:
        # Configure pytesseract language
        config = f"-l {language} --oem 1 --psm 3"
        
        all_text = ""
        
        for i, frame in enumerate(frames):
            # Convert frame to grayscale
            gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
            
            # Apply thresholding to enhance text
            _, thresh = cv2.threshold(gray, 150, 255, cv2.THRESH_BINARY_INV)
            
            # Perform OCR
            text = pytesseract.image_to_string(thresh, config=config)
            
            # Add non-empty text to results
            if text.strip():
                all_text += f"[Frame {i}] {text}\n"
        
        return {"success": True, "message": "OCR completed successfully", "text": all_text}
    except Exception as e:
        return {"success": False, "message": f"Error performing OCR: {str(e)}"}

def process_video_for_ocr(url, frame_interval=1, language="ita"):
    """
    Process YouTube video for OCR extraction
    """
    # First download the video
    download_result = download_youtube_video(url)
    
    if not download_result["success"]:
        return download_result["message"], ""
    
    # Extract frames from the video
    frames_result = extract_frames(download_result["path"], frame_interval)
    
    if not frames_result["success"]:
        return frames_result["message"], ""
    
    # Perform OCR on the extracted frames
    ocr_result = perform_ocr(frames_result["frames"], language)
    
    if not ocr_result["success"]:
        return ocr_result["message"], ""
    
    # Save the OCR text to a file
    video_id = download_result["video_id"]
    output_path = f"output/transcript-ocr/{video_id}_ocr.txt"
    
    with open(output_path, "w", encoding="utf-8") as f:
        f.write(ocr_result["text"])
    
    return f"OCR completed and saved to {output_path}", ocr_result["text"]

def create_ocr_interface(parent):
    """
    Creates the OCR extraction interface
    """
    with parent:
        gr.Markdown("### OCR Extraction using OpenCV and Tesseract")
        gr.Markdown("Enter a YouTube URL to extract text from video frames using OCR.")
        
        with gr.Row():
            url_input = gr.Textbox(
                label="YouTube URL",
                placeholder="https://www.youtube.com/watch?v=...",
            )
        
        with gr.Row():
            frame_interval = gr.Slider(
                minimum=0.5,
                maximum=10,
                value=1,
                step=0.5,
                label="Frame Interval (seconds)",
                info="Interval between frames to extract"
            )
            
            language_input = gr.Dropdown(
                label="OCR Language",
                choices=["ita", "eng", "fra", "spa", "deu"],
                value="ita",
                info="Language for OCR processing"
            )
            
            process_btn = gr.Button("Process Video")
        
        with gr.Row():
            status_output = gr.Textbox(label="Status")
        
        with gr.Row():
            ocr_output = gr.TextArea(
                label="OCR Text",
                placeholder="Extracted text will appear here...",
                lines=15,
                max_lines=30,
            )
        
        # Set up event handlers
        process_btn.click(
            process_video_for_ocr,
            inputs=[url_input, frame_interval, language_input],
            outputs=[status_output, ocr_output],
        )
        
        return ocr_output