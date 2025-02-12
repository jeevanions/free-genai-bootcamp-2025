
# References

https://docs.anthropic.com/en/docs/build-with-claude/prompt-engineering/overview


Writing effective prompts for Claude 3.5 Sonnet, or any large language model, requires clarity, specificity, and structure to ensure the model understands your intent and delivers the desired output. 

Here are some techniques to follow:

### 1. **Be Clear and Specific**
   - Clearly state what you want the model to do. Avoid vague or ambiguous language.
   - Example:
     - **Vague:** "Tell me about history."
     - **Specific:** "Provide a brief overview of the causes and key events of the American Civil War."

---

### 2. **Provide Context**
   - Give the model enough background information to generate a relevant response.
   - Example:
     - **Without context:** "Explain photosynthesis."
     - **With context:** "Explain the process of photosynthesis to a high school student, including the role of chlorophyll and sunlight."

---

### 3. **Use Direct Instructions**
   - Use action-oriented language to guide the model.
   - Example:
     - **Weak:** "Maybe you can write a poem."
     - **Strong:** "Write a 10-line poem about the ocean, using vivid imagery and a melancholic tone."

---

### 4. **Break Down Complex Tasks**
   - For multi-part tasks, break them into smaller, sequential steps.
   - Example:
     - **Complex:** "Analyze this data and write a report."
     - **Step-by-step:** "1. Summarize the key trends in this dataset. 2. Identify three insights. 3. Write a 200-word report explaining the insights."

---

### 5. **Specify Format and Tone**
   - Clearly define the desired format, style, or tone.
   - Example:
     - **Without format:** "Write a story."
     - **With format:** "Write a 300-word short story in the style of a mystery thriller, with a surprising twist at the end."

---

### 6. **Use Examples**
   - Provide examples to illustrate the type of output you want.
   - Example:
     - **Without example:** "Write a product description."
     - **With example:** "Write a product description like this: 'This sleek, ergonomic chair is designed for comfort and style, perfect for modern workspaces.' Now describe a smartwatch."

---

### 7. **Ask for Iterations or Refinements**
   - If the output isn’t perfect, ask the model to refine or adjust its response.
   - Example:
     - "That’s a good start, but can you make the tone more formal and add more details about the benefits?"

---

### 8. **Set Constraints**
   - Define limits such as word count, time period, or specific points to cover.
   - Example:
     - "Summarize the French Revolution in 150 words, focusing on its impact on European politics."

---

### 9. **Use Role-Playing**
   - Ask the model to adopt a specific role or perspective.
   - Example:
     - "You are a career counselor. Advise a recent college graduate on how to prepare for a job interview in the tech industry."

---

### 10. **Test and Iterate**
   - Experiment with different phrasings and structures to see what works best.
   - Example:
     - If the first prompt doesn’t yield the desired result, rephrase or add more details.

---

### Example Prompts for Claude 3.5 Sonnet:
1. **Creative Writing:**
   - "Write a 200-word science fiction story about a world where humans can upload their memories to the cloud. Include a moral dilemma about privacy."

2. **Analytical Task:**
   - "Analyze the strengths and weaknesses of remote work for businesses. Provide three benefits and three challenges, supported by examples."

3. **Educational Explanation:**
   - "Explain the concept of blockchain technology to a beginner. Use simple language and analogies to make it easy to understand."

4. **Problem-Solving:**
   - "You are a project manager. Outline a step-by-step plan to resolve a team conflict over resource allocation."

