import os
from dotenv import load_dotenv
import traceback
import requests
import json

# Load environment variables
load_dotenv()

print("\n=== Listing Azure OpenAI Deployments ===\n")

# Get environment variables
api_key = os.getenv("OPENAI_API_KEY")
api_version = os.getenv("OPENAI_API_VERSION")
api_base = os.getenv("OPENAI_API_BASE")

print(f"API Key exists: {bool(api_key)}")
print(f"API Version: {api_version}")
print(f"API Base URL: {api_base}")

# Extract resource name from API base URL
resource_name = None
if api_base:
    # Extract resource name from URL format: https://resource-name.openai.azure.com
    parts = api_base.replace("https://", "").split(".")
    if len(parts) > 0:
        resource_name = parts[0]

print(f"Extracted resource name: {resource_name}")

# Try to list deployments using the Azure OpenAI REST API
try:
    # Format 1: Try to list deployments using the standard endpoint
    endpoint = f"{api_base}/openai/deployments?api-version={api_version}"
    print(f"\nTrying to list deployments using endpoint: {endpoint}")
    
    headers = {
        "Content-Type": "application/json",
        "api-key": api_key
    }
    
    response = requests.get(endpoint, headers=headers)
    
    if response.status_code == 200:
        result = response.json()
        print("\nAvailable deployments:")
        for deployment in result.get("data", []):
            print(f"- {deployment.get('id')} (Model: {deployment.get('model')})")
    else:
        print(f"Error: Status code {response.status_code}")
        print(f"Response: {response.text}")
        
        # Try format 2 if format 1 fails
        if response.status_code == 404:
            # Format 2: Try without 'openai' in the path
            endpoint2 = f"{api_base}/deployments?api-version={api_version}"
            print(f"\nTrying alternative endpoint: {endpoint2}")
            
            response2 = requests.get(endpoint2, headers=headers)
            
            if response2.status_code == 200:
                result2 = response2.json()
                print("\nAvailable deployments:")
                for deployment in result2.get("data", []):
                    print(f"- {deployment.get('id')} (Model: {deployment.get('model')})")
            else:
                print(f"Error: Status code {response2.status_code}")
                print(f"Response: {response2.text}")
                
except Exception as e:
    print(f"Error listing deployments: {str(e)}")
    print(f"Full error details:\n{traceback.format_exc()}")

print("\n=== Suggestions ===\n")
print("1. Check that your API key has the correct permissions to list deployments")
print("2. Verify that the API version is compatible with your Azure OpenAI resource")
print("3. Ensure that the API base URL is correctly formatted")
print("4. If you know the deployment name, try using it directly in your application")
print("5. Check the Azure portal to confirm the deployment name and status")
print("\nCommon deployment names include: gpt-35-turbo, gpt-4, text-embedding-ada-002")
