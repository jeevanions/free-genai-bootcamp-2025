#!/bin/bash

# Script to stop Qdrant service using Docker Compose

# Change to the project root directory
cd "$(dirname "$0")/.."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed. Please install Docker first."
    echo "Visit https://docs.docker.com/get-docker/ for installation instructions."
    exit 1
fi

# Stop Qdrant service
echo "Stopping Qdrant service..."
docker compose down

echo "Qdrant service stopped."
