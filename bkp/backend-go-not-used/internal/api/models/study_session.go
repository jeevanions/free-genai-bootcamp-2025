package models

import (
	"fmt"
	"time"
)

type SessionStats struct {
	TotalSessions int64 `json:"total_sessions"`
	TotalWords    int   `json:"total_words"`
	TotalCorrect  int   `json:"total_correct"`
}

type CreateStudySessionRequest struct {
	GroupID         int64 `json:"group_id" binding:"required"`
	StudyActivityID int64 `json:"study_activity_id" binding:"required"`
	TotalWords      int   `json:"total_words" binding:"required,min=1"`
	CorrectWords    int   `json:"correct_words" binding:"required,min=0"`
	DurationSeconds int   `json:"duration_seconds" binding:"required,min=1"`
}

// Validate validates the CreateStudySessionRequest
func (r *CreateStudySessionRequest) Validate() error {
	// Validate IDs
	if r.GroupID <= 0 {
		return fmt.Errorf("group_id must be positive")
	}
	if r.StudyActivityID <= 0 {
		return fmt.Errorf("study_activity_id must be positive")
	}

	// Validate word counts
	if r.TotalWords < 1 {
		return fmt.Errorf("total_words must be at least 1")
	}
	if r.CorrectWords < 0 {
		return fmt.Errorf("correct_words must be non-negative")
	}
	if r.CorrectWords > r.TotalWords {
		return fmt.Errorf("correct_words cannot be greater than total_words")
	}

	// Validate duration
	if r.DurationSeconds < 1 {
		return fmt.Errorf("duration_seconds must be at least 1")
	}
	if r.DurationSeconds > 24*60*60 { // Max 24 hours
		return fmt.Errorf("duration_seconds cannot exceed 24 hours")
	}

	return nil
}

// Rename to avoid conflict
type StudySessionResponse struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	StudyActivityID int64     `json:"study_activity_id"`
	TotalWords      int       `json:"total_words"`
	CorrectWords    int       `json:"correct_words"`
	DurationSeconds int       `json:"duration_seconds"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	CreatedAt       time.Time `json:"created_at"`
}

// LastStudySessionResponse represents the response for the last study session endpoint
type LastStudySessionResponse struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	GroupName       string    `json:"group_name"`
	ActivityID      int64     `json:"activity_id"`
	ActivityName    string    `json:"activity_name"`
	TotalWords      int       `json:"total_words"`
	CorrectWords    int       `json:"correct_words"`
	DurationSeconds int       `json:"duration_seconds"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
}

type GroupStats struct {
	TotalSessions   int64   `json:"total_sessions"`
	TotalWords      int     `json:"total_words"`
	CorrectWords    int     `json:"correct_words"`
	AverageAccuracy float64 `json:"average_accuracy"`
}
