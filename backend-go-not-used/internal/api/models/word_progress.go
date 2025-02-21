package models

import (
	"fmt"
	"time"
)

// WordProgressResponse represents a word progress in the API response
type WordProgressResponse struct {
	ID              int64     `json:"id"`
	WordID          int64     `json:"word_id"`
	MasteryLevel    int       `json:"mastery_level"`
	NextReviewAt    time.Time `json:"next_review_at"`
	TotalAttempts   int       `json:"total_attempts"`
	CorrectAttempts int       `json:"correct_attempts"`
	AccuracyRate    float64   `json:"accuracy_rate"`
	Word            *Word     `json:"word,omitempty"`
}

// UpdateWordProgressRequest represents the request to update word progress
type UpdateWordProgressRequest struct {
	MasteryLevel    *int       `json:"mastery_level,omitempty" binding:"omitempty,min=0,max=5" example:"3"`
	NextReviewAt    *time.Time `json:"next_review_at,omitempty"`
	TotalAttempts   *int       `json:"total_attempts,omitempty" binding:"omitempty,min=0"`
	CorrectAttempts *int       `json:"correct_attempts,omitempty" binding:"omitempty,min=0"`
}

// Validate validates the UpdateWordProgressRequest
func (r *UpdateWordProgressRequest) Validate() error {
	if r.CorrectAttempts != nil && r.TotalAttempts != nil {
		if *r.CorrectAttempts > *r.TotalAttempts {
			return fmt.Errorf("correct_attempts cannot be greater than total_attempts")
		}
	}
	return nil
}

// WordProgressSummary represents a summary of word progress for a user
type WordProgressSummary struct {
	TotalWords           int     `json:"total_words"`
	WordsInProgress      int     `json:"words_in_progress"`
	WordsMastered        int     `json:"words_mastered"`
	AverageAccuracy      float64 `json:"average_accuracy"`
	AverageMasteryLevel float64 `json:"average_mastery_level"`
}

// WordProgressFilter represents filters for retrieving word progress
type WordProgressFilter struct {
	MasteryLevel    *int       `json:"mastery_level,omitempty" form:"mastery_level" binding:"omitempty,min=0,max=5"`
	ReviewDue       *bool      `json:"review_due,omitempty" form:"review_due"`
	ReviewDueBefore *time.Time `json:"review_due_before,omitempty" form:"review_due_before"`
	Limit           *int       `json:"limit,omitempty" form:"limit" binding:"omitempty,min=1,max=100"`
	Offset          *int       `json:"offset,omitempty" form:"offset" binding:"omitempty,min=0"`
}
