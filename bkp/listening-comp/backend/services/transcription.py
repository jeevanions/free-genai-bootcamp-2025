"""
Transcription service for processing YouTube videos.
"""
from pathlib import Path
from typing import Dict, Any, Optional
import os

from backend.components.media import MediaComponent
from backend.plugins.youtube import YouTubePlugin
from backend.plugins.ocr import OCRPlugin

class TranscriptionService:
    """Service for handling video transcription workflows."""
    
    def __init__(self):
        """Initialize the transcription service."""
        self.media = MediaComponent()
        self.youtube = YouTubePlugin()
        self.ocr = OCRPlugin()
    
    def get_youtube_transcript(self, youtube_url: str) -> Dict[str, Any]:
        """
        Get transcript from YouTube.
        
        Args:
            youtube_url: YouTube video URL
            
        Returns:
            Dictionary with transcript data and status
        """
        video_id = self.youtube._get_video_id(youtube_url)
        result = self.youtube.get_transcript(youtube_url)
        
        if result:
            return {
                "success": True,
                "video_id": video_id,
                "transcript": result["transcript"],
                "transcript_file": result["transcript_file"],
                "message": f"Successfully downloaded transcript for video {video_id}"
            }
        else:
            return {
                "success": False,
                "video_id": video_id,
                "message": f"Failed to download transcript for video {video_id}"
            }
    
    def generate_whisper_transcript(self, youtube_url: str) -> Dict[str, Any]:
        """
        Generate transcript using Whisper.
        
        Args:
            youtube_url: YouTube video URL
            
        Returns:
            Dictionary with transcript data and status
        """
        video_id = self.youtube._get_video_id(youtube_url)
        
        # Download video
        video_path = self.media.download_youtube_video(youtube_url)
        
        if not video_path:
            return {
                "success": False,
                "video_id": video_id,
                "message": f"Failed to download video {video_id}"
            }
        
        # Transcribe audio
        transcript_path = self.media.transcribe_audio(video_path, video_id)
        
        if not transcript_path:
            return {
                "success": False,
                "video_id": video_id,
                "message": f"Failed to transcribe audio for video {video_id}"
            }
        
        # Read transcript
        with open(transcript_path, "r", encoding="utf-8") as f:
            transcript = f.read()
        
        return {
            "success": True,
            "video_id": video_id,
            "transcript": transcript,
            "transcript_file": transcript_path,
            "message": f"Successfully generated transcript for video {video_id}"
        }
    
    def extract_ocr_text(self, youtube_url: str, interval: int = 5) -> Dict[str, Any]:
        """
        Extract text from video frames using OCR.
        
        Args:
            youtube_url: YouTube video URL
            interval: Interval in seconds between frames
            
        Returns:
            Dictionary with OCR data and status
        """
        video_id = self.youtube._get_video_id(youtube_url)
        
        # Download video
        video_path = self.media.download_youtube_video(youtube_url)
        
        if not video_path:
            return {
                "success": False,
                "video_id": video_id,
                "message": f"Failed to download video {video_id}"
            }
        
        # Extract frames
        frame_paths = self.media.extract_frames(video_path, interval)
        
        if not frame_paths:
            return {
                "success": False,
                "video_id": video_id,
                "message": f"Failed to extract frames from video {video_id}"
            }
        
        # Process frames with OCR
        ocr_result = self.ocr.process_frames(frame_paths, video_id)
        
        return {
            "success": True,
            "video_id": video_id,
            "text": ocr_result["combined_text"],
            "transcript_file": ocr_result["transcript_file"],
            "frame_count": ocr_result["frame_count"],
            "message": f"Successfully extracted text from {ocr_result['frame_count']} frames for video {video_id}"
        }
