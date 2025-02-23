package models

// WordResponse represents a word with its translations and grammatical details
// swagger:model
type WordResponse struct {
	// The unique identifier of the word
	// required: true
	ID int64 `json:"id" example:"1"`
	// The Italian word
	// required: true
	Italian string `json:"italian" example:"sorella"`
	// The English translation
	// required: true
	English string `json:"english" example:"sister"`
	// Grammatical details like type, gender, plural form
	// required: true
	Parts map[string]interface{} `json:"parts"`
	// Number of times the word was correctly answered
	// required: true
	CorrectCount int `json:"correct_count" example:"5"`
	// Number of times the word was incorrectly answered
	// required: true
	WrongCount int `json:"wrong_count" example:"2"`
}

type WordListResponse struct {
	Items      []WordResponse     `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

// WordReviewRequest represents a request to review a word in a study session
type WordReviewRequest struct {
	Correct bool `json:"correct" binding:"required"`
}

// WordReviewResponse represents a response to a word review request
type WordReviewResponse struct {
	Success bool  `json:"success"`
	WordID  int64 `json:"word_id"`
}

type ImportWordsRequest struct {
	GroupID int64    `json:"group_id"`
	Words   []string `json:"words"`
}

type ImportWordsResponse struct {
	ImportedCount int `json:"imported_count"`
}
