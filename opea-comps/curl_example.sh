#!/bin/bash

# Load environment variables
source .env

echo "Sending request to chatqna endpoint..."

# Format correctly following the successful direct call pattern
curl "http://${HOST_IP}:${MEGA_SERVICE_PORT}/v1/chatqna" \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "messages": [{"role": "user", "content": "What is responsible AI?"}],
        "max_tokens": 100,
        "stream": true
    }'

echo -e "\nRequest completed."
