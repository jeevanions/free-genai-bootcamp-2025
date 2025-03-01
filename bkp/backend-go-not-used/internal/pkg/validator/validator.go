package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates a struct using validator tags
func Validate(i interface{}) error {
	return validate.Struct(i)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}

type Word struct {
	Italian         string
	English         string
	PartsOfSpeech   string
	DifficultyLevel int
}

func ValidateWord(word *Word) error {
	if word.Italian == "" {
		return errors.New("italian word is required")
	}
	if word.English == "" {
		return errors.New("english translation is required")
	}
	if word.PartsOfSpeech == "" {
		return errors.New("parts of speech is required")
	}
	if word.DifficultyLevel < 1 || word.DifficultyLevel > 5 {
		return errors.New("difficulty level must be between 1 and 5")
	}
	return nil
}
