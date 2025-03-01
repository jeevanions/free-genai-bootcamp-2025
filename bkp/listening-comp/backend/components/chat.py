"""
Chat component for GPT-4 integration.
"""
import os
from typing import List, Dict, Any
from openai import AzureOpenAI
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

class ChatComponent:
    """Component for interacting with Azure OpenAI's GPT-4 model."""
    
    def __init__(self):
        """Initialize the chat component with Azure OpenAI credentials."""
        # Configure Azure OpenAI client
        self.client = AzureOpenAI(
            api_key=os.getenv("AZURE_OPENAI_API_KEY"),
            api_version="2023-05-15",  # Update this to the latest version if needed
            azure_endpoint=os.getenv("AZURE_OPENAI_ENDPOINT")
        )
        self.deployment_name = os.getenv("AZURE_OPENAI_DEPLOYMENT_NAME")
        
    def chat(self, messages: List[Dict[str, str]]) -> str:
        """
        Send a conversation to the GPT-4 model and get a response.
        
        Args:
            messages: List of message dictionaries with 'role' and 'content' keys
                     Example: [{"role": "user", "content": "Hello!"}]
        
        Returns:
            The model's response text
        """
        try:
            response = self.client.chat.completions.create(
                model=self.deployment_name,
                messages=messages,
                temperature=0.7,
                max_tokens=800,
                top_p=0.95,
                frequency_penalty=0,
                presence_penalty=0,
                stop=None
            )
            return response.choices[0].message.content
        except Exception as e:
            return f"Error communicating with OpenAI: {str(e)}"
    
    def clear_conversation(self) -> List[Dict[str, str]]:
        """
        Clear the conversation history.
        
        Returns:
            An empty conversation list with system prompt
        """
        return [{"role": "system", "content": "You are a helpful language learning assistant specializing in Italian."}]
