package models

import "time"

// Valid values for review types
var ValidReviewTypes = []string{"translation", "pronunciation"}

// WordReviewDetails contains additional details about a word review
type WordReviewDetails struct {
	ItalianWord  string `json:"italian_word"`
	EnglishWord  string `json:"english_word"`
	WordCategory string `json:"word_category"`
	Difficulty   int    `json:"difficulty"`
}

// WordReviewItem represents a word review in the database
type WordReviewItem struct {
	ID             int64            `db:"id" json:"id"`
	WordID         int64            `db:"word_id" json:"word_id"`
	StudySessionID int64            `db:"study_session_id" json:"study_session_id"`
	Correct        bool             `db:"correct" json:"correct"`
	ReviewType     string           `db:"review_type" json:"review_type"`
	CreatedAt      time.Time        `db:"created_at" json:"created_at"`
	WordDetails    WordReviewDetails `json:"word_details"`
}

// TableName returns the name of the table for the WordReviewItem model
func (WordReviewItem) TableName() string {
	return "word_review_items"
}

// WordReviewStats represents statistics for a word's review history
type WordReviewStats struct {
	WordID          int64   `json:"word_id"`
	TotalAttempts   int     `json:"total_attempts"`
	CorrectAttempts int     `json:"correct_attempts"`
	AccuracyRate    float64 `json:"accuracy_rate"`
	LastReviewedAt  string  `json:"last_reviewed_at"`
}

// CalculateAccuracy calculates the accuracy rate for the word review stats
func (w *WordReviewStats) CalculateAccuracy() {
	if w.TotalAttempts > 0 {
		w.AccuracyRate = float64(w.CorrectAttempts) / float64(w.TotalAttempts) * 100
	}
}
