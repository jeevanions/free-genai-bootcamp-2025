"""
Chat UI component for interacting with GPT-4.
"""
import gradio as gr
from backend.components.chat import ChatComponent

def create_chat_ui():
    """Create the chat UI component."""
    chat_component = ChatComponent()
    
    with gr.Column():
        chatbot = gr.Chatbot(
            height=450, 
            type="messages",
            show_copy_button=True,
            avatar_images=(None, "https://upload.wikimedia.org/wikipedia/en/0/03/Flag_of_Italy.svg"),
            bubble_full_width=False
        )
        
        with gr.Row():
            with gr.Column(scale=10):
                msg = gr.Textbox(
                    placeholder="Ask about Italian language, grammar, or culture...", 
                    label="Your message",
                    show_label=False,
                    container=False
                )
            with gr.Column(scale=1):
                clear = gr.Button("Clear", variant="secondary")
        
        with gr.Accordion("Debug Information", open=False):
            debug_info = gr.JSON(value={})
        
        # Store conversation state
        state = gr.State([{"role": "system", "content": "You are a helpful language learning assistant specializing in Italian. Provide accurate, concise information about Italian language, grammar, vocabulary, and culture. When appropriate, include examples and practice exercises."}])
        
        def user_message(message, history, conversation):
            """Process user message and update conversation."""
            if not message.strip():
                return "", history, conversation
                
            # Add user message to conversation
            conversation.append({"role": "user", "content": message})
            
            # Get response from GPT-4
            response = chat_component.chat(conversation)
            
            # Add assistant response to conversation
            conversation.append({"role": "assistant", "content": response})
            
            # Update chat history using the messages format
            history.append({"role": "user", "content": message})
            history.append({"role": "assistant", "content": response})
            
            # Update debug info
            return "", history, conversation, {"message_count": len(conversation), "last_update": "now"}
        
        def clear_chat():
            """Clear the chat history."""
            return None, chat_component.clear_conversation(), {"status": "cleared"}
        
        # Connect UI components
        msg.submit(user_message, [msg, chatbot, state], [msg, chatbot, state, debug_info])
        clear.click(clear_chat, None, [chatbot, state, debug_info])
    
    return gr.Column()
