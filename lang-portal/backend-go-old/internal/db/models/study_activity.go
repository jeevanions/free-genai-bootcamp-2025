package models

import (
	"time"
)

// Valid values for activity types
var ValidActivityTypes = []string{"vocabulary", "grammar", "pronunciation"}

// Valid values for activity categories
var ValidActivityCategories = []string{"basics", "travel", "food", "business", "culture", "daily_life"}

// StudyActivity represents a learning activity in the database
type StudyActivity struct {
	ID              int64     `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	Type            string    `db:"type" json:"type"`
	RequiresAudio   bool      `db:"requires_audio" json:"requires_audio"`
	DifficultyLevel int       `db:"difficulty_level" json:"difficulty_level"`
	Instructions    string    `db:"instructions" json:"instructions"`
	ThumbnailURL    string    `db:"thumbnail_url" json:"thumbnail_url,omitempty"`
	Category        string    `db:"category" json:"category"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

// TableName returns the name of the table for the StudyActivity model
func (StudyActivity) TableName() string {
	return "study_activities"
}
