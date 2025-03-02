import gradio as gr
import os
import json
import time
import glob
from openai import AzureOpenAI
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Initialize Azure OpenAI client
azure_endpoint = os.getenv("OPENAI_API_BASE")
if azure_endpoint and not azure_endpoint.startswith("https://"):
    azure_endpoint = f"https://{azure_endpoint}"

client = AzureOpenAI(
    api_key=os.getenv("OPENAI_API_KEY"),
    api_version=os.getenv("OPENAI_API_VERSION"),
    azure_endpoint=azure_endpoint,
)

# Schema for structured data
JSON_SCHEMA = """{
  "metadata": {
    "video_title": "Example: 'Italian for Beginners - At the Café'",
    "language_level": "A1",
    "topic": "Example: 'Ordering Food'",
    "video_timestamp": "00:00-05:30",
    "context": "Brief description of the scenario (e.g., 'A conversation between a customer and barista')"
  },
  "dialogue": [
    {
      "speaker": "Barista",
      "italian_text": "Buongiorno! Cosa desidera?",
      "english_translation": "Good morning! What would you like?",
      "audio_timestamp": "00:15",
      "key_phrases": ["Buongiorno", "Cosa desidera?"]
    },
    {
      "speaker": "Customer",
      "italian_text": "Un caffè, per favore.",
      "english_translation": "A coffee, please.",
      "audio_timestamp": "00:20",
      "key_phrases": ["Un caffè", "per favore"]
    }
  ],
  "vocabulary": [
    {
      "italian_term": "caffè",
      "part_of_speech": "noun",
      "english_translation": "coffee",
      "example_sentence": "Vorrei un caffè.",
      "ocr_reference": "Text from video visuals (if applicable)"
    }
  ],
  "grammar_concepts": [
    {
      "concept": "Polite requests with 'per favore'",
      "explanation": "Used to make polite requests in Italian",
      "examples": ["Un cappuccino, per favore."]
    }
  ],
  "exercises": [
    {
      "type": "multiple_choice",
      "question": "Cosa dice il barista?",
      "options": ["Buongiorno!", "Grazie!", "Arrivederci!"],
      "correct_answer": "Buongiorno!",
      "audio_timestamp": "00:15"
    }
  ],
  "cultural_notes": [
    {
      "topic": "Italian café culture",
      "description": "Italians often drink espresso standing at the bar",
      "connection_to_dialogue": "Reference to 'caffè' ordering"
    }
  ],
  "ocr_extracted_content": [
    {
      "text": "Menu: Caffè €1, Cornetto €2",
      "type": "menu",
      "timestamp": "00:30"
    }
  ]
}"""

def get_available_transcripts():
    """Get all available transcripts from output folders"""
    transcripts = []
    
    # YouTube transcripts
    yt_transcripts = glob.glob("output/transcript-yt/*.txt")
    for transcript in yt_transcripts:
        filename = os.path.basename(transcript)
        transcripts.append({"source": "YouTube", "path": transcript, "name": filename})
    
    # Whisper transcripts
    whisper_transcripts = glob.glob("output/transcript-vid/*.txt")
    for transcript in whisper_transcripts:
        filename = os.path.basename(transcript)
        transcripts.append({"source": "Whisper", "path": transcript, "name": filename})
    
    # OCR transcripts
    ocr_transcripts = glob.glob("output/transcript-ocr/*.txt")
    for transcript in ocr_transcripts:
        filename = os.path.basename(transcript)
        transcripts.append({"source": "OCR", "path": transcript, "name": filename})
    
    return transcripts

def format_transcript_options(transcripts):
    """Format transcript options for dropdown"""
    options = []
    for transcript in transcripts:
        options.append(f"{transcript['source']}: {transcript['name']}")
    return options

