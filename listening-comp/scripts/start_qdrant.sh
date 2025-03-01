#!/bin/bash

# Script to start Qdrant service using Docker Compose

# Change to the project root directory
cd "$(dirname "$0")/.."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed. Please install Docker first."
    echo "Visit https://docs.docker.com/get-docker/ for installation instructions."
    exit 1
fi

# Check if Docker Compose is installed
if ! docker compose version &> /dev/null; then
    echo "Error: Docker Compose is not installed. Please install Docker Compose first."
    echo "Visit https://docs.docker.com/compose/install/ for installation instructions."
    exit 1
fi

# Start Qdrant service
echo "Starting Qdrant service..."
docker compose up -d qdrant

# Wait for Qdrant to be ready
echo "Waiting for Qdrant to be ready..."
for i in {1..30}; do
    # Try health endpoint first (newer versions)
    if curl -s http://localhost:6333/health &> /dev/null; then
        echo "Qdrant is ready! (health endpoint)"
        echo "Qdrant UI is available at: http://localhost:6333/dashboard"
        exit 0
    fi
    
    # Fall back to collections endpoint (always available if Qdrant is running)
    if curl -s http://localhost:6333/collections &> /dev/null; then
        echo "Qdrant is ready! (collections endpoint)"
        echo "Qdrant UI is available at: http://localhost:6333/dashboard"
        exit 0
    fi
    
    echo "Waiting for Qdrant to start... ($i/30)"
    sleep 2
done

echo "Timed out waiting for Qdrant to start. Check Docker logs with 'docker compose logs qdrant'"
exit 1
