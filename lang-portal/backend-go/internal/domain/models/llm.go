package models

// GenerateWordsRequest represents a request to generate words for a thematic category
// swagger:model
type GenerateWordsRequest struct {
	// The thematic category for word generation (e.g., "family members", "food", etc.)
	// required: true
	Category string `json:"category" binding:"required" example:"family members"`
}

// GenerateWordsResponse represents the response from the LLM with generated words
// swagger:model
type GenerateWordsResponse struct {
	// List of generated Italian words with translations and grammatical details
	// required: true
	Words []WordResponse `json:"words"`
}

// AddWordsToGroupRequest represents a request to add words to a group
// swagger:model
type AddWordsToGroupRequest struct {
	// List of words to add to the group
	// required: true
	Words []WordResponse `json:"words" binding:"required"`
}

// AddWordsToGroupResponse represents the response after adding words to a group
// swagger:model
type AddWordsToGroupResponse struct {
	// Whether the operation was successful
	// required: true
	Success bool `json:"success" example:"true"`
	// Number of words successfully added to the group
	// required: true
	WordsAdded int `json:"words_added" example:"5"`
}
