"""
Main Gradio application for the Language Listening Comprehension App.
"""
import gradio as gr
from pathlib import Path
import os

from frontend.ui.components.chat_ui import create_chat_ui
from frontend.ui.components.youtube_ui import create_youtube_ui
from frontend.ui.components.whisper_ui import create_whisper_ui
from frontend.ui.components.ocr_ui import create_ocr_ui
from frontend.ui.components.rag_ui import create_rag_ui

# Create output directories if they don't exist
for dir_path in ["output/transcript-yt", "output/transcript-vid", "output/transcript-ocr"]:
    Path(dir_path).mkdir(parents=True, exist_ok=True)

# Custom CSS for a more professional look
custom_css = """
.header-container {
    text-align: center;
    margin-bottom: 20px;
    padding: 20px;
    background: linear-gradient(135deg, #006847 0%, #009246 100%);
}

/* Radio menu styling */
.menu-radio {
    width: 100%;
}

.menu-radio .gr-form {
    border: none !important;
    background: transparent !important;
    box-shadow: none !important;
}

.menu-radio .gr-block.gr-box {
    border: none !important;
    background: transparent !important;
    box-shadow: none !important;
    width: 100% !important;
}

.menu-radio .gr-radio-row {
    display: flex !important;
    align-items: center !important;
    padding: 10px 15px !important;
    margin: 5px 0 !important;
    border-radius: 5px !important;
    cursor: pointer !important;
    background-color: #f0f0f0 !important;
    border-left: 3px solid #009246 !important;
    transition: all 0.2s ease !important;
}

.menu-radio .gr-radio-row:hover {
    background-color: #e0f0e0 !important;
    transform: translateX(2px) !important;
}

.menu-radio .gr-radio-row[data-selected="true"] {
    background-color: #d0e8d0 !important;
    font-weight: bold !important;
}

.menu-radio .gr-radio-row input {
    margin-right: 10px !important;
}

/* Make sure the label is visible but styled appropriately */
.menu-radio label {
    font-weight: bold !important;
    color: #006847 !important;
    margin-bottom: 10px !important;
    font-size: 1.1em !important;
    display: block !important;
}

/* Menu instruction */
.menu-instruction {
    color: #006847 !important;
    margin-bottom: 10px !important;
    font-weight: bold !important;
    border-bottom: 1px solid #ddd !important;
    padding-bottom: 5px !important;
}

/* Layout containers */
.menu-container {
    padding: 15px;
    background-color: #f5f5f5;
    border-radius: 8px;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
    margin-right: 15px;
}

.content-container {
    background-color: #fff;
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
}

/* Content sections */
.content-section {
    margin-bottom: 30px;
}

/* Content titles */
.content-title {
    font-size: 1.5rem;
    font-weight: bold;
    color: #006847;
    margin-bottom: 15px;
    padding-bottom: 10px;
    border-bottom: 2px solid #009246;
}

/* Content description */
.content-description {
    background-color: #f9f9f9;
    border-left: 4px solid #009246;
    padding: 10px 15px;
    margin-bottom: 20px;
    border-radius: 0 4px 4px 0;
}

.menu-title {
    font-size: 1.2rem;
    font-weight: bold;
    color: #006847;
    margin-top: 20px;
    margin-bottom: 10px;
}

/* Form elements */
.gradio-container input, .gradio-container textarea, .gradio-container select {
    border: 1px solid #ccc !important;
    border-radius: 4px !important;
    padding: 8px !important;
}

.gradio-container button {
    background-color: #006847 !important;
    color: white !important;
    border: none !important;
    border-radius: 4px !important;
    padding: 8px 16px !important;
    cursor: pointer !important;
    transition: background-color 0.3s !important;
}

.gradio-container button:hover {
    background-color: #009246 !important;
}

/* Ensure content doesn't stretch */
.gradio-container .prose {
    max-width: 100% !important;
}

/* Labels for form elements */
.gradio-container label {
    font-weight: bold !important;
    color: #333 !important;
    margin-bottom: 5px !important;
}

    color: white;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}
.header-title {
    font-size: 2.2em !important;
    font-weight: 700 !important;
    margin-bottom: 5px !important;
}
.header-subtitle {
    font-size: 1.2em !important;
    font-weight: 400 !important;
}
/* Sidebar styling */
.sidebar-column {
    padding-right: 15px;
}

.menu-container, .info-container, .examples-container {
    background-color: #f5f5f5;
    border-radius: 8px;
    padding: 15px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    margin-bottom: 15px;
}

/* Example buttons */
.example-buttons {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.example-button {
    text-align: left !important;
    background-color: #fff !important;
    border: 1px solid #ddd !important;
    border-left: 3px solid #009246 !important;
    border-radius: 4px !important;
    padding: 8px 12px !important;
    font-size: 0.9em !important;
    transition: all 0.3s ease !important;
}

.example-button:hover {
    background-color: #f0f9f0 !important;
    border-left-color: #006847 !important;
    transform: translateX(2px) !important;
}
.menu-title {
    font-weight: 600 !important;
    color: #333;
    margin-bottom: 15px !important;
    text-align: center;
}
.menu-button {
    margin-bottom: 8px !important;
    border-left: 3px solid #CD212A !important;
    transition: all 0.3s ease;
}
.menu-button:hover {
    background-color: #f0f0f0 !important;
    transform: translateX(3px);
}
.menu-button.selected {
    background-color: #f0f0f0 !important;
    border-left: 3px solid #009246 !important;
}
.content-container {
    background-color: white;
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}
.content-title {
    color: #333;
    border-bottom: 2px solid #009246;
    padding-bottom: 10px;
    margin-bottom: 20px !important;
}
.footer {
    text-align: center;
    margin-top: 20px;
    font-size: 0.9em;
    color: #666;
}
"""

