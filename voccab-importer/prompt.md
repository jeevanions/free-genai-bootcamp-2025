# Vocabular Importer - React Application Specification

## Overview
Create a single-page React application that allows users to:
1. Select or create thematic categories for Italian vocabulary
2. Generate Italian-English word pairs for the selected category
3. Review and edit the generated words
4. Import the final word list to the backend database

## Technical Requirements
- Frontend: React with TypeScript
- Build Tool: Vite
- Styling: TailwindCSS
- State Management: React Hooks
- API Communication: Your choice of modern HTTP client library
- Backend Base URL: `http://localhost:8080/`

## UI Components

### Header
- Professional logo
- Title "Vocabular Importer"
- Styled with appropriate colors for visual appeal

### Main Content
1. Category Selection Area
   - Searchable dropdown for thematic categories
   - "Generate" button (enabled only when category is selected)
2. Word Review Area
   - JSON format textarea
   - "Import" button (enabled only when valid JSON is present)

## Detailed Feature Specifications

### Category Dropdown Functionality
1. Initial Load:
   - Endpoint: GET `/api/groups`
   - Response Schema:
   ```typescript
   interface GroupsResponse {
     items: Array<{
       id: number;
       name: string;
       word_count: number;
     }>;
     pagination: {
       current_page: number;
       items_per_page: number;
       total_items: number;
       total_pages: number;
     };
   }
   ```

2. Category Creation:
   - Trigger: User types non-existing category and presses Enter
   - Endpoint: POST `/api/groups`
   - Request Body: `{ name: string }`
   - Response: Returns created group ID

3. Interaction Features:
   - Filterable by typing
   - Single selection only
   - Clear selection option
   - Enter key selection support
   - Reselection of same category allowed

### Word Generation Process
1. Trigger: "Generate" button click
2. Endpoint: POST `/api/words/llm/generate-words`
3. Request Body:
   ```typescript
   interface GenerateWordsRequest {
     category: string;  // selected category name
   }
   ```
4. Response Schema:
   ```typescript
   interface GenerateWordsResponse {
     words: Array<{
       id: number;
       italian: string;
       english: string;
       parts: {
         gender: string;
         plural: string;
         type: string;
       };
       correct_count: number;
       wrong_count: number;
     }>;
   }
   ```

### Import Process
1. Validation:
   - Ensure textarea contains valid JSON
   - Display error if JSON is malformed
2. Submission:
   - Endpoint: POST `/api/words/import`
   - Request Body: JSON from textarea (validated)
   - Error Handling: Display appropriate error messages

## Documentation Requirements
Create a README.md including:
1. Project overview
2. Setup instructions for:
   - Development environment
   - Production deployment
3. Feature documentation
4. API endpoint documentation
```

Key Improvements Made:
1. Structured the content hierarchically
2. Added clear TypeScript interfaces for API responses
3. Separated concerns into distinct sections
4. Made the validation requirements explicit
5. Added error handling specifications
6. Organized the UI components more clearly
7. Made the API interaction flow more explicit


