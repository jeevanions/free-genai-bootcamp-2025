"""
Text-to-Speech plugin using gTTS.
"""
import os
from typing import Optional
from pathlib import Path
from gtts import gTTS
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

class TTSPlugin:
    """Plugin for Text-to-Speech conversion."""
    
    def __init__(self):
        """Initialize the TTS plugin."""
        self.output_dir = Path("/tmp/tts_output")
        self.output_dir.mkdir(parents=True, exist_ok=True)
    
    def text_to_speech(self, text: str, language: str = "it", filename: Optional[str] = None) -> str:
        """
        Convert text to speech.
        
        Args:
            text: Text to convert to speech
            language: Language code (default: Italian)
            filename: Optional filename for the output file
            
        Returns:
            Path to the generated audio file
        """
        try:
            # Generate a filename if not provided
            if filename is None:
                filename = f"speech_{hash(text) % 10000}.mp3"
            
            # Ensure filename has .mp3 extension
            if not filename.endswith(".mp3"):
                filename += ".mp3"
            
            # Full path to output file
            output_path = self.output_dir / filename
            
            # Generate speech
            tts = gTTS(text=text, lang=language, slow=False)
            tts.save(str(output_path))
            
            return str(output_path)
        
        except Exception as e:
            print(f"Error generating speech: {e}")
            return ""
