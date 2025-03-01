import cv2
import pytesseract
import numpy as np

def preprocess_image(image):
    """
    Preprocess image for OCR
    """
    # Convert to grayscale
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    
    # Apply thresholding
    _, thresh = cv2.threshold(gray, 150, 255, cv2.THRESH_BINARY_INV)
    
    # Apply noise reduction
    denoised = cv2.fastNlMeansDenoising(thresh, None, 10, 7, 21)
    
    return denoised

def extract_text_from_image(image, language="ita"):
    """
    Extract text from image using pytesseract
    """
    # Configure pytesseract
    config = f"-l {language} --oem 1 --psm 3"
    
    # Preprocess image
    processed_image = preprocess_image(image)
    
    # Perform OCR
    text = pytesseract.image_to_string(processed_image, config=config)
    
    return text

def extract_frames_from_video(video_path, frame_interval=1):
    """
    Extract frames from video at specified intervals
    """
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
    
    return frames