package models

type WordResponse struct {
	ID           int64                  `json:"id"`
	Italian      string                 `json:"italian"`
	English      string                 `json:"english"`
	Parts        map[string]interface{} `json:"parts"`
	CorrectCount int                    `json:"correct_count"`
	WrongCount   int                    `json:"wrong_count"`
}

type WordListResponse struct {
	Items      []WordResponse `json:"items"`
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

