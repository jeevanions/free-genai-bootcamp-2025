package models

import "time"

// User represents a user in the database
type User struct {
	ID           int64     `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"` // Never expose password hash in JSON
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the name of the table for the User model
func (User) TableName() string {
	return "users"
}

// Sanitize removes sensitive fields before returning to client
func (u *User) Sanitize() {
	u.PasswordHash = ""
}
