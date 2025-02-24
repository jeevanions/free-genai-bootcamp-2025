package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

type LLMServiceInterface interface {
	GenerateWords(category string) (*models.GenerateWordsResponse, error)
	GetGroupByID(id int64) (*models.GroupResponse, error)
	CreateWord(word *models.WordResponse) (int64, error)
	AddWordToGroup(wordID, groupID int64) error
	UpdateGroupWordsCount(groupID int64) error
}

type LLMService struct {
	repo repository.Repository
}

func NewLLMService(repo repository.Repository) *LLMService {
	return &LLMService{repo: repo}
}

func (s *LLMService) GenerateWords(category string) (*models.GenerateWordsResponse, error) {
	// Groq API endpoint and key
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY environment variable not set")
	}

	// Enhanced prompt for better word generation
	prompt := fmt.Sprintf(`Generate 10 Italian words for the thematic category: %s.
	For each word, provide:
	- The Italian word (with correct spelling and accents)
	- Accurate English translation
	- Detailed grammatical information including:
	  * Part of speech (noun, verb, adjective, etc.)
	  * Gender for nouns (masculine/feminine)
	  * Plural form for nouns
	  * Any irregular forms or important notes
	
	Format the response as a JSON array of objects. Each object should have this exact structure:
	[
		{
			"italian": "word",
			"english": "translation",
			"parts": {
				"type": "noun/verb/adjective",
				"gender": "masculine/feminine",
				"plural": "plural_form"
			}
		}
	]
	Do not include any explanations or additional text, only return the JSON array.`, category)

	// Groq API request configuration
	reqBody := map[string]interface{}{
		"model": "mixtral-8x7b-32768",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert Italian language teacher specializing in vocabulary.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"max_tokens":  1000,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract the generated words from the response
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("invalid response format")
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid message format")
	}

	content, ok := message["content"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid content format")
	}

	// Clean and parse the JSON content
	content = cleanJSONString(content)

	var words []models.WordResponse
	if err := json.Unmarshal([]byte(content), &words); err != nil {
		// Try to extract JSON array if content contains additional text
		if jsonStart := strings.Index(content, "["); jsonStart != -1 {
			if jsonEnd := strings.LastIndex(content, "]"); jsonEnd != -1 && jsonEnd > jsonStart {
				content = content[jsonStart : jsonEnd+1]
				if err := json.Unmarshal([]byte(content), &words); err != nil {
					return nil, fmt.Errorf("failed to parse extracted JSON array: %v", err)
				}
			}
		}
		if len(words) == 0 {
			return nil, fmt.Errorf("failed to parse generated words: %v", err)
		}
	}

	return &models.GenerateWordsResponse{Words: words}, nil
}

// cleanJSONString removes common issues in LLM-generated JSON
func cleanJSONString(s string) string {
	// Remove any markdown code block markers
	s = strings.ReplaceAll(s, "```json", "")
	s = strings.ReplaceAll(s, "```", "")

	// Remove any explanatory text before the first '['
	if idx := strings.Index(s, "["); idx != -1 {
		s = s[idx:]
	}

	// Remove any text after the last ']'
	if idx := strings.LastIndex(s, "]"); idx != -1 {
		s = s[:idx+1]
	}

	return strings.TrimSpace(s)
}

func (s *LLMService) GetGroupByID(id int64) (*models.GroupResponse, error) {
	group, err := s.repo.GetGroupByID(id)
	if err != nil {
		return nil, err
	}
	return &models.GroupResponse{
		ID:        group.ID,
		Name:      group.Name,
		WordCount: group.Stats.TotalWordCount,
	}, nil
}

func (s *LLMService) CreateWord(word *models.WordResponse) (int64, error) {
	return s.repo.CreateWord(word)
}

func (s *LLMService) AddWordToGroup(wordID, groupID int64) error {
	return s.repo.AddWordToGroup(wordID, groupID)
}

func (s *LLMService) UpdateGroupWordsCount(groupID int64) error {
	return s.repo.UpdateGroupWordsCount(groupID)
}
