import gradio as gr
import os
import time
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Try to initialize the OpenAI client with detailed error logging
try:
    import traceback
    from openai import AzureOpenAI
    
    # Print all environment variables for debugging (without the actual API key value)
    print("\n=== Azure OpenAI Configuration ===")
    print(f"API Key exists: {bool(os.getenv('OPENAI_API_KEY'))}")
    print(f"API Version: {os.getenv('OPENAI_API_VERSION')}")
    print(f"API Base URL: {os.getenv('OPENAI_API_BASE')}")
    print(f"Deployment Name: {os.getenv('OPENAI_API_DEPLOYMENT_NAME')}")
    
    # Check if all required environment variables are set
    if all([os.getenv("OPENAI_API_KEY"), 
            os.getenv("OPENAI_API_VERSION"), 
            os.getenv("OPENAI_API_BASE"),
            os.getenv("OPENAI_API_DEPLOYMENT_NAME")]):
        
        # Initialize Azure OpenAI client
        print("Initializing Azure OpenAI client...")    
        # Ensure the Azure endpoint has the https:// prefix
        azure_endpoint = os.getenv("OPENAI_API_BASE")
        if azure_endpoint and not azure_endpoint.startswith("https://"):
            azure_endpoint = f"https://{azure_endpoint}"
            print(f"Added https:// prefix to endpoint: {azure_endpoint}")
            
        client = AzureOpenAI(
            api_key=os.getenv("OPENAI_API_KEY"),
            api_version=os.getenv("OPENAI_API_VERSION"),
            azure_endpoint=azure_endpoint,
        )
        
        # Test the client with a simple completion to verify it works
        try:
            print(f"Testing connection with deployment name: {os.getenv('OPENAI_API_DEPLOYMENT_NAME')}")
            deployment_name = os.getenv("OPENAI_API_DEPLOYMENT_NAME")
            
            # FIXED: Use correct path for Azure OpenAI API
            test_response = client.chat.completions.create(
                model=deployment_name,
                messages=[{"role": "user", "content": "Hello, this is a test."}],
                max_tokens=10
            )
            print("Connection test successful!")
            USE_MOCK = False
        except Exception as api_error:
            print(f"Error testing API connection: {str(api_error)}")
            print(f"Full error details: {traceback.format_exc()}")
            
            # Check if this is a deployment name issue
            if "Resource not found" in str(api_error):
                print("\nPOSSIBLE FIX: The deployment name may be incorrect.")
                print("Please verify in your Azure OpenAI Studio that:")
                print(f"1. The deployment name '{deployment_name}' exists exactly as specified")
                print("2. Your Azure OpenAI resource has the correct permissions")
                print("3. Your API key has access to this specific deployment")
            
            USE_MOCK = True
    else:
        missing_vars = []
        if not os.getenv("OPENAI_API_KEY"): missing_vars.append("OPENAI_API_KEY")
        if not os.getenv("OPENAI_API_VERSION"): missing_vars.append("OPENAI_API_VERSION")
        if not os.getenv("OPENAI_API_BASE"): missing_vars.append("OPENAI_API_BASE")
        if not os.getenv("OPENAI_API_DEPLOYMENT_NAME"): missing_vars.append("OPENAI_API_DEPLOYMENT_NAME")
        
        print(f"Warning: Missing environment variables: {', '.join(missing_vars)}. Using mock responses.")
        USE_MOCK = True
except ImportError as e:
    print(f"Warning: ImportError - {str(e)}. Using mock responses.")
    USE_MOCK = True
except Exception as e:
    print(f"Warning: Unexpected error initializing Azure OpenAI client: {str(e)}")
    print(f"Full error details: {traceback.format_exc()}")
    USE_MOCK = True

if USE_MOCK:
    print("Using mock responses for chat functionality.")

