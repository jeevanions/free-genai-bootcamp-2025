# Business Goal: 
You are an Applied AI Engineer and you have been tasked to build a Language Listening Comprehension App. There are practice listening comprehension examples for language learning tests on youtube.

Pull the youtube content, and use that to generate out similar style listening comprehension.

# Technical Uncertainty:
Don’t know Italian!
Accessing or storing documents as vector store with Qdrant
TSS might not exist for my target language OR might not be good enough.
ASR might not exist for my target language OR might not be good enough.
Can you pull transcripts for the target videos?

# Technical Requirements:

- Vectorstore - Qdrant - Knowledge Base 
- LLM Model - GPT 4o Model [Using Azure OpenAI Services for model deployment] - LLM + Tool Use “Agent”
- Extract Video frames - OpenCV or ffmpeg
- OCR specifically trained for Italian on those frames - Tesseract, EasyOCR has good Italian language support and PaddleOCR is another strong multilingual OCR system
- use Whisper for Audio Transcription.
- Download youtube video - yt-dlp
- Frontend - Use gradio
- Text to Speach - gTTS (Google Text-to-Speech)
- use UV as python package manager.


# Functional requirements

## UI Design

Create a professional Italian-themed application using Gradio with the following design elements:

- **Header**: Centered header with Italian flag colors (green, white, red) and the title "Italian Language Learning Platform"
- **Layout**: Left sidebar for navigation and main content area for functionality display
- **Navigation**: Menu buttons in the left sidebar for accessing different components
- **Content Area**: Tab-based interface that changes based on the selected menu option
- **Color Scheme**: Professional design with Italian flag colors as accents
- **Footer**: Copyright information and platform details

## Navigation Structure

The application will have the following menu options in the left sidebar:

1. Chat with GPT-4 model
2. Raw Transcript using YouTube transcript downloader
3. Generate Transcript using Whisper
4. OCR Extraction using OpenCV and Tesseract
5. RAG Implementation
6. Interactive Learning App

On selecting each option from the side panel, the corresponding tab will be displayed in the main content area.

1. Chat with GPT-4 model 
   - Professional chat UI with Italian flag-themed avatar for the assistant
   - Message input area with placeholder text for Italian language queries
   - Clear chat button to reset the conversation
   - Debug information section (collapsed by default) for developers
   - Users will use this interface to try prompts created by feeding transcripts from other components
2. Raw Transcript using Youtube transcript downloader
   - Takes youtube video url as input
   - Youtube may provide transcripts for the videos. use YoutubeTranscriptDownloader
   - From the url download the transcript from Youtube and save it in the output folder under "transcript-yt"
3. Generate Transcript using Whisper
   - Takes youtube url as input
   - Download the video locally and generate the transcript using Whisper
   - The output will be stored in the output folder under "transcript-vid"
4. OCR Extraction using OpenCV and Tesseract
   - Simple UI that will extract text from the video.
   - The output will be stored in the output folder under "transcript-ocr"
5. Interactive Learning App
   - We will implement this later.

All files save shoud use the video id from the youtube url as prefix the filename
Each of the above component will be created in a backend folder can be run individually from cli  to test them out.
For frontend we can have a frontend folder that references the backend as library. 
Use python language.




