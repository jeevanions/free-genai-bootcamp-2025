import re
import subprocess
from youtube_transcript_api import YouTubeTranscriptApi

def extract_video_id(url):
    """
    Extract the YouTube video ID from a URL
    """
    youtube_regex = r'(?:youtube\.com\/(?:[^\/\n\s]+\/\S+\/|(?:v|e(?:mbed)?)\/|\S*?[?&]v=)|youtu\.be\/)([a-zA-Z0-9_-]{11})'
    match = re.search(youtube_regex, url)
    
    if match:
        return match.group(1)
    return None

def download_youtube_video(url, output_path):
    """
    Download YouTube video using yt-dlp
    """
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
        
        return {"success": True, "message": "Video downloaded successfully", "path": output_path}
    except subprocess.CalledProcessError as e:
        return {"success": False, "message": f"Error downloading video: {str(e)}"}
    except Exception as e:
        return {"success": False, "message": f"Unexpected error: {str(e)}"}

def get_youtube_transcript(video_id, language="it"):
    """
    Get transcript from YouTube for a given video ID
    """
    try:
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
        
        return {"success": True, "transcript": transcript_data}
    except Exception as e:
        return {"success": False, "message": f"Error getting transcript: {str(e)}"}