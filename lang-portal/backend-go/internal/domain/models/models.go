package models

import (
	"encoding/json"
	"time"
)

type Word struct {
	ID       int64           `json:"id"`
	Italian  string         `json:"italian"`
	English  string         `json:"english"`
	Parts    json.RawMessage `json:"parts,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

type Group struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	WordsCount int       `json:"words_count"`
	CreatedAt  time.Time `json:"created_at"`
}

type StudyActivity struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Description  string    `json:"description"`
	LaunchURL    *string   `json:"launch_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	StudyActivityID int64     `json:"study_activity_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type WordReviewItem struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

// StrPtr helper creates string pointer
func StrPtr(s string) *string {
	return &s
}
