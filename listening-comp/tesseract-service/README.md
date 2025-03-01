# Tesseract OCR Service

A FastAPI-based service that provides OCR capabilities using Tesseract, packaged in a Docker container.

## Features

- Image OCR: Extract text from uploaded images
- YouTube Video OCR: Extract text from YouTube videos by processing frames
- Multiple language support: English, Italian, French, Spanish, German

## Technology Stack

- Python 3.11
- FastAPI
- Tesseract OCR
- UV Package Manager (for dependency management)
- Docker

## Setup and Installation

### Prerequisites

- Docker and Docker Compose

### Running the Service

1. Clone the repository
2. Navigate to the project directory
3. Run the service using Docker Compose:

```bash
docker-compose up -d
```

The service will be available at http://localhost:8000

## API Endpoints

- `GET /`: Check if the service is running
- `GET /version`: Get Tesseract version information
- `POST /ocr/image`: Perform OCR on an uploaded image
- `POST /ocr/youtube`: Process a YouTube video and extract text using OCR

## Development

### Local Development Setup

1. Install UV package manager using pip:
```bash
pip install uv
```

2. Install dependencies (with virtual environment):
```bash
uv venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
uv pip install -r requirements.txt
```

Or without a virtual environment:
```bash
uv pip install --system -r requirements.txt
```

3. Run the application:
```bash
python app.py
```

## License

[MIT License](LICENSE)
