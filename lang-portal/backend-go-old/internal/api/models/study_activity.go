package models

import (
	"fmt"
	"time"
)

// Valid values for activity types
var ValidActivityTypes = []string{"vocabulary", "grammar", "pronunciation"}

// Valid values for activity categories
var ValidActivityCategories = []string{"basics", "travel", "food", "business", "culture", "daily_life"}

// CreateStudyActivityRequest represents the request to create a new study activity
type CreateStudyActivityRequest struct {
	Name            string `json:"name" binding:"required" example:"Basic Vocabulary Quiz"`
	Type            string `json:"type" binding:"required" example:"vocabulary"`
	DifficultyLevel int    `json:"difficulty_level" binding:"required,min=1,max=5" example:"1"`
	RequiresAudio   bool   `json:"requires_audio" binding:"required" example:"false"`
	Instructions    string `json:"instructions" binding:"required" example:"Match the Italian words with their English translations"`
	ThumbnailURL    string `json:"thumbnail_url,omitempty" example:"https://example.com/thumbnail.jpg"`
	Category        string `json:"category" binding:"required" example:"basics"`
}

// Validate validates the CreateStudyActivityRequest
func (r *CreateStudyActivityRequest) Validate() error {
	// Validate type
	validType := false
	for _, t := range ValidActivityTypes {
		if r.Type == t {
			validType = true
			break
		}
	}
	if !validType {
		return fmt.Errorf("invalid type: must be one of %v", ValidActivityTypes)
	}

	// Validate difficulty level
	if r.DifficultyLevel < 1 || r.DifficultyLevel > 5 {
		return fmt.Errorf("difficulty_level must be between 1 and 5")
	}

	// Validate category
	validCategory := false
	for _, c := range ValidActivityCategories {
		if r.Category == c {
			validCategory = true
			break
		}
	}
	if !validCategory {
		return fmt.Errorf("invalid category: must be one of %v", ValidActivityCategories)
	}

	return nil
}

// StudyActivity represents a study activity in the API response
type StudyActivity struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	DifficultyLevel int       `json:"difficulty_level"`
	RequiresAudio   bool      `json:"requires_audio"`
	Instructions    string    `json:"instructions"`
	ThumbnailURL    string    `json:"thumbnail_url,omitempty"`
	Category        string    `json:"category"`
	CreatedAt       time.Time `json:"created_at"`
}
