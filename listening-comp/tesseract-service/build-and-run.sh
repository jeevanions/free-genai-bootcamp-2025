#!/bin/bash

# Build the Docker image
echo "Building Docker image..."
docker-compose build

# Check if build was successful
if [ $? -eq 0 ]; then
    echo "Build successful! Starting container..."
    docker-compose up -d
    echo "Container started. Service available at http://localhost:8000"
else
    echo "Build failed. Please check the error messages above."
fi