def create_chat_interface(parent):
    """
    Creates the chat interface with GPT-4 model
    """
    with parent:
        # Header for the chat section
        with gr.Group(elem_classes="chat-with-assistant"):
            gr.Markdown("## Chat with our Italian Language Assistant")
            gr.Markdown("Ask questions about Italian language, request translations, or get help with grammar.", elem_classes="chat-description")
        
        # Chatbot component
        chatbot = gr.Chatbot(
            height=400,
            bubble_full_width=False,
            show_label=False,
            type="messages",
            elem_classes="chatbot"
        )
        
        # Input area
        with gr.Row(elem_classes="chat-input-container"):
            msg = gr.Textbox(
                placeholder="Type your Italian language question here...",
                scale=9,
                container=False,
                elem_classes="chat-input"
            )
            submit_btn = gr.Button("Send", scale=1, variant="primary", elem_classes="send-btn")
        
        clear_btn = gr.Button("Clear Chat", size="sm", variant="secondary", elem_classes="clear-btn")
        
        def respond(message, chat_history):
            try:
                # Add user message to chat history
                chat_history.append({"role": "user", "content": message})
                
                if USE_MOCK:
                    # Simulate API call with a delay
                    time.sleep(1)
                    
                    # Generate mock response based on the message
                    if "hello" in message.lower() or "hi" in message.lower() or "ciao" in message.lower():
                        bot_message = "Ciao! How can I help you with Italian language learning today?"
                    elif "translate" in message.lower():
                        bot_message = "I'd be happy to help with translation! What would you like me to translate?"
                    elif "grammar" in message.lower():
                        bot_message = "Italian grammar has some interesting features. What specific aspect would you like to know about?"
                    elif "vocabulary" in message.lower():
                        bot_message = "Building vocabulary is essential for language learning. Would you like some common Italian words or phrases?"
                    else:
                        bot_message = "I'm a mock Italian language assistant. Please set up your Azure OpenAI API credentials to use the full functionality. In the meantime, you can ask me about translations, grammar, or vocabulary."
                else:
                    try:
                        # Prepare conversation history for the API
                        messages = []
                        for msg in chat_history:
                            messages.append({"role": msg["role"], "content": msg["content"]})
                        
                        print(f"Sending request to Azure OpenAI with deployment: {os.getenv('OPENAI_API_DEPLOYMENT_NAME')}")
                        # Call the OpenAI API
                        deployment_name = os.getenv("OPENAI_API_DEPLOYMENT_NAME")
                        
                        # FIXED: Make sure we're sending the correct message format
                        # The previous code had a logic error - it was excluding the current user message
                        response = client.chat.completions.create(
                            model=deployment_name,
                            messages=messages,  # Include all messages including the current one
                            temperature=0.7,
                            max_tokens=800,
                        )
                        
                        # Extract the response text
                        bot_message = response.choices[0].message.content
                        print("Successfully received response from Azure OpenAI")
                    except Exception as api_error:
                        import traceback
                        print(f"Error calling Azure OpenAI API: {str(api_error)}")
                        print(f"Full error details: {traceback.format_exc()}")
                        bot_message = f"Error: {str(api_error)}. Please check the console logs for details."
                
                # Add assistant response to chat history
                chat_history.append({"role": "assistant", "content": bot_message})
                
                return "", chat_history
            except Exception as e:
                import traceback
                error_message = f"Error: {str(e)}. Please check your Azure OpenAI API configuration."
                print(f"Unexpected error in respond function: {str(e)}")
                print(f"Full error details: {traceback.format_exc()}")
                chat_history.append({"role": "assistant", "content": error_message})
                return "", chat_history
        
        # Set up event handlers
        submit_btn.click(respond, [msg, chatbot], [msg, chatbot])
        msg.submit(respond, [msg, chatbot], [msg, chatbot])
        clear_btn.click(lambda: [], None, chatbot, queue=False)
        
        return chatbot


