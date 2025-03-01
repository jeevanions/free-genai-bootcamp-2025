"""
Media processing component for video and audio handling.
"""
import os
import cv2
import whisper
import subprocess
from pathlib import Path
from typing import List, Optional, Tuple
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

class MediaComponent:
    """Component for processing video and audio content."""
    
    def __init__(self):
        """Initialize the media component with configuration."""
        self.transcript_vid_dir = Path(os.getenv("TRANSCRIPT_VID_DIR", "output/transcript-vid"))
        self.transcript_ocr_dir = Path(os.getenv("TRANSCRIPT_OCR_DIR", "output/transcript-ocr"))
        self.whisper_model_size = os.getenv("WHISPER_MODEL", "base")
        
        # Create output directories if they don't exist
        self.transcript_vid_dir.mkdir(parents=True, exist_ok=True)
        self.transcript_ocr_dir.mkdir(parents=True, exist_ok=True)
        
        # Load Whisper model
        self.whisper_model = None  # Lazy loading to save memory
    
    def _get_video_id(self, youtube_url: str) -> str:
        """
        Extract video ID from YouTube URL.
        
        Args:
            youtube_url: YouTube video URL
            
        Returns:
            YouTube video ID
        """
        if "youtu.be" in youtube_url:
            return youtube_url.split("/")[-1].split("?")[0]
        elif "youtube.com" in youtube_url:
            if "v=" in youtube_url:
                return youtube_url.split("v=")[1].split("&")[0]
            elif "embed" in youtube_url:
                return youtube_url.split("/")[-1].split("?")[0]
        return youtube_url  # Return the URL if no ID found
    
    def download_youtube_video(self, youtube_url: str) -> Optional[str]:
        """
        Download a YouTube video using yt-dlp.
        
        Args:
            youtube_url: YouTube video URL
            
        Returns:
            Path to downloaded video file or None if download failed
        """
        video_id = self._get_video_id(youtube_url)
        output_path = f"/tmp/{video_id}.mp4"
        
        try:
            # Use yt-dlp to download the video
            command = [
                "yt-dlp",
                "-f", "best[height<=720]",  # Limit resolution to 720p
                "-o", output_path,
                youtube_url
            ]
            
            subprocess.run(command, check=True, capture_output=True)
            return output_path
        except subprocess.CalledProcessError as e:
            print(f"Error downloading video: {e}")
            return None
    
    def transcribe_audio(self, video_path: str, video_id: str) -> Optional[str]:
        """
        Transcribe audio from a video file using Whisper.
        
        Args:
            video_path: Path to the video file
            video_id: YouTube video ID for output file naming
            
        Returns:
            Path to the transcript file or None if transcription failed
        """
        try:
            # Lazy load the Whisper model
            if self.whisper_model is None:
                self.whisper_model = whisper.load_model(self.whisper_model_size)
            
            # Transcribe the audio
            result = self.whisper_model.transcribe(video_path)
            
            # Save the transcript
            output_file = self.transcript_vid_dir / f"{video_id}.txt"
            with open(output_file, "w", encoding="utf-8") as f:
                f.write(result["text"])
            
            return str(output_file)
        except Exception as e:
            print(f"Error transcribing audio: {e}")
            return None
    
    def extract_frames(self, video_path: str, interval: int = 5) -> List[str]:
        """
        Extract frames from a video at specified intervals.
        
        Args:
            video_path: Path to the video file
            interval: Interval in seconds between frames
            
        Returns:
            List of paths to extracted frame images
        """
        try:
            video = cv2.VideoCapture(video_path)
            fps = video.get(cv2.CAP_PROP_FPS)
            frame_interval = int(fps * interval)
            
            count = 0
            frame_paths = []
            video_id = Path(video_path).stem
            
            # Create a temporary directory for frames
            frames_dir = Path(f"/tmp/frames_{video_id}")
            frames_dir.mkdir(parents=True, exist_ok=True)
            
            while video.isOpened():
                ret, frame = video.read()
                if not ret:
                    break
                
                if count % frame_interval == 0:
                    frame_path = frames_dir / f"frame_{count}.jpg"
                    cv2.imwrite(str(frame_path), frame)
                    frame_paths.append(str(frame_path))
                
                count += 1
            
            video.release()
            return frame_paths
        except Exception as e:
            print(f"Error extracting frames: {e}")
            return []
