FROM python:3.11-slim

# Install Tesseract OCR and required dependencies
RUN apt-get update && apt-get install -y \
    tesseract-ocr \
    libtesseract-dev \
    tesseract-ocr-ita \
    tesseract-ocr-eng \
    tesseract-ocr-fra \
    tesseract-ocr-spa \
    tesseract-ocr-deu \
    ffmpeg \
    libsm6 \
    libxext6 \
    libgl1 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements.txt
COPY requirements.txt .

# Install UV package manager
RUN pip install uv

# Install dependencies using UV
RUN uv pip install --system -r requirements.txt

# Create directories for files
RUN mkdir -p /app/output/transcript-vid /app/output/transcript-ocr

# Set working directory
WORKDIR /app

# Copy application code
COPY app.py .

# Expose the port the app will run on
EXPOSE 8000

# Command to run the application
CMD ["python", "app.py"]