import gradio as gr
import os
from gtts import gTTS
import tempfile
import random
from langchain_openai import AzureChatOpenAI
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Initialize LLM
llm = AzureChatOpenAI(
    openai_api_key=os.getenv("OPENAI_API_KEY"),
    openai_api_version=os.getenv("OPENAI_API_VERSION"),
    openai_api_base=os.getenv("OPENAI_API_BASE"),
    openai_api_type=os.getenv("OPENAI_API_TYPE"),
    deployment_name=os.getenv("OPENAI_API_DEPLOYMENT_NAME"),
    temperature=0.7,
)

def text_to_speech(text, lang="it"):
    """
    Convert text to speech using gTTS
    """
    try:
        # Create a temporary file
        with tempfile.NamedTemporaryFile(delete=False, suffix=".mp3") as temp_file:
            temp_path = temp_file.name
        
        # Generate speech
        tts = gTTS(text=text, lang=lang, slow=False)
        tts.save(temp_path)
        
        return {"success": True, "path": temp_path}
    except Exception as e:
        return {"success": False, "message": f"Error generating speech: {str(e)}"}

def generate_exercise(difficulty, exercise_type):
    """
    Generate a language exercise using GPT-4
    """
    try:
        # Create prompt based on exercise type and difficulty
        if exercise_type == "listening_comprehension":
            prompt = f"""
            Create an Italian listening comprehension exercise at {difficulty} level.
            Include:
            1. A short paragraph in Italian (5-8 sentences)
            2. 3 multiple-choice questions about the content
            3. The correct answers
            
            Format your response as a JSON object with these keys:
            - text: the Italian paragraph
            - questions: array of question objects with 'question', 'options' (array), and 'correct_index'
            """
        elif exercise_type == "vocabulary":
            prompt = f"""
            Create an Italian vocabulary exercise at {difficulty} level.
            Include:
            1. 5 Italian words with their English translations
            2. A sample sentence for each word
            3. Audio pronunciation hints (describe how to pronounce)
            
            Format your response as a JSON object with these keys:
            - words: array of word objects with 'italian', 'english', 'sentence', and 'pronunciation'
            """
        elif exercise_type == "grammar":
            prompt = f"""
            Create an Italian grammar exercise at {difficulty} level.
            Include:
            1. A brief explanation of a grammar rule
            2. 5 fill-in-the-blank sentences practicing this rule
            3. The correct answers
            
            Format your response as a JSON object with these keys:
            - rule: explanation of the grammar rule
            - sentences: array of sentence objects with 'incomplete', 'complete', and 'explanation'
            """
        else:
            prompt = f"""
            Create an Italian conversation practice exercise at {difficulty} level.
            Include:
            1. A dialogue between two people (5-8 exchanges)
            2. English translations for each line
            3. 3 questions about the dialogue
            
            Format your response as a JSON object with these keys:
            - dialogue: array of exchange objects with 'speaker', 'italian', and 'english'
            - questions: array of question objects with 'question', 'answer'
            """
        
        # Call the LLM
        response = llm.invoke(prompt)
        
        return {"success": True, "exercise": response}
    except Exception as e:
        return {"success": False, "message": f"Error generating exercise: {str(e)}"}

def create_interactive_learning_interface(parent):
    """
    Creates the interactive learning interface
    """
    with parent:
        gr.Markdown("### Interactive Italian Learning")
        gr.Markdown("Generate and practice with interactive Italian language exercises.")
        
        with gr.Tabs():
            with gr.Tab("Generate Exercise"):
                with gr.Row():
                    difficulty_dropdown = gr.Dropdown(
                        label="Difficulty Level",
                        choices=["beginner", "intermediate", "advanced"],
                        value="beginner",
                        info="Select the difficulty level of the exercise"
                    )
                    
                    exercise_type_dropdown = gr.Dropdown(
                        label="Exercise Type",
                        choices=[
                            "listening_comprehension",
                            "vocabulary",
                            "grammar",
                            "conversation"
                        ],
                        value="listening_comprehension",
                        info="Select the type of exercise"
                    )
                
                with gr.Row():
                    generate_btn = gr.Button("Generate Exercise")
                
                with gr.Row():
                    exercise_output = gr.Textbox(
                        label="Exercise",
                        lines=15,
                    )
            
            with gr.Tab("Text to Speech"):
                with gr.Row():
                    tts_input = gr.Textbox(
                        label="Italian Text",
                        placeholder="Enter Italian text to convert to speech...",
                        lines=5,
                    )
                
                with gr.Row():
                    tts_btn = gr.Button("Generate Speech")
                
                with gr.Row():
                    audio_output = gr.Audio(label="Audio Output")
        
        # Set up event handlers
        generate_btn.click(
            lambda difficulty, exercise_type: generate_exercise(difficulty, exercise_type)["exercise"],
            inputs=[difficulty_dropdown, exercise_type_dropdown],
            outputs=[exercise_output],
        )
        
        tts_btn.click(
            lambda text: text_to_speech(text)["path"] if text_to_speech(text)["success"] else None,
            inputs=[tts_input],
            outputs=[audio_output],
        )
        
        return None