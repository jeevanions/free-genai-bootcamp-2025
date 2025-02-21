package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request RegisterUserRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: RegisterUserRequest{
				Username: "john_doe",
				Email:    "john@example.com",
				Password: "SecurePass123!",
			},
			wantErr: false,
		},
		{
			name: "invalid username - too short",
			request: RegisterUserRequest{
				Username: "jo",
				Email:    "john@example.com",
				Password: "SecurePass123!",
			},
			wantErr: true,
		},
		{
			name: "invalid username - special characters",
			request: RegisterUserRequest{
				Username: "john@doe",
				Email:    "john@example.com",
				Password: "SecurePass123!",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			request: RegisterUserRequest{
				Username: "john_doe",
				Email:    "invalid-email",
				Password: "SecurePass123!",
			},
			wantErr: true,
		},
		{
			name: "weak password - too short",
			request: RegisterUserRequest{
				Username: "john_doe",
				Email:    "john@example.com",
				Password: "Weak1!",
			},
			wantErr: true,
		},
		{
			name: "weak password - no uppercase",
			request: RegisterUserRequest{
				Username: "john_doe",
				Email:    "john@example.com",
				Password: "securepass123!",
			},
			wantErr: true,
		},
		{
			name: "weak password - no lowercase",
			request: RegisterUserRequest{
				Username: "john_doe",
				Email:    "john@example.com",
				Password: "SECUREPASS123!",
			},
			wantErr: true,
		},
		{
			name: "weak password - no number",
			request: RegisterUserRequest{
				Username: "john_doe",
				Email:    "john@example.com",
				Password: "SecurePass!!!",
			},
			wantErr: true,
		},
		{
			name: "weak password - no special character",
			request: RegisterUserRequest{
				Username: "john_doe",
				Email:    "john@example.com",
				Password: "SecurePass123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateUserRequest_Validate(t *testing.T) {
	validEmail := "newemail@example.com"
	invalidEmail := "invalid-email"
	validPassword := "NewSecurePass123!"
	weakPassword := "weak"

	tests := []struct {
		name    string
		request UpdateUserRequest
		wantErr bool
	}{
		{
			name: "valid request - email only",
			request: UpdateUserRequest{
				Email: &validEmail,
			},
			wantErr: false,
		},
		{
			name: "valid request - password only",
			request: UpdateUserRequest{
				Password: &validPassword,
			},
			wantErr: false,
		},
		{
			name: "valid request - both fields",
			request: UpdateUserRequest{
				Email:    &validEmail,
				Password: &validPassword,
			},
			wantErr: false,
		},
		{
			name: "invalid email format",
			request: UpdateUserRequest{
				Email: &invalidEmail,
			},
			wantErr: true,
		},
		{
			name: "weak password",
			request: UpdateUserRequest{
				Password: &weakPassword,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
