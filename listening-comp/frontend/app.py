import gradio as gr
import os
import sys
from dotenv import load_dotenv

# Add the backend directory to the path so we can import from it
sys.path.append(os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'backend'))

from backend.components.chat import create_chat_interface
from backend.components.youtube_transcript import create_youtube_transcript_interface
from backend.components.whisper_transcript import create_whisper_transcript_interface
from backend.components.ocr_extraction import create_ocr_interface
from backend.components.rag_implementation import create_rag_interface
from backend.components.interactive_learning import create_interactive_learning_interface

# Load environment variables
load_dotenv()

# Create output directories if they don't exist
os.makedirs("output/transcript-yt", exist_ok=True)
os.makedirs("output/transcript-vid", exist_ok=True)
os.makedirs("output/transcript-ocr", exist_ok=True)

# Define Italian flag colors
GREEN = "#008C45"
WHITE = "#F4F5F0"
RED = "#CD212A"

# Function to change tabs
def change_tab(n):
    return gr.Tabs(selected=n)

# Create the Gradio app with Italian-themed styling
with gr.Blocks(
    css="""
    #app-header {
        background: linear-gradient(90deg, #008C45 33%, #F4F5F0 33%, #F4F5F0 66%, #CD212A 66%);
        padding: 1rem;
        margin-bottom: 1rem;
        border-radius: 8px;
        text-align: center;
    }
    #app-header h1 {
        color: #1A1A1A;
        text-shadow: 1px 1px 2px rgba(0,0,0,0.1);
    }
    .footer {
        text-align: center;
        margin-top: 2rem;
        padding: 1rem;
        border-top: 1px solid #ddd;
        font-size: 0.8rem;
    }
    .sidebar-btn {
        background-color: #F4F5F0;
        border: 1px solid #ddd;
        margin: 0.2rem 0;
        transition: all 0.3s;
    }
    .sidebar-btn:hover {
        background-color: #008C45;
        color: white;
    }
    """
) as app:
    # Header
    with gr.Row(elem_id="app-header"):
        gr.HTML("<h1>🇮🇹 Italian Language Learning Platform 🇮🇹</h1>")
    
    # Add a hidden state variable to track active tab
    active_tab = gr.State(value=0)
    
    # Main layout with sidebar and content area
    with gr.Row():
        # Sidebar for navigation
        with gr.Column(scale=1):
            gr.Markdown("### Navigation")
            
            with gr.Group():
                chat_btn = gr.Button("Chat with GPT-4", elem_classes="sidebar-btn")
                yt_transcript_btn = gr.Button("YouTube Transcript", elem_classes="sidebar-btn")
                whisper_btn = gr.Button("Whisper Transcript", elem_classes="sidebar-btn")
                ocr_btn = gr.Button("OCR Extraction", elem_classes="sidebar-btn")
                rag_btn = gr.Button("RAG Implementation", elem_classes="sidebar-btn")
                learning_btn = gr.Button("Interactive Learning", elem_classes="sidebar-btn")
        
        # Main content area with tabs
        with gr.Column(scale=4):
            tabs = gr.Tabs() 
            with tabs:
                with gr.TabItem("Chat with GPT-4", id=0):
                    create_chat_interface(gr.Group())
                    
                with gr.TabItem("YouTube Transcript", id=1):
                    create_youtube_transcript_interface(gr.Group())
                    
                with gr.TabItem("Whisper Transcript", id=2):
                    create_whisper_transcript_interface(gr.Group())
                    
                with gr.TabItem("OCR Extraction", id=3):
                    create_ocr_interface(gr.Group())
                    
                with gr.TabItem("RAG Implementation", id=4):
                    create_rag_interface(gr.Group())
                    
                with gr.TabItem("Interactive Learning", id=5):
                    create_interactive_learning_interface(gr.Group())
    
    # Footer
    with gr.Row(elem_classes="footer"):
        gr.Markdown("© 2025 Italian Language Learning Platform. All rights reserved.")
    
    # Connect buttons to tab selection function using the example approach
    chat_btn.click(lambda: 0, outputs=active_tab).then(
        change_tab, inputs=[active_tab], outputs=tabs
    )
    
    yt_transcript_btn.click(lambda: 1, outputs=active_tab).then(
        change_tab, inputs=[active_tab], outputs=tabs
    )
    
    whisper_btn.click(lambda: 2, outputs=active_tab).then(
        change_tab, inputs=[active_tab], outputs=tabs
    )
    
    ocr_btn.click(lambda: 3, outputs=active_tab).then(
        change_tab, inputs=[active_tab], outputs=tabs
    )
    
    rag_btn.click(lambda: 4, outputs=active_tab).then(
        change_tab, inputs=[active_tab], outputs=tabs
    )
    
    learning_btn.click(lambda: 5, outputs=active_tab).then(
        change_tab, inputs=[active_tab], outputs=tabs
    )

# Add API endpoint for the frontend to start the backend
def start_backend():
    return {"url": "http://localhost:7861", "message": "Backend started successfully"}

app.queue()

if __name__ == "__main__":
    app.launch(share=False, server_name="0.0.0.0", server_port=7861)