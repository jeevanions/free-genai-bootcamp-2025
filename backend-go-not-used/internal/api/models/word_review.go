package models

import (
	"fmt"
	"time"
)

// Valid values for review types
var ValidReviewTypes = []string{"translation", "pronunciation"}

// CreateWordReviewRequest represents the request to create a new word review
type CreateWordReviewRequest struct {
	WordID         int64  `json:"word_id" binding:"required" example:"1"`
	StudySessionID int64  `json:"study_session_id" binding:"required" example:"1"`
	Correct        bool   `json:"correct" binding:"required" example:"true"`
	ReviewType     string `json:"review_type" binding:"required" example:"translation"`
}

// Validate validates the CreateWordReviewRequest
func (r *CreateWordReviewRequest) Validate() error {
	// Validate review type
	validType := false
	for _, t := range ValidReviewTypes {
		if r.ReviewType == t {
			validType = true
			break
		}
	}
	if !validType {
		return fmt.Errorf("invalid review_type: must be one of %v", ValidReviewTypes)
	}

	return nil
}

// WordReviewResponse represents a word review in the API response
type WordReviewResponse struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	ReviewType     string    `json:"review_type"`
	CreatedAt      time.Time `json:"created_at"`
	Word           *Word     `json:"word,omitempty"`
}

// WordReviewStats represents statistics for a word's review history
type WordReviewStats struct {
	WordID          int64   `json:"word_id"`
	TotalAttempts   int     `json:"total_attempts"`
	CorrectAttempts int     `json:"correct_attempts"`
	AccuracyRate    float64 `json:"accuracy_rate"`
	LastReviewedAt  string  `json:"last_reviewed_at"`
}
