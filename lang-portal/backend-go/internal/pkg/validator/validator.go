package validator

import (
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
