import gradio as gr
import os
import re
import cv2
import requests
import json
import traceback
from PIL import Image
import numpy as np

# Tesseract service URL - change to your Docker container address
TESSERACT_SERVICE_URL = os.getenv("TESSERACT_SERVICE_URL", "http://localhost:8000")

def extract_video_id(url):
    """
    Extract the YouTube video ID from a URL
    """
    youtube_regex = r'(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})'
    match = re.search(youtube_regex, url)
    
    if match:
        return match.group(1)
    return None

def process_video_for_ocr(url, frame_interval=1, language="ita"):
    """
    Process YouTube video for OCR extraction by calling the Tesseract service
    """
    try:
        print(f"Processing video: {url}")
        print(f"Connecting to Tesseract service at: {TESSERACT_SERVICE_URL}")
        
        # Check if service is available
        try:
            response = requests.get(f"{TESSERACT_SERVICE_URL}/version")
            if response.status_code == 200:
                data = response.json()
                print(f"Tesseract Service is running. Version: {data.get('version', 'unknown')}")
            else:
                return f"Tesseract Service not available: Status code {response.status_code}", ""
        except requests.RequestException as e:
            return f"Could not connect to Tesseract service: {str(e)}", ""
        
        # Call the OCR service for YouTube video processing
        payload = {
            "url": url,
            "frame_interval": float(frame_interval),
            "language": language
        }
        
        print(f"Sending request to Tesseract service with parameters: {payload}")
        response = requests.post(
            f"{TESSERACT_SERVICE_URL}/ocr/youtube",
            json=payload
        )
        
        if response.status_code != 200:
            return f"Error from Tesseract service: Status code {response.status_code}", ""
        
        result = response.json()
        
        if not result.get("success", False):
            return result.get("message", "Unknown error"), ""
        
        # Save a local copy of the OCR text
        video_id = extract_video_id(url)
        if video_id:
            os.makedirs("output/transcript-ocr", exist_ok=True)
            output_path = f"output/transcript-ocr/{video_id}_ocr.txt"
            
            with open(output_path, "w", encoding="utf-8") as f:
                f.write(result.get("text", ""))
            
            print(f"Saved local copy of OCR results to {output_path}")
        
        return result.get("message", "OCR completed"), result.get("text", "")
    except Exception as e:
        error_msg = f"Error in process_video_for_ocr: {str(e)}"
        print(error_msg)
        print(traceback.format_exc())
        return error_msg, ""

def create_ocr_interface(parent):
    """
    Creates the OCR extraction interface using the Tesseract Docker service
    """
    with parent:
        # Header for the OCR extraction section
        with gr.Group(elem_classes="chat-with-assistant"):
            gr.Markdown("## OCR Extraction using Tesseract Docker Service", elem_classes="dark-header")
            gr.Markdown("Enter a YouTube URL to extract text from video frames using OCR.", elem_classes="chat-description dark-description")
        
        with gr.Group():
            with gr.Row():
                url_input = gr.Textbox(
                    label="YouTube URL",
                    placeholder="https://www.youtube.com/watch?v=...",
                    elem_classes="dark-textbox"
                )
            
            with gr.Row():
                frame_interval = gr.Slider(
                    minimum=0.5,
                    maximum=10,
                    value=1,
                    step=0.5,
                    label="Frame Interval (seconds)",
                    info="Interval between frames to extract",
                    elem_classes="dark-slider"
                )
                
                language_input = gr.Dropdown(
                    label="OCR Language",
                    choices=["ita", "eng", "fra", "spa", "deu"],
                    value="ita",
                    info="Language for OCR processing",
                    elem_classes="dark-dropdown"
                )
                
                process_btn = gr.Button("Process Video", variant="primary", elem_classes="send-btn")
            
            with gr.Row():
                status_output = gr.Textbox(label="Status", elem_classes="status-box dark-textbox")
            
            with gr.Row():
                ocr_output = gr.TextArea(
                    label="OCR Text",
                    placeholder="Extracted text will appear here...",
                    lines=15,
                    max_lines=30,
                    elem_classes="dark-textarea"
                )
            
            # Debug output
            with gr.Accordion("Debug Information", open=False):
                debug_info = gr.Textbox(
                    label="Processing Log",
                    value="Debug information will appear here during processing",
                    lines=10,
                    elem_classes="dark-textbox"
                )
        
        # Set up event handlers
        def process_with_debug(url, interval, language):
            try:
                # Redirect print output to capture logs
                import io
                import sys
                old_stdout = sys.stdout
                new_stdout = io.StringIO()
                sys.stdout = new_stdout
                
                # Process the video
                status, transcript = process_video_for_ocr(url, interval, language)
                
                # Get the captured output
                sys.stdout = old_stdout
                debug_text = new_stdout.getvalue()
                
                return status, transcript, debug_text
            except Exception as e:
                sys.stdout = old_stdout
                error_details = traceback.format_exc()
                print(f"Error in UI handler: {str(e)}")
                print(f"Full error details: {error_details}")
                return f"Error: {str(e)}", "", error_details
        
        process_btn.click(
            process_with_debug,
            inputs=[url_input, frame_interval, language_input],
            outputs=[status_output, ocr_output, debug_info],
        )
        
        return ocr_output