#!/usr/bin/env python

import requests
import os
import json

EMBEDDING_SERVICE_HOST_IP = os.getenv("EMBEDDING_SERVICE_HOST_IP", "0.0.0.0")
EMBEDDING_SERVICE_PORT = int(os.getenv("EMBEDDING_SERVICE_PORT", 8007))

# Test with different formats
formats = [
    {"input": "What is Deep Learning?"}, # HF format
    {"inputs": "What is Deep Learning?"}, # Alternative format
    {"text": "What is Deep Learning?"} # Original format
]

for fmt in formats:
    print(f"\nTesting with format: {fmt}")
    try:
        response = requests.post(
            f"http://{EMBEDDING_SERVICE_HOST_IP}:{EMBEDDING_SERVICE_PORT}/v1/embeddings",
            json=fmt,
            headers={"Content-Type": "application/json"}
        )
        
        print(f"Status code: {response.status_code}")
        print(f"Content type: {response.headers.get('Content-Type')}")
        
        if response.status_code == 200:
            print("Success!")
            print(f"Response: {json.dumps(response.json(), indent=2)[:100]}...")
        else:
            print(f"Error: {response.text}")
    except Exception as e:
        print(f"Exception: {e}")
