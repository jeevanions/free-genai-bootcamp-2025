import os
from dotenv import load_dotenv
import traceback
import json
import requests
from openai import AzureOpenAI

# Load environment variables
load_dotenv()

print("\n=== Checking Azure OpenAI Deployments ===\n")

# Print all environment variables for debugging (without the actual API key value)
api_key = os.getenv("OPENAI_API_KEY")
api_version = os.getenv("OPENAI_API_VERSION")
api_base = os.getenv("OPENAI_API_BASE")
deployment_name = os.getenv("OPENAI_API_DEPLOYMENT_NAME")

print(f"API Key exists: {bool(api_key)}")
print(f"API Version: {api_version}")
print(f"API Base URL: {api_base}")
print(f"Deployment Name: {deployment_name}")

# Try multiple approaches to connect to Azure OpenAI

# Define different API versions to try
api_versions_to_try = [
    api_version,  # Try the one from the environment first
    "2023-05-15",
    "2023-07-01-preview",
    "2023-12-01-preview"
]

# Remove empty or None values
api_versions_to_try = [v for v in api_versions_to_try if v]

# Make sure we have at least one version to try
if not api_versions_to_try:
    api_versions_to_try = ["2023-05-15"]

# Remove duplicates while preserving order
seen = set()
api_versions_to_try = [x for x in api_versions_to_try if not (x in seen or seen.add(x))]

print(f"\nWill try the following API versions: {api_versions_to_try}\n")

# Approach 1: Using the OpenAI SDK
print("\n=== Approach 1: Using the OpenAI SDK ===\n")

for version in api_versions_to_try:
    print(f"\nTrying with API version: {version}")
    try:
        # Initialize Azure OpenAI client
        print("Initializing Azure OpenAI client...")
        client = AzureOpenAI(
            api_key=api_key,
            api_version=version,
            azure_endpoint=api_base,
        )
    
        # Try a test completion with the specified deployment
        print(f"Testing connection with deployment name: {deployment_name}")
        
        try:
            test_response = client.chat.completions.create(
                model=deployment_name,
                messages=[{"role": "user", "content": "Hello, this is a test."}],
                max_tokens=10
            )
            print(f"Success! Response: {test_response.choices[0].message.content}")
            print(f"\n*** SUCCESS with API version {version} ***\n")
            break  # Exit the loop if successful
        except Exception as e:
            print(f"Error with deployment '{deployment_name}' using API version {version}: {str(e)}")
            # Only print full error details for the first version to avoid clutter
            if version == api_versions_to_try[0]:
                print(f"Full error details:\n{traceback.format_exc()}")
            
    except Exception as e:
        print(f"Error initializing client with API version {version}: {str(e)}")
        # Only print full error details for the first version to avoid clutter
        if version == api_versions_to_try[0]:
            print(f"Full error details:\n{traceback.format_exc()}")

# Approach 2: Using direct REST API calls
print("\n=== Approach 2: Using direct REST API calls ===\n")

# Try each API version
for version in api_versions_to_try:
    print(f"\nTrying direct API calls with API version: {version}")
    
    try:
        # Format 1: Standard Azure OpenAI endpoint format
        endpoint = f"{api_base}/openai/deployments/{deployment_name}/chat/completions?api-version={version}"
        print(f"Trying endpoint format 1: {endpoint}")
        
        headers = {
            "Content-Type": "application/json",
            "api-key": api_key
        }
        
        payload = {
            "messages": [{"role": "user", "content": "Hello, this is a test."}],
            "max_tokens": 10
        }
        
        response = requests.post(endpoint, headers=headers, json=payload)
        
        if response.status_code == 200:
            result = response.json()
            print(f"Success! Response: {result['choices'][0]['message']['content']}")
            print(f"\n*** SUCCESS with API version {version} and endpoint format 1 ***\n")
            break  # Exit the loop if successful
        else:
            print(f"Error: Status code {response.status_code}")
            print(f"Response: {response.text}")
            
            # Try format 2 if format 1 fails
            if response.status_code == 404:
                # Format 2: Alternative endpoint format without 'openai'
                endpoint2 = f"{api_base}/deployments/{deployment_name}/chat/completions?api-version={version}"
                print(f"\nTrying endpoint format 2: {endpoint2}")
                
                response2 = requests.post(endpoint2, headers=headers, json=payload)
                
                if response2.status_code == 200:
                    result2 = response2.json()
                    print(f"Success! Response: {result2['choices'][0]['message']['content']}")
                    print(f"\n*** SUCCESS with API version {version} and endpoint format 2 ***\n")
                    break  # Exit the loop if successful
                else:
                    print(f"Error: Status code {response2.status_code}")
                    print(f"Response: {response2.text}")
                    
                    # Try format 3 if format 2 fails
                    if response2.status_code == 404:
                        # Format 3: Try with a trailing slash in the base URL
                        if not api_base.endswith("/"):
                            api_base_with_slash = f"{api_base}/"
                            endpoint3 = f"{api_base_with_slash}openai/deployments/{deployment_name}/chat/completions?api-version={version}"
                            print(f"\nTrying endpoint format 3: {endpoint3}")
                            
                            response3 = requests.post(endpoint3, headers=headers, json=payload)
                            
                            if response3.status_code == 200:
                                result3 = response3.json()
                                print(f"Success! Response: {result3['choices'][0]['message']['content']}")
                                print(f"\n*** SUCCESS with API version {version} and endpoint format 3 ***\n")
                                break  # Exit the loop if successful
                            else:
                                print(f"Error: Status code {response3.status_code}")
                                print(f"Response: {response3.text}")
    except Exception as e:
        print(f"Error with API version {version}: {e}")

# Provide suggestions based on the results
print("\n=== Suggestions ===\n")
print("If all approaches failed, check the following:")
print("1. Verify that the deployment name is correct and exists in your Azure OpenAI resource")
print("2. Confirm that your API key has the correct permissions")
print("3. Check that the API version is compatible with your Azure OpenAI resource")
print("4. Ensure that the API base URL is correctly formatted")
print("\nIf one of the approaches succeeded, update your code to use that approach.")
