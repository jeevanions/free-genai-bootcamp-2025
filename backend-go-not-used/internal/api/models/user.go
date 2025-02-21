package models

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`)
)

// RegisterUserRequest represents the request to register a new user
type RegisterUserRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Email    string `json:"email" binding:"required" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"SecurePass123!"`
}

// Validate validates the RegisterUserRequest
func (r *RegisterUserRequest) Validate() error {
	// Validate username
	if !usernameRegex.MatchString(r.Username) {
		return fmt.Errorf("invalid username: must be 3-30 characters long and contain only letters, numbers, underscores, and hyphens")
	}

	// Validate email
	if !emailRegex.MatchString(r.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Validate password strength
	if err := validatePassword(r.Password); err != nil {
		return err
	}

	return nil
}

// LoginRequest represents the request to login
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"SecurePass123!"`
}

// UpdateUserRequest represents the request to update user details
type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" example:"newemail@example.com"`
	Password *string `json:"password,omitempty" example:"NewSecurePass123!"`
}

// Validate validates the UpdateUserRequest
func (r *UpdateUserRequest) Validate() error {
	if r.Email != nil {
		if !emailRegex.MatchString(*r.Email) {
			return fmt.Errorf("invalid email format")
		}
	}

	if r.Password != nil {
		if err := validatePassword(*r.Password); err != nil {
			return err
		}
	}

	return nil
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// validatePassword checks if the password meets security requirements
func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	return nil
}
