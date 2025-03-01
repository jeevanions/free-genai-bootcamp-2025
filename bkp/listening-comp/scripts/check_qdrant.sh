#!/bin/bash

# Script to check if Qdrant is running

# Change to the project root directory
cd "$(dirname "$0")/.."

# Get Qdrant URL from .env file if it exists
QDRANT_URL="http://localhost:6333"
if [ -f .env ]; then
    # Extract QDRANT_URL from .env file
    ENV_QDRANT_URL=$(grep "QDRANT_URL" .env | cut -d '=' -f2)
    if [ ! -z "$ENV_QDRANT_URL" ]; then
        QDRANT_URL=$ENV_QDRANT_URL
    fi
fi

echo "Checking Qdrant status at $QDRANT_URL..."

# Check if Qdrant is running
QDRANT_RUNNING=false

# Try health endpoint first (newer versions)
if curl -s "$QDRANT_URL/health" > /dev/null; then
    QDRANT_RUNNING=true
    echo "✅ Qdrant is running and ready! (health endpoint)"
fi

# If health endpoint failed, try collections endpoint
if [ "$QDRANT_RUNNING" = false ] && curl -s "$QDRANT_URL/collections" > /dev/null; then
    QDRANT_RUNNING=true
    echo "✅ Qdrant is running and ready! (collections endpoint)"
fi

if [ "$QDRANT_RUNNING" = true ]; then
    echo "Dashboard available at: $QDRANT_URL/dashboard"
    
    # Get collections
    COLLECTIONS=$(curl -s "$QDRANT_URL/collections")
    if [ $? -eq 0 ]; then
        echo "Collections:"
        echo "$COLLECTIONS" | grep -o '"name":"[^"]*"' | cut -d '"' -f 4
    fi
    
    exit 0
else
    echo "❌ Qdrant is not running or not accessible at $QDRANT_URL"
    
    # Check if Docker is running
    if ! docker ps > /dev/null 2>&1; then
        echo "Docker is not running. Please start Docker first."
        exit 1
    fi
    
    # Check if Qdrant container exists
    if docker ps -a | grep -q qdrant; then
        echo "Qdrant container exists but is not running. You can start it with:"
        echo "./scripts/start_qdrant.sh"
    else
        echo "Qdrant container does not exist. You can create and start it with:"
        echo "./scripts/start_qdrant.sh"
    fi
    
    exit 1
fi
