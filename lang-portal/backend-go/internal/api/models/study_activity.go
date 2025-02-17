package models

import "time"

type CreateStudyActivityRequest struct {
	Name            string `json:"name" binding:"required"`
	Type            string `json:"type" binding:"required"`
	DifficultyLevel int    `json:"difficulty_level" binding:"required,min=1,max=5"`
	RequiresAudio   bool   `json:"requires_audio" binding:"required"`
	Instructions    string `json:"instructions" binding:"required"`
}

type StudyActivity struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	DifficultyLevel int       `json:"difficulty_level"`
	RequiresAudio   bool      `json:"requires_audio"`
	Instructions    string    `json:"instructions"`
	CreatedAt       time.Time `json:"created_at"`
}