def create_app():
    """Create and configure the Gradio application."""
    
    with gr.Blocks(title="Italian Language Learning Platform", 
                  theme=gr.themes.Soft(primary_hue="green"),
                  css=custom_css) as app:
        
        # Header with Italian flag colors and centered text
        with gr.Row(elem_classes="header-container"):
            with gr.Column():
                gr.Markdown(
                    """
                    <div class="header-title">🎧 Italian Language Learning Platform</div>
                    <div class="header-subtitle">Improve your Italian listening comprehension with AI-powered tools</div>
                    """
                )
        
        with gr.Row():
            # Left sidebar for navigation with menu items styled like the inspiration image
            with gr.Column(scale=1, elem_classes="sidebar-column"):
                with gr.Group(elem_classes="menu-container"):
                    gr.Markdown("## Navigation Menu", elem_classes="menu-title")
                    
                    # Store the currently selected button state
                    selected_btn = gr.State("chat")
                    
                    # Use Radio component for menu selection
                    menu_options = [
                        "Chat with GPT-4o", 
                        "Raw Transcript", 
                        "Structured Data", 
                        "OCR Extraction", 
                        "RAG Implementation", 
                        "Interactive Learning"
                    ]
                    
                    gr.Markdown("### Select an option below:", elem_classes="menu-instruction")
                    menu_radio = gr.Radio(
                        choices=menu_options,
                        value="Chat with GPT-4o",
                        label="Navigation Options",
                        elem_classes="menu-radio",
                        container=True,
                        interactive=True,
                        visible=True
                    )
                
                # Removed Current Focus and Examples sections
            
            # Main content area
            with gr.Column(scale=3, elem_classes="content-container"):
                # Use a state to track which content to show
                content_state = gr.State("chat")
                
                # Chat UI
                with gr.Group(visible=True, elem_classes="content-section") as chat_ui:
                    gr.Markdown("# Chat with GPT-4o", elem_classes="content-title")
                    with gr.Group(elem_classes="content-description"):
                        gr.Markdown("Start by exploring GPT-4o's base Italian language capabilities. Try asking questions about Italian grammar, vocabulary, or cultural aspects.")
                    create_chat_ui()
                
                # YouTube UI
                with gr.Group(visible=False, elem_classes="content-section") as youtube_ui:
                    gr.Markdown("# Raw Transcript", elem_classes="content-title")
                    with gr.Group(elem_classes="content-description"):
                        gr.Markdown("Extract raw transcripts from Italian YouTube videos to practice your listening skills.")
                    create_youtube_ui()
                
                # Whisper UI
                with gr.Group(visible=False, elem_classes="content-section") as whisper_ui:
                    gr.Markdown("# Structured Data", elem_classes="content-title")
                    with gr.Group(elem_classes="content-description"):
                        gr.Markdown("Generate structured transcripts from Italian audio using Whisper.")
                    create_whisper_ui()
                
                # OCR UI
                with gr.Group(visible=False, elem_classes="content-section") as ocr_ui:
                    gr.Markdown("# OCR Extraction", elem_classes="content-title")
                    with gr.Group(elem_classes="content-description"):
                        gr.Markdown("Extract Italian text from images using OCR technology.")
                    create_ocr_ui()
                
                # RAG UI
                with gr.Group(visible=False, elem_classes="content-section") as rag_ui:
                    gr.Markdown("# RAG Implementation", elem_classes="content-title")
                    with gr.Group(elem_classes="content-description"):
                        gr.Markdown("Enhance your learning with Retrieval-Augmented Generation for Italian language content.")
                    create_rag_ui()
                
                # Learning UI
                with gr.Group(visible=False, elem_classes="content-section") as learning_ui:
                    gr.Markdown("# Interactive Learning", elem_classes="content-title")
                    with gr.Group(elem_classes="content-description"):
                        gr.Markdown("Practice your Italian with interactive exercises and quizzes.")
                        gr.Markdown("This feature will be implemented in a future update.")
                
                # Map radio values to internal state values
                def get_state_from_radio(radio_value):
                    mapping = {
                        "Chat with GPT-4o": "chat",
                        "Raw Transcript": "youtube",
                        "Structured Data": "whisper",
                        "OCR Extraction": "ocr",
                        "RAG Implementation": "rag",
                        "Interactive Learning": "learning"
                    }
                    return mapping.get(radio_value, "chat")
                
                # Function to update UI visibility based on selection
                def update_ui_visibility(selected):
                    chat_visible = selected == "chat"
                    youtube_visible = selected == "youtube"
                    whisper_visible = selected == "whisper"
                    ocr_visible = selected == "ocr"
                    rag_visible = selected == "rag"
                    learning_visible = selected == "learning"
                    return chat_visible, youtube_visible, whisper_visible, ocr_visible, rag_visible, learning_visible
                
                # Set up button click events
                # Create visibility outputs for each UI component
                chat_visible = gr.Checkbox(value=True, visible=False)
                youtube_visible = gr.Checkbox(value=False, visible=False)
                whisper_visible = gr.Checkbox(value=False, visible=False)
                ocr_visible = gr.Checkbox(value=False, visible=False)
                rag_visible = gr.Checkbox(value=False, visible=False)
                learning_visible = gr.Checkbox(value=False, visible=False)
                
                # Connect visibility checkboxes to UI components
                chat_visible.change(lambda x: gr.update(visible=x), inputs=[chat_visible], outputs=[chat_ui])
                youtube_visible.change(lambda x: gr.update(visible=x), inputs=[youtube_visible], outputs=[youtube_ui])
                whisper_visible.change(lambda x: gr.update(visible=x), inputs=[whisper_visible], outputs=[whisper_ui])
                ocr_visible.change(lambda x: gr.update(visible=x), inputs=[ocr_visible], outputs=[ocr_ui])
                rag_visible.change(lambda x: gr.update(visible=x), inputs=[rag_visible], outputs=[rag_ui])
                learning_visible.change(lambda x: gr.update(visible=x), inputs=[learning_visible], outputs=[learning_ui])
                
                # Set up radio change event
                menu_radio.change(
                    # First set the internal state
                    get_state_from_radio,
                    inputs=[menu_radio],
                    outputs=[selected_btn]
                ).then(
                    # Then update UI visibility based on the state
                    update_ui_visibility,
                    inputs=[selected_btn],
                    outputs=[chat_visible, youtube_visible, whisper_visible, ocr_visible, rag_visible, learning_visible]
                )
        
        # Footer
        with gr.Row(elem_classes="footer"):
            gr.Markdown("© 2025 Italian Language Learning Platform | Powered by AI")
    
    return app

if __name__ == "__main__":
    app = create_app()
    app.launch(share=False)
