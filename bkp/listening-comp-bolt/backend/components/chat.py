import gradio as gr
import os
from openai import AzureOpenAI
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Initialize Azure OpenAI client
client = AzureOpenAI(
    api_key=os.getenv("OPENAI_API_KEY"),
    api_version=os.getenv("OPENAI_API_VERSION"),
    azure_endpoint=os.getenv("OPENAI_API_BASE"),
)

def chat_with_gpt(message, history):
    """
    Function to interact with GPT-4 model via Azure OpenAI
    """
    try:
        # Prepare conversation history for the API
        messages = []
        for human, assistant in history:
            messages.append({"role": "user", "content": human})
            messages.append({"role": "assistant", "content": assistant})
        
        # Add the new message
        messages.append({"role": "user", "content": message})
        
        # Call the OpenAI API
        response = client.chat.completions.create(
            model=os.getenv("OPENAI_API_DEPLOYMENT_NAME"),
            messages=messages,
            temperature=0.7,
            max_tokens=800,
        )
        
        # Extract and return the response
        return response.choices[0].message.content
    except Exception as e:
        return f"Error: {str(e)}"

def create_chat_interface(parent):
    """
    Creates the chat interface with GPT-4 model
    """
    with parent:
        gr.Markdown("### Chat with our Italian Language Assistant")
        gr.Markdown("Ask questions about Italian language, request translations, or get help with grammar.")
        
        chatbot = gr.Chatbot(
            avatar_images=["https://images.unsplash.com/photo-1560785496-3c9d27877182?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=100&q=80", None],
            height=400,
            bubble_full_width=False,
        )
        
        with gr.Row():
            msg = gr.Textbox(
                placeholder="Type your Italian language question here...",
                scale=9,
                container=False,
            )
            submit_btn = gr.Button("Send", scale=1)
        
        clear_btn = gr.Button("Clear Chat")
        
        with gr.Accordion("Debug Information", open=False):
            debug_info = gr.Textbox(label="API Response Details", interactive=False)
        
        # Set up event handlers
        submit_btn.click(chat_with_gpt, [msg, chatbot], [chatbot])
        msg.submit(chat_with_gpt, [msg, chatbot], [chatbot])
        clear_btn.click(lambda: None, None, chatbot, queue=False)
        
        return chatbot