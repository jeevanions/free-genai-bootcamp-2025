package models

import "time"

type CreateStudySessionRequest struct {
	GroupID         int64 `json:"group_id" binding:"required"`
	StudyActivityID int64 `json:"study_activity_id" binding:"required"`
	TotalWords      int   `json:"total_words" binding:"required,min=1"`
	CorrectWords    int   `json:"correct_words" binding:"required,min=0"`
	DurationSeconds int   `json:"duration_seconds" binding:"required,min=1"`
}

type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	StudyActivityID int64     `json:"study_activity_id"`
	TotalWords      int       `json:"total_words"`
	CorrectWords    int       `json:"correct_words"`
	DurationSeconds int       `json:"duration_seconds"`
	CreatedAt       time.Time `json:"created_at"`
}

type GroupStats struct {
	TotalSessions   int64   `json:"total_sessions"`
	TotalWords      int     `json:"total_words"`
	CorrectWords    int     `json:"correct_words"`
	AverageAccuracy float64 `json:"average_accuracy"`
}
