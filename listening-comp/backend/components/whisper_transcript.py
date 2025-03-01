import gradio as gr
import os
import re
import subprocess
import openai
import traceback
from datetime import timedelta
from dotenv import load_dotenv
from pathlib import Path

# Load environment variables
load_dotenv()

# Initialize OpenAI client with detailed environment inspection
try:
    # Print configuration details for Whisper-specific environment variables
    print("\n=== Azure Whisper Configuration ===")
    print(f"Whisper API Key exists: {bool(os.getenv('OPENAI_WHISPER_API_KEY'))}")
    print(f"Whisper API Version: {os.getenv('OPENAI_WHISPER_API_VERSION')}")
    print(f"Whisper API Base URL: {os.getenv('OPENAI_WHISPER_API_BASE')}")
    print(f"Whisper Deployment Name: {os.getenv('OPENAI_WHISPER_API_DEPLOYMENT_NAME')}")
    print(f"Whisper Model Version: {os.getenv('OPENAI_WHISPER_MODEL_VERSIPM')}")
    
    # Initialize the Azure OpenAI client specifically for Whisper
    # Using the dedicated Whisper environment variables
    
    # Ensure the Azure endpoint has the https:// prefix
    whisper_endpoint = os.getenv("OPENAI_WHISPER_API_BASE")
    if whisper_endpoint:
        # Remove any trailing slashes
        whisper_endpoint = whisper_endpoint.rstrip('/')
        
        # Add https:// prefix if not present
        if not whisper_endpoint.startswith("https://"):
            whisper_endpoint = f"https://{whisper_endpoint}"
            print(f"Added https:// prefix to Whisper endpoint: {whisper_endpoint}")
    
    whisper_client = openai.AzureOpenAI(
        api_key=os.getenv("OPENAI_WHISPER_API_KEY"),
        api_version=os.getenv("OPENAI_WHISPER_API_VERSION"),
        azure_endpoint=whisper_endpoint,
    )
    
    print("Azure Whisper client initialized successfully")
    
