# Vocabular Importer

A React application for managing and importing Italian vocabulary with thematic categories.

## Features

- Select or create thematic categories for Italian vocabulary
- Generate Italian-English word pairs for selected categories
- Review and edit generated words in JSON format
- Import finalized word lists to the backend database

## Tech Stack

- React with TypeScript
- Vite for build tooling
- TailwindCSS for styling
- Axios for API communication
- React Select for advanced dropdown functionality

## Development Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run dev
   ```

3. Open http://localhost:5173 in your browser

## Production Deployment

1. Build the application:
   ```bash
   npm run build
   ```

2. The built files will be in the `dist` directory, ready to be served by any static file server

## API Endpoints

### Groups

- `GET /api/groups`
  - Retrieves list of vocabulary groups
  - Returns: List of groups with pagination

- `POST /api/groups`
  - Creates a new vocabulary group
  - Body: `{ name: string }`
  - Returns: Created group details

### Words

- `POST /api/words/llm/generate-words`
  - Generates words for a category
  - Body: `{ category: string }`
  - Returns: List of generated words

- `POST /api/words/import`
  - Imports words into the system
  - Body: Array of word objects
  - Returns: Success confirmation