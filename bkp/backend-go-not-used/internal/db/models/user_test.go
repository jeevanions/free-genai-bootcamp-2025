package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_TableName(t *testing.T) {
	user := User{}
	assert.Equal(t, "users", user.TableName())
}

func TestUser_Fields(t *testing.T) {
	now := time.Now()
	user := User{
		ID:           1,
		Username:     "john_doe",
		Email:        "john@example.com",
		PasswordHash: "hashed_password",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Test field types and values
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "john_doe", user.Username)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "hashed_password", user.PasswordHash)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)
}

func TestUser_Sanitize(t *testing.T) {
	user := User{
		ID:           1,
		Username:     "john_doe",
		Email:        "john@example.com",
		PasswordHash: "hashed_password",
	}

	// Test that Sanitize removes sensitive fields
	user.Sanitize()
	assert.Empty(t, user.PasswordHash)

	// Other fields should remain unchanged
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "john_doe", user.Username)
	assert.Equal(t, "john@example.com", user.Email)
}
