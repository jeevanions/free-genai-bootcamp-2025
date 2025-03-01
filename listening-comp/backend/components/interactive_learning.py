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
try:
    llm = AzureChatOpenAI(
        api_key=os.getenv("OPENAI_API_KEY"),
        api_version=os.getenv("OPENAI_API_VERSION"),
        azure_endpoint=os.getenv("OPENAI_API_BASE"),
        deployment_name=os.getenv("OPENAI_API_DEPLOYMENT_NAME"),
        temperature=0.7,
    )
    USE_MOCK = False
except Exception as e:
    print(f"Warning: Error initializing Azure Chat OpenAI in interactive_learning: {e}")
    # Create a mock LLM for testing
    USE_MOCK = True
    
    class MockLLM:
        def invoke(self, prompt):
            # Return a simple mock response based on the exercise type
            if "listening comprehension" in prompt.lower():
                return "Here's a mock listening comprehension exercise for testing."
            elif "vocabulary" in prompt.lower():
                return "Here's a mock vocabulary exercise for testing."
            elif "grammar" in prompt.lower():
                return "Here's a mock grammar exercise for testing."
            else:
                return "Here's a mock conversation exercise for testing."
    
    llm = MockLLM()

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
        if USE_MOCK:
            # If we're using the mock LLM, provide a formatted mock response
            if exercise_type == "listening_comprehension":
                response = """{
                    "text": "Ciao! Mi chiamo Marco e sono di Roma. Mi piace studiare le lingue straniere. Parlo italiano, inglese e un po' di francese. Nel tempo libero, mi piace leggere libri e guardare film.",
                    "questions": [
                        {
                            "question": "Come si chiama la persona che parla?",
                            "options": ["Paolo", "Marco", "Giovanni"],
                            "correct_index": 1
                        },
                        {
                            "question": "Di dove è?",
                            "options": ["Milano", "Napoli", "Roma"],
                            "correct_index": 2
                        },
                        {
                            "question": "Quali lingue parla?",
                            "options": ["Italiano e spagnolo", "Italiano, inglese e un po' di francese", "Italiano, francese e tedesco"],
                            "correct_index": 1
                        }
                    ]
                }"""
            elif exercise_type == "vocabulary":
                response = """{
                    "words": [
                        {
                            "italian": "casa",
                            "english": "house",
                            "sentence": "La mia casa è piccola ma accogliente.",
                            "pronunciation": "KAH-sah"
                        },
                        {
                            "italian": "libro",
                            "english": "book",
                            "sentence": "Ho letto un libro interessante.",
                            "pronunciation": "LEE-bro"
                        },
                        {
                            "italian": "amico",
                            "english": "friend",
                            "sentence": "Il mio amico si chiama Giovanni.",
                            "pronunciation": "ah-MEE-ko"
                        },
                        {
                            "italian": "mangiare",
                            "english": "to eat",
                            "sentence": "Mi piace mangiare la pizza.",
                            "pronunciation": "man-JAH-reh"
                        },
                        {
                            "italian": "bello",
                            "english": "beautiful",
                            "sentence": "Che bella giornata!",
                            "pronunciation": "BEL-lo"
                        }
                    ]
                }"""
            else:
                response = "Mock exercise content for testing without API credentials."
        else:
            # Use the real LLM
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