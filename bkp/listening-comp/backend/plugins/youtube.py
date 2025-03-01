"""
YouTube plugin for downloading transcripts and videos.
"""
import os
from typing import Optional, Dict, Any, List
from pathlib import Path
import json
import subprocess
from youtube_transcript_api import YouTubeTranscriptApi, _errors
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

class YouTubePlugin:
    """Plugin for interacting with YouTube content."""
    
    def __init__(self):
        """Initialize the YouTube plugin with configuration."""
        self.transcript_yt_dir = Path(os.getenv("TRANSCRIPT_YT_DIR", "output/transcript-yt"))
        
        # Create output directory if it doesn't exist
        self.transcript_yt_dir.mkdir(parents=True, exist_ok=True)
    
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
    
    def get_transcript(self, youtube_url: str, languages: List[str] = ['it', 'en']) -> Optional[Dict[str, Any]]:
        """
        Get transcript for a YouTube video.
        
        Args:
            youtube_url: YouTube video URL
            languages: List of language codes to try, in order of preference
            
        Returns:
            Dictionary with transcript data or None if not available
        """
        video_id = self._get_video_id(youtube_url)
        
        try:
            # Try to get transcript in preferred languages
            transcript_list = None
            language_found = None
            
            for lang in languages:
                try:
                    transcript_list = YouTubeTranscriptApi.get_transcript(video_id, languages=[lang])
                    language_found = lang
                    break
                except _errors.NoTranscriptFound:
                    continue
            
            if transcript_list is None:
                # If no preferred language found, try to get any available transcript
                try:
                    transcript_list = YouTubeTranscriptApi.get_transcript(video_id)
                    language_found = "unknown"
                except Exception as e:
                    print(f"No transcript available: {e}")
                    return None
            
            # Process transcript
            full_text = " ".join([item['text'] for item in transcript_list])
            
            # Save transcript
            output_file = self.transcript_yt_dir / f"{video_id}.txt"
            with open(output_file, "w", encoding="utf-8") as f:
                f.write(full_text)
            
            # Save raw transcript data
            raw_output_file = self.transcript_yt_dir / f"{video_id}_raw.json"
            with open(raw_output_file, "w", encoding="utf-8") as f:
                json.dump(transcript_list, f, ensure_ascii=False, indent=2)
            
            return {
                "video_id": video_id,
                "language": language_found,
                "transcript": full_text,
                "transcript_file": str(output_file),
                "raw_transcript_file": str(raw_output_file),
                "segments": transcript_list
            }
        
        except Exception as e:
            print(f"Error getting transcript: {e}")
            return None
