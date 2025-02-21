import json
import random
from datetime import datetime, timedelta

# A1 Level basic words (sample - you'll need to expand this)
a1_words = {
    "Greetings and Basic Expressions": [
        {"italian": "ciao", "english": "hello/bye", "parts": {"type": "interjection", "usage": ["greeting", "farewell"]}},
        {"italian": "grazie", "english": "thank you", "parts": {"type": "interjection", "usage": ["gratitude"]}},
        {"italian": "prego", "english": "you're welcome", "parts": {"type": "interjection", "usage": ["response"]}},
        # Add more words here
    ],
    "Numbers and Time": [
        {"italian": "uno", "english": "one", "parts": {"type": "number"}},
        {"italian": "due", "english": "two", "parts": {"type": "number"}},
        {"italian": "ora", "english": "hour", "parts": {"type": "noun", "gender": "feminine"}},
        # Add more words here
    ],
    # Add more categories
}

# A2 Level words (sample - you'll need to expand this)
a2_words = {
    "House and Furniture": [
        {"italian": "tavolo", "english": "table", "parts": {"type": "noun", "gender": "masculine"}},
        {"italian": "sedia", "english": "chair", "parts": {"type": "noun", "gender": "feminine"}},
        # Add more words here
    ],
    "Travel and Transportation": [
        {"italian": "biglietto", "english": "ticket", "parts": {"type": "noun", "gender": "masculine"}},
        {"italian": "treno", "english": "train", "parts": {"type": "noun", "gender": "masculine"}},
        # Add more words here
    ],
    # Add more categories
}

def generate_word_data():
    words = []
    word_id = 1
    
    # Process A1 words
    for category, word_list in a1_words.items():
        for word in word_list:
            word_data = {
                "id": word_id,
                "italian": word["italian"],
                "english": word["english"],
                "parts": word["parts"],
                "created_at": (datetime.now() - timedelta(days=random.randint(0, 30))).isoformat()
            }
            words.append(word_data)
            word_id += 1
    
    # Process A2 words
    for category, word_list in a2_words.items():
        for word in word_list:
            word_data = {
                "id": word_id,
                "italian": word["italian"],
                "english": word["english"],
                "parts": word["parts"],
                "created_at": (datetime.now() - timedelta(days=random.randint(0, 30))).isoformat()
            }
            words.append(word_data)
            word_id += 1
    
    return {"words": words}

# Generate and save word data
word_data = generate_word_data()
with open('words.json', 'w', encoding='utf-8') as f:
    json.dump(word_data, f, ensure_ascii=False, indent=2)

print(f"Generated {len(word_data['words'])} words")
