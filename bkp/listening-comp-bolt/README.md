# Italian Language Listening Comprehension App

A comprehensive platform for learning Italian through listening comprehension exercises, powered by AI.

## Features

- **YouTube Transcript Extraction**: Download transcripts from Italian YouTube videos
- **Whisper Transcription**: Generate accurate transcripts using OpenAI's Whisper model
- **OCR Extraction**: Extract text from video frames using specialized Italian OCR
- **RAG Implementation**: Use retrieval-augmented generation for contextual learning
- **Interactive Learning**: Practice with AI-generated exercises tailored to your level

## Technical Stack

- **Frontend**: React with Tailwind CSS
- **Backend**: Python with Gradio
- **Vector Store**: Qdrant
- **LLM**: GPT-4o via Azure OpenAI Services
- **OCR**: Tesseract with Italian language support
- **Audio Transcription**: Whisper
- **YouTube Integration**: yt-dlp
- **Text-to-Speech**: gTTS (Google Text-to-Speech)
- **Package Management**: uv (Python)

## Getting Started

### Prerequisites

- Node.js and npm
- Python 3.8+
- OpenAI API key (Azure)
- Qdrant (local or cloud)

### Installation

1. Clone the repository
2. Install frontend dependencies:
   ```
   npm install
   ```
3. Install backend dependencies:
   ```
   uv pip install -r requirements.txt
   ```
4. Copy `.env.example` to `.env` and fill in your API keys

### Running the Application

1. Start the frontend:
   ```
   npm run dev
   ```
2. Start the backend:
   ```
   npm run start-backend
   ```
3. Open your browser and navigate to the URL shown in the console

## Project Structure

- `/src` - React frontend
- `/backend` - Python backend
  - `/components` - Gradio interface components
  - `/utils` - Utility functions for YouTube, OCR, and vector store operations
- `/output` - Output directory for transcripts and processed data
  - `/transcript-yt` - YouTube transcripts
  - `/transcript-vid` - Whisper transcripts
  - `/transcript-ocr` - OCR extracted text

## Usage

1. **YouTube Transcript**: Enter a YouTube URL to download its transcript
2. **Whisper Transcript**: Enter a YouTube URL to download the video and generate a transcript
3. **OCR Extraction**: Extract text from video frames
4. **RAG Implementation**: Add documents to the vector store and query them
5. **Interactive Learning**: Generate and practice with interactive exercises

## License

This project is licensed under the MIT License - see the LICENSE file for details.