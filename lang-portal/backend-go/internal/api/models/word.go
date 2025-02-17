package models

import "time"

// Word represents a vocabulary word in the Italian learning system
// @Description Word model containing Italian vocabulary information
type Word struct {
	ID              int64     `json:"id"`
	Italian         string    `json:"italian"`
	English         string    `json:"english"`
	PartsOfSpeech   string    `json:"parts_of_speech"`
	Gender          *string   `json:"gender,omitempty"`
	Number          *string   `json:"number,omitempty"`
	DifficultyLevel int       `json:"difficulty_level"`
	VerbConjugation *string   `json:"verb_conjugation,omitempty"`
	Notes           *string   `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreateWordRequest represents the request body for creating a new word
// @Description Request model for creating a new word
type CreateWordRequest struct {
	Italian         string  `json:"italian" binding:"required"`
	English         string  `json:"english" binding:"required"`
	PartsOfSpeech   string  `json:"parts_of_speech" binding:"required"`
	Gender          *string `json:"gender,omitempty"`
	Number          *string `json:"number,omitempty"`
	DifficultyLevel int     `json:"difficulty_level" binding:"required,min=1,max=5"`
	VerbConjugation *string `json:"verb_conjugation,omitempty"`
	Notes           *string `json:"notes,omitempty"`
}

// UpdateWordRequest represents the request body for updating a word
// @Description Request model for updating an existing word
type UpdateWordRequest struct {
	Italian         *string `json:"italian,omitempty" example:"ciao"`
	English         *string `json:"english,omitempty" example:"hello"`
	PartsOfSpeech   *string `json:"parts_of_speech,omitempty" example:"interjection"`
	Gender          *string `json:"gender,omitempty" example:"masculine"`
	Number          *string `json:"number,omitempty" example:"singular"`
	DifficultyLevel *int    `json:"difficulty_level,omitempty" binding:"omitempty,min=1,max=5" example:"1"`
	VerbConjugation *string `json:"verb_conjugation,omitempty"`
	Notes           *string `json:"notes,omitempty" example:"Common greeting"`
}
