import os
import re
import cv2
import pytesseract
import subprocess
import tempfile
from typing import List, Optional
from fastapi import FastAPI, File, UploadFile, Form, HTTPException
from fastapi.responses import JSONResponse
import uvicorn
from pydantic import BaseModel

app = FastAPI(title="Tesseract OCR Service")

# Ensure output directories exist
os.makedirs("/app/output/transcript-vid", exist_ok=True)
os.makedirs("/app/output/transcript-ocr", exist_ok=True)

class OCRResponse(BaseModel):
    success: bool
    message: str
    text: Optional[str] = None

@app.get("/")
def read_root():
    return {"status": "Tesseract OCR Service is running"}

@app.get("/version")
def get_version():
    """Get Tesseract version information"""
    try:
        version = pytesseract.get_tesseract_version()
        return {"success": True, "version": str(version)}
    except Exception as e:
        return {"success": False, "error": str(e)}

@app.post("/ocr/image", response_model=OCRResponse)
async def ocr_image(
    file: UploadFile = File(...),
    language: str = Form("eng")
):
    """Perform OCR on an uploaded image"""
    try:
        # Save uploaded file to temp location
        temp_file = tempfile.NamedTemporaryFile(delete=False, suffix=".jpg")
        temp_file.write(await file.read())
        temp_file.close()
        
        # Perform OCR
        img = cv2.imread(temp_file.name)
        gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
        _, thresh = cv2.threshold(gray, 150, 255, cv2.THRESH_BINARY_INV)
        
        text = pytesseract.image_to_string(thresh, lang=language)
        
        # Clean up
        os.unlink(temp_file.name)
        
        return OCRResponse(
            success=True,
            message="OCR completed successfully",
            text=text
        )
    except Exception as e:
        # Clean up in case of error
        if temp_file and os.path.exists(temp_file.name):
            os.unlink(temp_file.name)
        return OCRResponse(
            success=False,
            message=f"Error performing OCR: {str(e)}"
        )

def extract_video_id(url):
    """Extract the YouTube video ID from a URL"""
    youtube_regex = r'(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})'
    match = re.search(youtube_regex, url)
    
    if match:
        return match.group(1)
    return None

def download_youtube_video(url):
    """Download YouTube video using yt-dlp"""
    video_id = extract_video_id(url)
    if not video_id:
        return {"success": False, "message": "Invalid YouTube URL"}
    
    output_path = f"/app/output/transcript-vid/{video_id}.mp4"
    
    # Check if video already exists
    if os.path.exists(output_path):
        return {"success": True, "message": "Video already downloaded", "path": output_path, "video_id": video_id}
    
    try:
        # Create command to download video
        command = [
            "yt-dlp",
            "-f", "mp4",
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
    """Extract frames from video at specified intervals"""
    try:
        # Open the video file
        video = cv2.VideoCapture(video_path)
        
        # Get video properties
        fps = video.get(cv2.CAP_PROP_FPS)
        
        # Calculate frame extraction interval
        frame_step = int(fps * frame_interval)
        if frame_step <= 0:
            frame_step = 1
        
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

def perform_ocr(frames, language="eng"):
    """Perform OCR on extracted frames"""
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

class VideoOCRRequest(BaseModel):
    url: str
    frame_interval: float = 1.0
    language: str = "eng"

@app.post("/ocr/youtube", response_model=OCRResponse)
async def ocr_youtube(request: VideoOCRRequest):
    """
    Process a YouTube video and extract text using OCR
    
    - url: YouTube URL
    - frame_interval: Interval between frames to extract (in seconds)
    - language: Language for OCR (eng, ita, fra, spa, deu)
    """
    try:
        # Download the video
        download_result = download_youtube_video(request.url)
        
        if not download_result["success"]:
            return OCRResponse(
                success=False,
                message=download_result["message"]
            )
        
        # Extract frames from the video
        frames_result = extract_frames(download_result["path"], request.frame_interval)
        
        if not frames_result["success"]:
            return OCRResponse(
                success=False,
                message=frames_result["message"]
            )
        
        # Perform OCR on the extracted frames
        ocr_result = perform_ocr(frames_result["frames"], request.language)
        
        if not ocr_result["success"]:
            return OCRResponse(
                success=False,
                message=ocr_result["message"]
            )
        
        # Save the OCR text to a file
        video_id = download_result["video_id"]
        output_path = f"/app/output/transcript-ocr/{video_id}_ocr.txt"
        
        with open(output_path, "w", encoding="utf-8") as f:
            f.write(ocr_result["text"])
        
        return OCRResponse(
            success=True,
            message=f"OCR completed and saved to {output_path}",
            text=ocr_result["text"]
        )
    except Exception as e:
        return OCRResponse(
            success=False,
            message=f"Error: {str(e)}"
        )

if __name__ == "__main__":
    uvicorn.run("app:app", host="0.0.0.0", port=8000, reload=True)