def get_transcript_content(selected_option, transcripts):
    """Get the content of the selected transcript"""
    if not selected_option:
        return ""
    
    # Parse the selected option to get source and name
    parts = selected_option.split(": ", 1)
    if len(parts) != 2:
        return ""
    
    source, name = parts
    
    # Find the matching transcript
    for transcript in transcripts:
        if transcript['source'] == source and transcript['name'] == name:
            try:
                with open(transcript['path'], 'r', encoding='utf-8') as f:
                    return f.read()
            except Exception as e:
                return f"Error reading transcript: {str(e)}"
    
    return "Transcript not found"

def generate_structured_data(transcript_text, status_box):
    """Generate structured data from transcript using Azure OpenAI"""
    if not transcript_text:
        return "Please select a transcript first", None
    
    try:
        # Update status - using value property instead of update method
        status_message = "Generating structured data... This may take a minute."
        
        # Prepare the prompt
        prompt = f"""Convert the following Italian A1 listening comprehension video transcript and OCR text into structured JSON data. Follow these rules:

Identify speakers and separate dialogue lines
Extract vocabulary with translations and examples
Highlight grammar concepts from context
Create simple multiple-choice questions
Note cultural references
Include OCR text with timestamps

Use this JSON schema:
{JSON_SCHEMA}

Transcript:
{transcript_text}"""

        # Call Azure OpenAI
        deployment_name = os.getenv("OPENAI_API_DEPLOYMENT_NAME")
        response = client.chat.completions.create(
            model=deployment_name,
            messages=[{"role": "system", "content": "You are a helpful assistant that converts Italian language learning content into structured data."},
                     {"role": "user", "content": prompt}],
            temperature=0.7,
            max_tokens=4000
        )
        
        # Extract the JSON response
        json_response = response.choices[0].message.content
        
        # Try to parse the JSON to ensure it's valid
        try:
            # Extract just the JSON part if there's any markdown formatting
            if "```json" in json_response:
                json_response = json_response.split("```json")[1].split("```")[0].strip()
            elif "```" in json_response:
                json_response = json_response.split("```")[1].split("```")[0].strip()
            
            parsed_json = json.loads(json_response)
            formatted_json = json.dumps(parsed_json, indent=2)
            
            # Save the structured data
            timestamp = int(time.time())
            output_path = f"output/structured-data/structured_{timestamp}.json"
            os.makedirs("output/structured-data", exist_ok=True)
            
            with open(output_path, 'w', encoding='utf-8') as f:
                f.write(formatted_json)
            
            success_message = f"✅ Structured data generated successfully and saved to {output_path}"
            return formatted_json, parsed_json
            
        except json.JSONDecodeError as e:
            error_message = f"❌ Error parsing JSON response: {str(e)}"
            return json_response, None
            
    except Exception as e:
        error_message = f"❌ Error generating structured data: {str(e)}"
        return error_message, None

