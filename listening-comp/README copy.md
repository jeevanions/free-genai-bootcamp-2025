# Language Listening Comprehension App

An AI-powered language learning application that generates listening comprehension exercises from YouTube content.

## Setup

1. Install dependencies:
```bash
uv venv
source .venv/bin/activate
uv pip install -e .
```

2. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your API keys and configurations
```

3. Start Qdrant service (required for RAG functionality):
```bash
# Make sure Docker is installed and running
./scripts/start_qdrant.sh
```

See [Qdrant Setup Guide](docs/qdrant_setup.md) for detailed instructions.

4. Run the application:
```bash
python -m frontend.ui.app
```

## Project Structure

```
listening-comp/
├── backend/
│   ├── components/     # OPEA components
│   │   ├── chat.py    # GPT-4 chat integration
│   │   ├── rag.py     # RAG implementation
│   │   └── media.py   # Media processing
│   ├── plugins/       # External integrations
│   │   ├── youtube.py # YouTube processing
│   │   ├── ocr.py     # OCR processing
│   │   └── tts.py     # Text-to-Speech
│   └── services/      # Business logic
│       ├── transcription.py
│       └── knowledge_base.py
├── frontend/
│   └── ui/           # Gradio interfaces
│       ├── app.py    # Main application
│       └── components/ # UI components
└── output/           # Generated content
    ├── transcript-yt/  # YouTube transcripts
    ├── transcript-vid/ # Whisper transcripts
    └── transcript-ocr/ # OCR output
```

## Features

1. Chat with GPT-4 model
2. YouTube transcript downloading
3. Audio transcription with Whisper
4. OCR extraction from video frames
5. RAG-based knowledge retrieval with Qdrant vector database
6. Interactive learning interface

## Development

- Run tests: `pytest`
- Format code: `black .`
- Check types: `mypy .`

## Docker Services

### Qdrant Vector Database

The application uses Qdrant for vector storage and retrieval in the RAG component. Docker is required to run Qdrant.

- Start Qdrant: `./scripts/start_qdrant.sh`
- Stop Qdrant: `./scripts/stop_qdrant.sh`
- Access Qdrant dashboard: http://localhost:6333/dashboard

The Docker configuration is defined in `docker-compose.yml`. Qdrant data is persisted in a Docker volume.

For detailed setup instructions, troubleshooting, and more information, see the [Qdrant Setup Guide](docs/qdrant_setup.md).

## License

MIT
