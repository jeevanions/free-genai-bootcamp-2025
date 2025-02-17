package models

import "time"

type CreateWordReviewRequest struct {
	WordID         int64 `json:"word_id" binding:"required"`
	StudySessionID int64 `json:"study_session_id" binding:"required"`
	Correct        bool  `json:"correct" binding:"required"`
}

type WordReviewItem struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	WordDetails    *Word     `json:"word_details,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type WordStats struct {
	TotalReviews   int64   `json:"total_reviews"`
	CorrectReviews int64   `json:"correct_reviews"`
	Accuracy       float64 `json:"accuracy"`
}