def chat_with_gpt(message, chat_history=None):
    """
    Function to chat with GPT-4 without the UI
    This can be used for CLI testing or programmatic access
    
    Args:
        message (str): The user's message
        chat_history (list, optional): Previous chat history. Defaults to None.
        
    Returns:
        tuple: Updated chat history and the assistant's response
    """
    if chat_history is None:
        chat_history = []
    
    try:
        # Add user message to chat history
        chat_history.append({"role": "user", "content": message})
        
        if USE_MOCK:
            # Simulate API call with a delay
            time.sleep(1)
            
            # Generate mock response based on the message
            if "hello" in message.lower() or "hi" in message.lower() or "ciao" in message.lower():
                bot_message = "Ciao! How can I help you with Italian language learning today?"
            elif "translate" in message.lower():
                bot_message = "I'd be happy to help with translation! What would you like me to translate?"
            elif "grammar" in message.lower():
                bot_message = "Italian grammar has some interesting features. What specific aspect would you like to know about?"
            elif "vocabulary" in message.lower():
                bot_message = "Building vocabulary is essential for language learning. Would you like some common Italian words or phrases?"
            else:
                bot_message = "I'm a mock Italian language assistant. Please set up your Azure OpenAI API credentials to use the full functionality. In the meantime, you can ask me about translations, grammar, or vocabulary."
        else:
            try:
                # Prepare conversation history for the API
                messages = []
                for msg in chat_history:
                    messages.append({"role": msg["role"], "content": msg["content"]})
                
                print(f"Sending request to Azure OpenAI with deployment: {os.getenv('OPENAI_API_DEPLOYMENT_NAME')}")
                # Call the OpenAI API
                deployment_name = os.getenv("OPENAI_API_DEPLOYMENT_NAME")
                
                response = client.chat.completions.create(
                    model=deployment_name,
                    messages=messages,
                    temperature=0.7,
                    max_tokens=800,
                )
                
                # Extract the response text
                bot_message = response.choices[0].message.content
                print("Successfully received response from Azure OpenAI")
            except Exception as api_error:
                import traceback
                print(f"Error calling Azure OpenAI API: {str(api_error)}")
                print(f"Full error details: {traceback.format_exc()}")
                bot_message = f"Error: {str(api_error)}. Please check the console logs for details."
        
        # Add assistant response to chat history
        chat_history.append({"role": "assistant", "content": bot_message})
        
        return chat_history, bot_message
    except Exception as e:
        import traceback
        error_message = f"Error: {str(e)}. Please check your Azure OpenAI API configuration."
        print(f"Unexpected error in chat_with_gpt function: {str(e)}")
        print(f"Full error details: {traceback.format_exc()}")
        chat_history.append({"role": "assistant", "content": error_message})
        return chat_history, error_message


if __name__ == "__main__":
    import argparse
    
    parser = argparse.ArgumentParser(description="Italian Language Assistant CLI")
    parser.add_argument("--interactive", action="store_true", help="Run in interactive mode")
    parser.add_argument("--message", type=str, help="Single message to send to the assistant")
    args = parser.parse_args()
    
    if args.interactive:
        print("=== Italian Language Assistant CLI ===")
        print("Type 'exit' or 'quit' to end the conversation.")
        print("Type 'clear' to clear the chat history.")
        print("")
        
        chat_history = []
        while True:
            try:
                user_input = input("You: ")
                if user_input.lower() in ["exit", "quit"]:
                    print("\nArrivederci! 👋")
                    break
                elif user_input.lower() == "clear":
                    chat_history = []
                    print("Chat history cleared.")
                    continue
                elif not user_input.strip():
                    continue
                
                chat_history, response = chat_with_gpt(user_input, chat_history)
                print(f"\nAssistant: {response}\n")
            except KeyboardInterrupt:
                print("\nArrivederci! 👋")
                break
            except Exception as e:
                print(f"\nError: {str(e)}\n")
    
    elif args.message:
        _, response = chat_with_gpt(args.message)
        print(f"\nAssistant: {response}\n")
    
    else:
        print("Please specify either --interactive or --message.")
        print("Example usage:")
        print("  python chat.py --interactive")
        print("  python chat.py --message 'Translate hello to Italian'")