package models

import "time"

// WordStats represents statistics for a word's review history
type WordStats struct {
	WordID          int64     `json:"word_id"`
	TotalAttempts   int       `json:"total_attempts"`
	CorrectAttempts int       `json:"correct_attempts"`
	LastReviewedAt  time.Time `json:"last_reviewed_at"`
	MasteryLevel    float64   `json:"mastery_level"`
	SuccessRate     float64   `json:"success_rate"`
	NextReviewDue   time.Time `json:"next_review_due"`
}

// CalculateAccuracy calculates and updates the success rate for the word stats
func (s *WordStats) CalculateAccuracy() {
	if s.TotalAttempts > 0 {
		s.SuccessRate = float64(s.CorrectAttempts) / float64(s.TotalAttempts) * 100
	} else {
		s.SuccessRate = 0
	}
}
