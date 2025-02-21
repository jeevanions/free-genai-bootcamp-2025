package models

import "time"

// WordProgress represents a user's progress with a specific word
type WordProgress struct {
	ID              int64     `db:"id" json:"id"`
	UserID          int64     `db:"user_id" json:"user_id"`
	WordID          int64     `db:"word_id" json:"word_id"`
	MasteryLevel    int       `db:"mastery_level" json:"mastery_level"`
	NextReviewAt    time.Time `db:"next_review_at" json:"next_review_at"`
	TotalAttempts   int       `db:"total_attempts" json:"total_attempts"`
	CorrectAttempts int       `db:"correct_attempts" json:"correct_attempts"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the name of the table for the WordProgress model
func (WordProgress) TableName() string {
	return "word_progress"
}

// CalculateAccuracy calculates the accuracy rate for the word progress
func (w *WordProgress) CalculateAccuracy() float64 {
	if w.TotalAttempts > 0 {
		return float64(w.CorrectAttempts) / float64(w.TotalAttempts) * 100
	}
	return 0
}

// UpdateMasteryLevel updates the mastery level based on accuracy and attempts
func (w *WordProgress) UpdateMasteryLevel() {
	accuracy := w.CalculateAccuracy()
	minAttempts := 5 // Minimum attempts required for mastery level increase

	switch {
	case w.TotalAttempts < minAttempts:
		// Keep current mastery level if not enough attempts
		return
	case accuracy >= 90:
		w.MasteryLevel = 5
	case accuracy >= 80:
		w.MasteryLevel = 4
	case accuracy >= 70:
		w.MasteryLevel = 3
	case accuracy >= 60:
		w.MasteryLevel = 2
	case accuracy >= 50:
		w.MasteryLevel = 1
	default:
		w.MasteryLevel = 0
	}
}

// UpdateNextReviewTime calculates the next review time based on mastery level
func (w *WordProgress) UpdateNextReviewTime() {
	now := time.Now()
	switch w.MasteryLevel {
	case 5:
		w.NextReviewAt = now.AddDate(0, 1, 0)  // 1 month
	case 4:
		w.NextReviewAt = now.AddDate(0, 0, 14) // 2 weeks
	case 3:
		w.NextReviewAt = now.AddDate(0, 0, 7)  // 1 week
	case 2:
		w.NextReviewAt = now.AddDate(0, 0, 3)  // 3 days
	case 1:
		w.NextReviewAt = now.AddDate(0, 0, 1)  // 1 day
	default:
		w.NextReviewAt = now // Review immediately
	}
}
