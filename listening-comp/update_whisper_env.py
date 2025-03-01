import os
import re
from dotenv import load_dotenv

# First load the current environment variables
load_dotenv()

print("\n=== Updating .env file with Whisper deployment configuration ===\n")

# Check if .env file exists
env_path = os.path.join(os.getcwd(), '.env')
if not os.path.exists(env_path):
    print("Error: .env file not found. Please create one first.")
    exit(1)

# Read the current .env file
with open(env_path, 'r') as file:
    env_content = file.read()

# Get the current Whisper deployment name
current_whisper_deployment = os.getenv("OPENAI_WHISPER_DEPLOYMENT_NAME")
print(f"Current Whisper deployment name: {current_whisper_deployment or 'Not set'}")

# Check if the user wants to update it
print("\nIMPORTANT: In Azure OpenAI, you need to create a deployment of the whisper-1 model")
print("and give it a name. The default name 'whisper' might not exist in your Azure resource.")
print("\nTo fix this:")
print("1. Go to Azure OpenAI Studio")
print("2. Click on 'Deployments'")
print("3. Create a new deployment using the whisper-1 model")
print("4. Give it a name (e.g., 'whisper')")
print("5. Update your .env file with the deployment name")

# Check if the OPENAI_WHISPER_DEPLOYMENT_NAME variable exists in the .env file
whisper_pattern = re.compile(r'OPENAI_WHISPER_DEPLOYMENT_NAME=.*')
if whisper_pattern.search(env_content):
    print("\nThe OPENAI_WHISPER_DEPLOYMENT_NAME variable already exists in your .env file.")
    print("You can update it manually by editing the .env file.")
else:
    # Add the OPENAI_WHISPER_DEPLOYMENT_NAME variable to the .env file
    print("\nAdding OPENAI_WHISPER_DEPLOYMENT_NAME variable to your .env file...")
    with open(env_path, 'a') as file:
        file.write("\n# The deployment name for Whisper transcription (required for transcription features)\n")
        file.write("# This is the name you gave to your whisper-1 model deployment in Azure OpenAI Studio\n")
        file.write("OPENAI_WHISPER_DEPLOYMENT_NAME=whisper\n")
    
    print("Added OPENAI_WHISPER_DEPLOYMENT_NAME=whisper to your .env file.")
    print("If your deployment has a different name, please update it manually.")

print("\nDone! Now you need to create a deployment of the whisper-1 model in Azure OpenAI Studio.")
print("Make sure the deployment name matches the OPENAI_WHISPER_DEPLOYMENT_NAME in your .env file.")