def create_structured_data_interface(parent):
    """Create the structured data interface"""
    with parent:
        gr.Markdown("## Generate Structured Data from Transcripts", elem_classes="content-section")
        gr.Markdown("""
        Convert transcripts from YouTube, Whisper, or OCR into structured JSON data for language learning.
        This tool uses Azure OpenAI to analyze the content and organize it into a structured format.
        """, elem_classes="content-section")
        
        # Get available transcripts
        transcripts = get_available_transcripts()
        transcript_options = format_transcript_options(transcripts)
        
        with gr.Row():
            with gr.Column(scale=1):
                gr.Markdown("### Select Transcript", elem_classes="nav-header")
                transcript_dropdown = gr.Dropdown(
                    choices=transcript_options,
                    label="Available Transcripts",
                    info="Select a transcript to convert to structured data"
                )
                
                preview_btn = gr.Button("Preview Transcript", variant="secondary")
                generate_btn = gr.Button("Generate Structured Data", variant="primary")
                status_box = gr.Textbox(label="Status", interactive=False)
                
            with gr.Column(scale=2):
                with gr.Tabs():
                    with gr.TabItem("Transcript Preview"):
                        transcript_preview = gr.TextArea(
                            label="Transcript Content",
                            interactive=False,
                            lines=15
                        )
                    
                    with gr.TabItem("Structured Data (JSON)"):
                        json_output = gr.JSON(label="Structured Data")
                    
                    with gr.TabItem("Dialogue"):
                        dialogue_output = gr.Dataframe(
                            headers=["Speaker", "Italian Text", "English Translation", "Timestamp"],
                            label="Dialogue"
                        )
                    
                    with gr.TabItem("Vocabulary"):
                        vocabulary_output = gr.Dataframe(
                            headers=["Term", "Part of Speech", "Translation", "Example"],
                            label="Vocabulary"
                        )
                    
                    with gr.TabItem("Grammar & Exercises"):
                        grammar_output = gr.Dataframe(
                            headers=["Concept", "Explanation", "Examples"],
                            label="Grammar Concepts"
                        )
                        
                        gr.Markdown("### Practice Exercises", elem_classes="nav-header")
                        exercises_output = gr.HTML(label="Exercises")
        
        # Set up event handlers
        preview_btn.click(
            lambda x: get_transcript_content(x, transcripts),
            inputs=[transcript_dropdown],
            outputs=[transcript_preview]
        )
        
        def process_structured_data(transcript_option):
            # Get transcript content
            transcript_content = get_transcript_content(transcript_option, transcripts)
            if not transcript_content:
                return "Please select a valid transcript", None, [], [], [], ""
            
            # Update status message first
            status_box.value = "Generating structured data... This may take a minute."
            
            # Generate structured data
            json_text, json_data = generate_structured_data(transcript_content, status_box)
            
            # Update status message with result
            if json_data:
                status_box.value = f"✅ Structured data generated successfully"
            else:
                status_box.value = f"❌ Error generating structured data"
            
            if json_data:
                # Process dialogue
                dialogue_data = []
                if "dialogue" in json_data:
                    for entry in json_data["dialogue"]:
                        dialogue_data.append([
                            entry.get("speaker", ""),
                            entry.get("italian_text", ""),
                            entry.get("english_translation", ""),
                            entry.get("audio_timestamp", "")
                        ])
                
                # Process vocabulary
                vocabulary_data = []
                if "vocabulary" in json_data:
                    for entry in json_data["vocabulary"]:
                        vocabulary_data.append([
                            entry.get("italian_term", ""),
                            entry.get("part_of_speech", ""),
                            entry.get("english_translation", ""),
                            entry.get("example_sentence", "")
                        ])
                
                # Process grammar
                grammar_data = []
                if "grammar_concepts" in json_data:
                    for entry in json_data["grammar_concepts"]:
                        grammar_data.append([
                            entry.get("concept", ""),
                            entry.get("explanation", ""),
                            ", ".join(entry.get("examples", []))
                        ])
                
                # Process exercises
                exercises_html = "<div class='content-section'>"
                if "exercises" in json_data:
                    for i, exercise in enumerate(json_data["exercises"]):
                        exercises_html += f"<div style='margin-bottom: 1.5rem;'>"
                        exercises_html += f"<p><strong>Question {i+1}:</strong> {exercise.get('question', '')}</p>"
                        exercises_html += "<ul>"
                        for option in exercise.get("options", []):
                            correct = "✓ " if option == exercise.get("correct_answer", "") else ""
                            exercises_html += f"<li>{correct}{option}</li>"
                        exercises_html += "</ul>"
                        exercises_html += f"<p><small>Timestamp: {exercise.get('audio_timestamp', '')}</small></p>"
                        exercises_html += "</div>"
                exercises_html += "</div>"
                
                return json_text, json_data, dialogue_data, vocabulary_data, grammar_data, exercises_html
            else:
                return json_text, None, [], [], [], ""
        
        generate_btn.click(
            process_structured_data,
            inputs=[transcript_dropdown],
            outputs=[json_output, json_output, dialogue_output, vocabulary_output, grammar_output, exercises_output]
        )

if __name__ == "__main__":
    with gr.Blocks() as demo:
        create_structured_data_interface(gr.Group())
    demo.launch()
