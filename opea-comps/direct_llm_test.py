#!/usr/bin/env python3
import os
import json
import requests
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

LLM_SERVER_HOST_IP = os.getenv("LLM_SERVICE_HOST_IP", "0.0.0.0")
LLM_SERVER_PORT = os.getenv("LLM_SERVICE_PORT", "8008")
LLM_MODEL = os.getenv("LLM_MODEL_ID", "llama3")

url = f"http://{LLM_SERVER_HOST_IP}:{LLM_SERVER_PORT}/v1/chat/completions"

# Prepare the request
payload = {
    "model": LLM_MODEL,
    "messages": [{"role": "user", "content": "What is responsible AI?"}],
    "max_tokens": 100,
    "stream": False
}

print(f"Sending direct request to LLM service at {url}")
print(f"Payload: {json.dumps(payload, indent=2)}")

try:
    response = requests.post(url, json=payload, timeout=10)
    print(f"Status code: {response.status_code}")
    
    if response.status_code == 200:
        result = response.json()
        print("\nResponse:")
        print(json.dumps(result, indent=2))
    else:
        print(f"Error: {response.text}")

except Exception as e:
    print(f"Request failed: {str(e)}")
