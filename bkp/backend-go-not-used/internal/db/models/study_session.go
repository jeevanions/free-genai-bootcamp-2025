package models

import (
	"time"
)

// StudySession represents a study session in the database
type StudySession struct {
	ID              int64     `db:"id"`
	GroupID         int64     `db:"group_id"`
	StudyActivityID int64     `db:"study_activity_id"`
	TotalWords      int       `db:"total_words"`
	CorrectWords    int       `db:"correct_words"`
	DurationSeconds int       `db:"duration_seconds"`
	StartTime       time.Time `db:"start_time"`
	EndTime         time.Time `db:"end_time"`
	CreatedAt       time.Time `db:"created_at"`
}

// TableName returns the name of the table for the StudySession model
func (StudySession) TableName() string {
	return "study_sessions"
}
