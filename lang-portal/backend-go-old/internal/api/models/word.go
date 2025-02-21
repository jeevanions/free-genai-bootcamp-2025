package models

import (
	"fmt"
	"time"
)

// Valid values for parts of speech
var ValidPartsOfSpeech = []string{"noun", "verb", "adjective", "adverb", "preposition", "conjunction", "interjection"}

// Valid values for gender
var ValidGenders = []string{"masculine", "feminine", "neuter"}

// Valid values for number
var ValidNumbers = []string{"singular", "plural"}

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

// Validate validates the CreateWordRequest
func (r *CreateWordRequest) Validate() error {
	// Validate parts of speech
	validPOS := false
	for _, pos := range ValidPartsOfSpeech {
		if r.PartsOfSpeech == pos {
			validPOS = true
			break
		}
	}
	if !validPOS {
		return fmt.Errorf("invalid parts_of_speech: must be one of %v", ValidPartsOfSpeech)
	}

	// Validate difficulty level
	if r.DifficultyLevel < 1 || r.DifficultyLevel > 5 {
		return fmt.Errorf("difficulty_level must be between 1 and 5")
	}

	// Validate gender (only for nouns)
	if r.Gender != nil {
		if r.PartsOfSpeech != "noun" {
			return fmt.Errorf("gender can only be set for nouns")
		}
		validGender := false
		for _, g := range ValidGenders {
			if *r.Gender == g {
				validGender = true
				break
			}
		}
		if !validGender {
			return fmt.Errorf("invalid gender: must be one of %v", ValidGenders)
		}
	}

	// Validate number (only for nouns)
	if r.Number != nil {
		if r.PartsOfSpeech != "noun" {
			return fmt.Errorf("number can only be set for nouns")
		}
		validNumber := false
		for _, n := range ValidNumbers {
			if *r.Number == n {
				validNumber = true
				break
			}
		}
		if !validNumber {
			return fmt.Errorf("invalid number: must be one of %v", ValidNumbers)
		}
	}

	// Validate verb conjugation (only for verbs)
	if r.VerbConjugation != nil && r.PartsOfSpeech != "verb" {
		return fmt.Errorf("verb_conjugation can only be set for verbs")
	}

	return nil
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

// Validate validates the UpdateWordRequest
func (r *UpdateWordRequest) Validate() error {
	// If parts of speech is being updated, validate it
	if r.PartsOfSpeech != nil {
		validPOS := false
		for _, pos := range ValidPartsOfSpeech {
			if *r.PartsOfSpeech == pos {
				validPOS = true
				break
			}
		}
		if !validPOS {
			return fmt.Errorf("invalid parts_of_speech: must be one of %v", ValidPartsOfSpeech)
		}
	}

	// If difficulty level is being updated, validate it
	if r.DifficultyLevel != nil {
		if *r.DifficultyLevel < 1 || *r.DifficultyLevel > 5 {
			return fmt.Errorf("difficulty_level must be between 1 and 5")
		}
	}

	// If gender is being updated, validate it
	if r.Gender != nil {
		if r.PartsOfSpeech != nil && *r.PartsOfSpeech != "noun" {
			return fmt.Errorf("gender can only be set for nouns")
		}
		validGender := false
		for _, g := range ValidGenders {
			if *r.Gender == g {
				validGender = true
				break
			}
		}
		if !validGender {
			return fmt.Errorf("invalid gender: must be one of %v", ValidGenders)
		}
	}

	// If number is being updated, validate it
	if r.Number != nil {
		if r.PartsOfSpeech != nil && *r.PartsOfSpeech != "noun" {
			return fmt.Errorf("number can only be set for nouns")
		}
		validNumber := false
		for _, n := range ValidNumbers {
			if *r.Number == n {
				validNumber = true
				break
			}
		}
		if !validNumber {
			return fmt.Errorf("invalid number: must be one of %v", ValidNumbers)
		}
	}

	// If verb conjugation is being updated, validate it
	if r.VerbConjugation != nil {
		if r.PartsOfSpeech != nil && *r.PartsOfSpeech != "verb" {
			return fmt.Errorf("verb_conjugation can only be set for verbs")
		}
	}

	return nil
}
