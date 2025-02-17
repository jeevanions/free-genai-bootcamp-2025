package models

import "time"

type WordReviewItem struct {
	ID             int64              `json:"id"`
	WordID         int64              `json:"word_id"`
	StudySessionID int64              `json:"study_session_id"`
	Correct        bool               `json:"correct"`
	CreatedAt      time.Time          `json:"created_at"`
	WordDetails    *WordReviewDetails `json:"word_details,omitempty"`
}

type WordReviewDetails struct {
	Italian string `json:"italian"`
	English string `json:"english"`
}

type WordStats struct {
	TotalReviews   int64   `json:"total_reviews"`
	CorrectReviews int64   `json:"correct_reviews"`
	Accuracy       float64 `json:"accuracy"`
}

func (w *WordStats) CalculateAccuracy() {
	if w.TotalReviews > 0 {
		w.Accuracy = float64(w.CorrectReviews) / float64(w.TotalReviews) * 100
	}
}
