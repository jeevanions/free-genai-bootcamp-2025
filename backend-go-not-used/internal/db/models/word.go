package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Word represents a word in the database
type Word struct {
	ID              int64           `db:"id"`
	Italian         string          `db:"italian"`
	English         string          `db:"english"`
	PartsOfSpeech   string          `db:"parts_of_speech"`
	Gender          sql.NullString  `db:"gender"`
	Number          sql.NullString  `db:"number"`
	DifficultyLevel int             `db:"difficulty_level"`
	VerbConjugation json.RawMessage `db:"verb_conjugation"`
	Notes           sql.NullString  `db:"notes"`
	CreatedAt       time.Time       `db:"created_at"`
}

// TableName returns the name of the table for the Word model
func (Word) TableName() string {
	return "words"
}
