"""
OCR plugin for extracting text from video frames.
"""
import os
from typing import List, Dict, Any, Optional
from pathlib import Path
import json
import easyocr
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

class OCRPlugin:
    """Plugin for OCR processing of video frames."""
    
    def __init__(self):
        """Initialize the OCR plugin with configuration."""
        self.transcript_ocr_dir = Path(os.getenv("TRANSCRIPT_OCR_DIR", "output/transcript-ocr"))
        self.ocr_language = os.getenv("OCR_LANGUAGE", "it")
        
        # Create output directory if it doesn't exist
        self.transcript_ocr_dir.mkdir(parents=True, exist_ok=True)
        
        # Initialize OCR reader
        self.reader = None  # Lazy loading to save memory
    
    def _get_reader(self):
        """Lazy load the OCR reader."""
        if self.reader is None:
            # Initialize with Italian and English
            languages = [self.ocr_language]
            if self.ocr_language != "en":
                languages.append("en")
            
            self.reader = easyocr.Reader(languages)
        
        return self.reader
    
    def process_frames(self, frame_paths: List[str], video_id: str) -> Dict[str, Any]:
        """
        Process video frames with OCR.
        
        Args:
            frame_paths: List of paths to frame images
            video_id: YouTube video ID for output file naming
            
        Returns:
            Dictionary with OCR results
        """
        reader = self._get_reader()
        
        # Process each frame
        results = []
        combined_text = ""
        
        for i, frame_path in enumerate(frame_paths):
            try:
                # Run OCR on the frame
                ocr_result = reader.readtext(frame_path)
                
                # Extract text
                frame_text = " ".join([item[1] for item in ocr_result])
                
                # Add to results
                results.append({
                    "frame_index": i,
                    "frame_path": frame_path,
                    "text": frame_text,
                    "confidence": [item[2] for item in ocr_result],
                    "boxes": [item[0] for item in ocr_result]
                })
                
                # Add to combined text if not empty
                if frame_text.strip():
                    combined_text += frame_text + "\n\n"
            
            except Exception as e:
                print(f"Error processing frame {frame_path}: {e}")
        
        # Save combined text
        output_file = self.transcript_ocr_dir / f"{video_id}.txt"
        with open(output_file, "w", encoding="utf-8") as f:
            f.write(combined_text)
        
        # Save detailed results
        detailed_output_file = self.transcript_ocr_dir / f"{video_id}_detailed.json"
        with open(detailed_output_file, "w", encoding="utf-8") as f:
            # Convert non-serializable types to strings
            serializable_results = []
            for result in results:
                serializable_result = {
                    "frame_index": result["frame_index"],
                    "frame_path": result["frame_path"],
                    "text": result["text"],
                    "confidence": result["confidence"],
                    "boxes": [[str(point) for point in box] for box in result["boxes"]]
                }
                serializable_results.append(serializable_result)
            
            json.dump(serializable_results, f, ensure_ascii=False, indent=2)
        
        return {
            "video_id": video_id,
            "combined_text": combined_text,
            "transcript_file": str(output_file),
            "detailed_file": str(detailed_output_file),
            "frame_count": len(frame_paths),
            "processed_count": len(results)
        }
