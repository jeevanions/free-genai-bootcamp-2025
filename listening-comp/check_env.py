import os
from dotenv import load_dotenv
import sys

# Load environment variables
load_dotenv()

# Check Azure OpenAI configuration
print("\n=== Checking Azure OpenAI configuration ===\n")
api_key = os.getenv("OPENAI_API_KEY")
api_version = os.getenv("OPENAI_API_VERSION")
api_base = os.getenv("OPENAI_API_BASE")
deployment_name = os.getenv("OPENAI_API_DEPLOYMENT_NAME")

# Define status symbols
OK = "✅"
WARNING = "⚠️"
ERROR = "❌"

# Check each environment variable
all_ok = True

# API Key
if api_key:
    print(f"{OK} API Key: Set")
else:
    print(f"{ERROR} API Key: Not set")
    all_ok = False

# API Version
if api_version:
    print(f"{OK} API Version: {api_version}")
else:
    print(f"{ERROR} API Version: Not set (recommended: '2023-05-15')")
    all_ok = False

# API Base URL
if api_base:
    if api_base.startswith("https://") and ".openai.azure.com" in api_base:
        print(f"{OK} API Base URL: {api_base}")
    else:
        print(f"{WARNING} API Base URL: {api_base}")
        print(f"   The URL should be in the format: https://YOUR_RESOURCE_NAME.openai.azure.com")
        all_ok = False
else:
    print(f"{ERROR} API Base URL: Not set")
    print(f"   Should be in the format: https://YOUR_RESOURCE_NAME.openai.azure.com")
    all_ok = False

# Deployment Name
if deployment_name:
    print(f"{OK} Deployment Name: {deployment_name}")
    print(f"   Note: Make sure this EXACTLY matches a deployment in your Azure OpenAI resource")
    print(f"   Common deployment names: gpt-35-turbo, gpt-4, text-embedding-ada-002")
else:
    print(f"{ERROR} Deployment Name: Not set")
    print(f"   This is the name you gave to your deployed model in the Azure OpenAI Studio")
    all_ok = False

# Summary
print("\n=== Summary ===\n")
if all_ok:
    print(f"{OK} All Azure OpenAI environment variables are properly configured!")
    print("\nYou can now run the application with:")
    print("python frontend/app.py")
    print("\nIf you still get a 404 error, run the check_azure_deployments.py script:")
    print("python check_azure_deployments.py")
else:
    print(f"{WARNING} Some environment variables are missing or incorrectly formatted.")
    print("\nPlease update your .env file with the correct values. Example:")
    print("\nOPENAI_API_KEY=your_api_key_here")
    print("OPENAI_API_VERSION=2023-05-15")
    print("OPENAI_API_BASE=https://your_resource_name.openai.azure.com")
    print("OPENAI_API_DEPLOYMENT_NAME=your_deployment_name")
    print("\nIMPORTANT: The OPENAI_API_DEPLOYMENT_NAME must match EXACTLY with a deployment")
    print("name in your Azure OpenAI resource. Common names include: gpt-35-turbo, gpt-4, etc.")
    print("\nTo check your available deployments, run:")
    print("python check_azure_deployments.py")
    print("\nNote: The application will still run with mock responses for testing.")
