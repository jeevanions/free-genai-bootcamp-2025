package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Word struct {
	ID              int64           `json:"id"`
	Italian         string          `json:"italian"`
	English         string          `json:"english"`
	PartsOfSpeech   string          `json:"parts_of_speech"`
	Gender          sql.NullString  `json:"gender,omitempty"`
	Number          sql.NullString  `json:"number,omitempty"`
	DifficultyLevel int             `json:"difficulty_level"`
	VerbConjugation json.RawMessage `json:"verb_conjugation,omitempty"`
	Notes           sql.NullString  `json:"notes,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

type Group struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	DifficultyLevel int       `json:"difficulty_level"`
	Category        string    `json:"category"`
	CreatedAt       time.Time `json:"created_at"`
}

type StudyActivity struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	RequiresAudio   bool      `json:"requires_audio"`
	DifficultyLevel int       `json:"difficulty_level"`
	Instructions    string    `json:"instructions"`
	CreatedAt       time.Time `json:"created_at"`
}

type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	StudyActivityID int64     `json:"study_activity_id"`
	TotalWords      int       `json:"total_words"`
	CorrectWords    int       `json:"correct_words"`
	DurationSeconds int       `json:"duration_seconds"`
	CreatedAt       time.Time `json:"created_at"`
}
