import os
import re
from dotenv import load_dotenv

# First load the current environment variables
load_dotenv()

print("\n=== Updating .env file with correct API version ===\n")

# Check if .env file exists
env_path = os.path.join(os.getcwd(), '.env')
if not os.path.exists(env_path):
    print("Error: .env file not found. Please create one first.")
    exit(1)

# Read the current .env file
with open(env_path, 'r') as file:
    env_content = file.read()

# Get the current API version
current_api_version = os.getenv("OPENAI_API_VERSION")
print(f"Current API version: {current_api_version}")

# Set the correct API version
correct_api_version = "2023-05-15"
print(f"Setting API version to: {correct_api_version}")

# Update the API version in the .env file
if current_api_version != correct_api_version:
    # Check if the API version line exists
    api_version_pattern = re.compile(r'OPENAI_API_VERSION=.*')
    if api_version_pattern.search(env_content):
        # Replace the existing API version
        updated_content = api_version_pattern.sub(f'OPENAI_API_VERSION={correct_api_version}', env_content)
        
        # Write the updated content back to the .env file
        with open(env_path, 'w') as file:
            file.write(updated_content)
        
        print(f"Updated API version in .env file to {correct_api_version}")
    else:
        # Add the API version if it doesn't exist
        with open(env_path, 'a') as file:
            file.write(f'\nOPENAI_API_VERSION={correct_api_version}\n')
        
        print(f"Added API version to .env file: {correct_api_version}")
else:
    print("API version is already set correctly.")

print("\nDone! Now run the application with:")
print("python frontend/app.py")
