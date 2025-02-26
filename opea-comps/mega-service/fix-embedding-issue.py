#!/usr/bin/env python

import requests
import os
import json

# Get environment variables
EMBEDDING_SERVICE_HOST_IP = os.getenv("EMBEDDING_SERVICE_HOST_IP", "0.0.0.0")
EMBEDDING_SERVICE_PORT = int(os.getenv("EMBEDDING_SERVICE_PORT", 8007))
HOST_IP = os.getenv("HOST_IP", "0.0.0.0")
endpoint = f"http://{EMBEDDING_SERVICE_HOST_IP}:{EMBEDDING_SERVICE_PORT}/v1/embeddings"

# Based on the curl test in initial-setup.sh and insights from chatapp.py
correct_format = {"input": "What is Deep Learning?"}

print(f"\n1. Testing with the format from initial-setup.sh: {correct_format}")
try:
    response = requests.post(
        endpoint,
        json=correct_format,
        headers={"Content-Type": "application/json"}
    )
    
    print(f"Status code: {response.status_code}")
    print(f"Content type: {response.headers.get('Content-Type')}")
    
    if response.status_code == 200:
        json_response = response.json()
        print("Success!")
        print(f"Response type: {type(json_response)}")
        print(f"Response first few entries: {str(json_response)[:100]}...")
    else:
        print(f"Error: {response.text}")
except Exception as e:
    print(f"Exception: {e}")

# Now check what chatapp.py is actually sending
print("\n2. Checking environment variables:")
print(f"  EMBEDDING_SERVICE_HOST_IP: {EMBEDDING_SERVICE_HOST_IP}")
print(f"  EMBEDDING_SERVICE_PORT: {EMBEDDING_SERVICE_PORT}")
print(f"  HOST_IP: {HOST_IP}")
print(f"  Full endpoint: {endpoint}")

print("\n3. Recommendations:")
print("- Make sure the embedding service is running")
print("- Check all environment variables are set correctly")
print("- Verify in docker logs if the embedding service is receiving requests")
print("- Run 'docker logs tei-embedding-service' to see any errors")