except Exception as e:
    print(f"Error initializing Azure Whisper client: {str(e)}")
    print(f"Full error details: {traceback.format_exc()}")

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
    Download YouTube video using yt-dlp with improved path handling
    """
    video_id = extract_video_id(url)
    if not video_id:
        return {"success": False, "message": "Invalid YouTube URL"}
    
    # Create output directory using absolute path to avoid path resolution issues
    output_dir = os.path.abspath("output/transcript-vid")
    os.makedirs(output_dir, exist_ok=True)
    
    # Use absolute path for output
    output_path = os.path.join(output_dir, f"{video_id}.mp4")
    print(f"Will download to absolute path: {output_path}")
    
    try:
        # Check if yt-dlp is available
        try:
            subprocess.run(["yt-dlp", "--version"], check=True, capture_output=True)
            print("yt-dlp is installed and accessible")
        except (subprocess.CalledProcessError, FileNotFoundError) as e:
            return {"success": False, "message": f"yt-dlp is not installed or not accessible: {str(e)}"}
        
        # Create command to download video - using specific format to avoid potential issues
        command = [
            "yt-dlp",
            "--format", "mp4",  # Simpler format specification
            "--output", output_path,
            url
        ]
        
        print(f"Executing command: {' '.join(command)}")
        
        # Execute the command
        result = subprocess.run(command, check=True, capture_output=True, text=True)
        print(f"yt-dlp stdout: {result.stdout}")
        
        # Verify the file exists and has content
        if not os.path.exists(output_path):
            print(f"File not found after download. Checking for other files in directory...")
            # List all files in the directory to see if yt-dlp used a different filename
            dir_files = os.listdir(output_dir)
            print(f"Files in {output_dir}: {dir_files}")
            
            # Try to find a matching file that might have a different extension or naming pattern
            matching_files = [f for f in dir_files if video_id in f]
            if matching_files:
                print(f"Found potential matches: {matching_files}")
                # Use the first matching file as our output
                output_path = os.path.join(output_dir, matching_files[0])
                print(f"Using alternative file: {output_path}")
            else:
                return {"success": False, "message": f"Download completed but file not found at {output_path} and no matching files found"}
        
        file_size = os.path.getsize(output_path)
        print(f"Downloaded file size: {file_size} bytes")
        
        if file_size == 0:
            return {"success": False, "message": f"Download completed but file is empty (0 bytes)"}
            
        print(f"Video downloaded successfully to {output_path}")
        return {"success": True, "message": "Video downloaded successfully", "path": output_path, "video_id": video_id}
    except subprocess.CalledProcessError as e:
        error_msg = f"Error downloading video: {str(e)}"
        print(error_msg)
        print(f"Command output: {e.stdout if hasattr(e, 'stdout') else 'No output'}")
        print(f"Command error: {e.stderr if hasattr(e, 'stderr') else 'No error output'}")
        print(f"Full error details: {traceback.format_exc()}")
        return {"success": False, "message": error_msg}
    except Exception as e:
        error_msg = f"Unexpected error during download: {str(e)}"
        print(error_msg)
        print(f"Full error details: {traceback.format_exc()}")
        return {"success": False, "message": error_msg}

def transcribe_with_whisper(video_path, language="it"):
    """
    Transcribe video using OpenAI's Whisper API with custom environment variables
    """
    try:
        print(f"Starting transcription of {video_path} with language {language}")
        
        # Check if file exists and has content
        if not os.path.exists(video_path):
            error_msg = f"Error: File {video_path} does not exist"
            print(error_msg)
            return {"success": False, "message": error_msg}
            
        file_size = os.path.getsize(video_path)
        print(f"File size: {file_size} bytes")
        
        if file_size == 0:
            error_msg = f"Error: File {video_path} is empty (0 bytes)"
            print(error_msg)
            return {"success": False, "message": error_msg}
        
        # Extract audio from video to ensure we have a valid audio file (mp3 format is more reliable with Whisper)
        print("Extracting audio from video to ensure we have a valid audio file...")
        output_dir = os.path.dirname(video_path)
        audio_path = os.path.join(output_dir, f"{os.path.basename(video_path).split('.')[0]}.mp3")
        
        try:
            subprocess.run([
                "ffmpeg", 
                "-i", video_path, 
                "-q:a", "0", 
                "-map", "a", 
                "-y",  # Overwrite output file if it exists
                audio_path
            ], check=True, capture_output=True)
            
            print(f"Audio extracted to {audio_path}")
            # Use the extracted audio file instead
            transcription_file = audio_path
        except (subprocess.SubprocessError, FileNotFoundError) as e:
            print(f"Warning: Failed to extract audio: {str(e)}. Will try using the original file.")
            transcription_file = video_path
        
        # Open the audio file
        with open(transcription_file, "rb") as audio_file:
            print(f"File opened successfully: {transcription_file}")
            
            # Get the Whisper deployment name from environment variable
            whisper_deployment = os.getenv("OPENAI_WHISPER_API_DEPLOYMENT_NAME", "whisper")
            print(f"Using Whisper deployment name: {whisper_deployment}")
            
            try:
                print(f"Calling Azure Whisper API with deployment name: {whisper_deployment}")
                
                # Use the whisper_client specifically for Whisper API calls
                response = whisper_client.audio.transcriptions.create(
                    model=whisper_deployment,  # Use the deployment name from env var
                    file=audio_file,
                    language=language
                )
                
                print("Transcription completed successfully!")
                print(f"Response type: {type(response)}")
                
                # Extract the transcript text
                if hasattr(response, 'text'):
                    transcript_text = response.text
                elif isinstance(response, dict) and 'text' in response:
                    transcript_text = response['text']
                else:
                    transcript_text = str(response)
                
                print(f"Transcript preview: {transcript_text[:100]}...")
                
                # Save the transcript to a file
                video_id = os.path.basename(video_path).split('.')[0]
                output_path = os.path.join(output_dir, f"{video_id}_whisper.txt")
                with open(output_path, "w", encoding="utf-8") as f:
                    f.write(transcript_text)
                
                print(f"Transcript saved to {output_path}")
                
                return {
                    "success": True,
                    "message": f"Transcription completed and saved to {output_path}",
                    "transcript": transcript_text
                }
            except openai.NotFoundError as e:
                error_msg = f"NotFoundError: {str(e)} - The Whisper deployment '{whisper_deployment}' was not found"
                print(error_msg)
                print(f"Full error details: {traceback.format_exc()}")
                
                print("\n=== TROUBLESHOOTING STEPS FOR 404 ERROR ===")
                print(f"1. Verify the deployment name '{whisper_deployment}' is correct")
                print(f"2. Check if the Azure endpoint '{os.getenv('OPENAI_WHISPER_API_BASE')}' is correct")
                print(f"3. Ensure your API version '{os.getenv('OPENAI_WHISPER_API_VERSION')}' supports Whisper")
                print("4. Verify the API key has permissions to access this deployment")
                
                return {"success": False, "message": error_msg}
            except Exception as api_error:
                error_msg = f"API Error: {str(api_error)}"
                print(error_msg)
                print(f"Full error details: {traceback.format_exc()}")
                return {"success": False, "message": error_msg}
            
    except Exception as e:
        error_msg = f"Error during transcription: {str(e)}"
        print(error_msg)
        print(f"Full error details: {traceback.format_exc()}")
        return {"success": False, "message": error_msg}

def process_youtube_video(url, language="it"):
    """
    Download YouTube video and transcribe it with Whisper
    """
    try:
        print(f"Processing YouTube video: {url}")
        
        # First download the video
        download_result = download_youtube_video(url)
        
        if not download_result["success"]:
            return download_result["message"], ""
        
        # Then transcribe it
        transcribe_result = transcribe_with_whisper(download_result["path"], language)
        
        if not transcribe_result["success"]:
            return transcribe_result["message"], ""
        
        return transcribe_result["message"], transcribe_result["transcript"]
        
    except Exception as e:
        error_msg = f"Unexpected error in process_youtube_video: {str(e)}"
        print(error_msg)
        print(f"Full error details: {traceback.format_exc()}")
        return error_msg, ""

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
        
        # Debug output
        with gr.Accordion("Debug Information", open=True):
            debug_info = gr.Textbox(
                label="Error Details",
                value="Debug information will appear here when errors occur",
                lines=10
            )
        
        # Set up event handlers
        def process_with_debug(url, language):
            try:
                status, transcript = process_youtube_video(url, language)
                debug_text = "No errors occurred during processing."
                return status, transcript, debug_text
            except Exception as e:
                error_details = traceback.format_exc()
                print(f"Error in UI handler: {str(e)}")
                print(f"Full error details: {error_details}")
                return f"Error: {str(e)}", "", error_details
        
        process_btn.click(
            process_with_debug,
            inputs=[url_input, language_input],
            outputs=[status_output, transcript_output, debug_info],
        )
        
        return transcript_output