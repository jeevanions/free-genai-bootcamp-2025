package models

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
