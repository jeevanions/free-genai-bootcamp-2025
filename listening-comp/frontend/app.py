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
from backend.components.structured_data import create_structured_data_interface

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

# Create the Gradio app with Italian-themed styling
with gr.Blocks(
    theme=gr.themes.Default(),
    css="""
    body {
        background-color: #111827;
        color: #fff;
    }
    
    /* Header styling */
    #app-header {
        background: linear-gradient(90deg, #008C45 33%, #F4F5F0 33%, #F4F5F0 66%, #CD212A 66%);
        padding: 1.5rem;
        margin-bottom: 0.5rem;
        border-radius: 10px;
        text-align: center;
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
    }
    
    #app-header h1 {
        color: #1A1A1A;
        font-weight: 700;
        margin: 0;
        font-size: 1.8rem;
    }
    
    #app-header p {
        color: #555;
        margin-top: 0.5rem;
        font-size: 1rem;
    }
    
    /* Tab buttons styling */
    .tab-buttons-container {
        display: flex;
        margin-bottom: 0.5rem;
        background-color: #1e293b;
        border-radius: 10px;
        overflow: hidden;
    }
    
    .tab-btn {
        flex: 1;
        background-color: #1e293b !important;
        color: white !important;
        border: none !important;
        padding: 0.7rem 0.5rem !important;
        font-size: 0.85rem !important;
        text-align: center !important;
        border-radius: 0 !important;
        margin: 0 !important;
        transition: background-color 0.3s;
    }
    
    .tab-btn:hover {
        background-color: #2d3748 !important;
    }
    
    .tab-btn.active-tab {
        background-color: #008C45 !important;
    }
    
    /* Main content area */
    .main-content-area {
        background-color: white;
        border-radius: 10px;
        overflow: hidden;
    }
    
    /* Hide default tab navigation */
    .tabs > .tab-nav {
        display: none !important;
    }
    
    /* Footer styling */
    .footer {
        background-color: #1e293b;
        color: white;
        text-align: center;
        padding: 1rem;
        border-radius: 10px;
        margin-top: 1rem;
    }
    
    /* Chat interface styling */
    .chatbot {
        height: 400px;
        background-color: #1e293b !important;
        border-radius: 0 !important;
    }
    
    .chatbot .message.user {
        background-color: #2d3748 !important;
    }
    
    .chatbot .message.bot {
        background-color: #374151 !important;
    }
    
    .chat-input-container {
        background-color: white !important;
        padding: 0.5rem !important;
        border-top: 1px solid #e5e7eb;
    }
    
    .chat-input {
        border: 1px solid #d1d5db !important;
        border-radius: 0.375rem !important;
    }
    
    .send-btn {
        background-color: #008C45 !important;
        color: white !important;
    }
    
    .clear-btn {
        background-color: #f3f4f6 !important;
        color: #374151 !important;
        border: 1px solid #d1d5db !important;
    }
    
    /* Form elements styling */
    input, select, textarea {
        background-color: white !important;
        border: 1px solid #d1d5db !important;
        border-radius: 0.375rem !important;
    }
    
    button.primary {
        background-color: #008C45 !important;
        color: white !important;
    }
    
    /* Tab content styling */
    .tab-content {
        padding: 1.5rem;
    }
    
    /* Gradio component overrides */
    .gradio-container {
        max-width: 100% !important;
    }
    
    .dark .gr-button-primary {
        background-color: #008C45 !important;
    }
    
    .dark .gr-button-secondary {
        background-color: #1e293b !important;
    }
    
    .dark .gr-input, .dark .gr-textarea, .dark .gr-dropdown {
        background-color: white !important;
        color: #111827 !important;
    }
    
    .dark .gr-panel {
        background-color: white !important;
        color: #111827 !important;
    }
    
    .dark .gr-box {
        background-color: white !important;
    }
    
    .dark .gr-form {
        background-color: white !important;
    }
    
    .dark .gr-padded {
        background-color: white !important;
    }
    
    .dark .gr-compact {
        background-color: white !important;
    }
    
    /* Custom styling for chat interface */
    .chat-with-assistant {
        background-color: #1e293b;
        color: white;
        padding: 1rem;
        border-radius: 0.5rem 0.5rem 0 0;
    }
    
    .chat-description {
        color: #d1d5db;
        margin-bottom: 1rem;
    }
    """
) as app:
    # Header
    with gr.Row(elem_id="app-header"):
        gr.HTML("""
        <h1>🇮🇹 Italian Language Learning Platform 🇮🇹</h1>
        <p>Enhance your Italian language skills with AI-powered tools</p>
        """)
    
    # Tab buttons
    with gr.Row(elem_classes="tab-buttons-container"):
        chat_btn = gr.Button("💬 Chat with\nGPT", elem_classes="tab-btn active-tab")
        youtube_btn = gr.Button("📺 YouTube\nTranscript", elem_classes="tab-btn")
        whisper_btn = gr.Button("🎤 Whisper\nTranscript", elem_classes="tab-btn")
        ocr_btn = gr.Button("👁️ OCR\nExtraction", elem_classes="tab-btn")
        rag_btn = gr.Button("🔍 RAG\nImplementation", elem_classes="tab-btn")
        interactive_btn = gr.Button("🎮 Interactive\nLearning", elem_classes="tab-btn")
        structured_btn = gr.Button("📊 Structured\nData", elem_classes="tab-btn")
    
    # Main content area with tabs
    with gr.Row(elem_classes="main-content-area"):
        with gr.Tabs() as tabs:
            with gr.TabItem("Chat with GPT", id=0, elem_classes="tab-content"):
                chat_interface = create_chat_interface(gr.Column())
            
            with gr.TabItem("YouTube Transcript", id=1, elem_classes="tab-content"):
                youtube_interface = create_youtube_transcript_interface(gr.Column())
            
            with gr.TabItem("Whisper Transcript", id=2, elem_classes="tab-content"):
                whisper_interface = create_whisper_transcript_interface(gr.Column())
            
            with gr.TabItem("OCR Extraction", id=3, elem_classes="tab-content"):
                ocr_interface = create_ocr_interface(gr.Column())
            
            with gr.TabItem("RAG Implementation", id=4, elem_classes="tab-content"):
                rag_interface = create_rag_interface(gr.Column())
            
            with gr.TabItem("Interactive Learning", id=5, elem_classes="tab-content"):
                interactive_interface = create_interactive_learning_interface(gr.Column())
            
            with gr.TabItem("Structured Data", id=6, elem_classes="tab-content"):
                structured_interface = create_structured_data_interface(gr.Column())
    
    # Footer
    with gr.Row(elem_classes="footer"):
        gr.HTML("""
        <div>
            <p>© 2025 Italian Language Learning Platform. All rights reserved.</p>
            <p style="margin-top: 0.5rem;">Powered by Gradio, OpenAI, and Qdrant</p>
        </div>
        """)
    
    # Set up event handlers for navigation buttons
    tab_buttons = [chat_btn, youtube_btn, whisper_btn, ocr_btn, rag_btn, interactive_btn, structured_btn]
    
    # Function to update button styles based on active tab
    def update_button_styles(tab_index):
        return [
            gr.update(elem_classes="tab-btn active-tab" if i == tab_index else "tab-btn")
            for i in range(len(tab_buttons))
        ]
    
    # Function to change tabs and update button styles
    def change_tab(tab_index):
        # First update the tab selection
        tab_update = gr.update(selected=tab_index)
        # Then update all button styles
        button_updates = update_button_styles(tab_index)
        # Return all updates
        return [tab_update] + button_updates
    
    # Connect each button to its corresponding tab and update button styles
    chat_btn.click(lambda: change_tab(0), None, [tabs] + tab_buttons)
    youtube_btn.click(lambda: change_tab(1), None, [tabs] + tab_buttons)
    whisper_btn.click(lambda: change_tab(2), None, [tabs] + tab_buttons)
    ocr_btn.click(lambda: change_tab(3), None, [tabs] + tab_buttons)
    rag_btn.click(lambda: change_tab(4), None, [tabs] + tab_buttons)
    interactive_btn.click(lambda: change_tab(5), None, [tabs] + tab_buttons)
    structured_btn.click(lambda: change_tab(6), None, [tabs] + tab_buttons)

# Launch the app
if __name__ == "__main__":
    app.launch(server_name="0.0.0.0", server_port=7861)