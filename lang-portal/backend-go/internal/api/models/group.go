package models

import "time"

type CreateGroupRequest struct {
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description" binding:"required"`
	DifficultyLevel int    `json:"difficulty_level" binding:"required,min=1,max=5"`
	Category        string `json:"category" binding:"required"`
}

type Group struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	DifficultyLevel int       `json:"difficulty_level"`
	Category        string    `json:"category"`
	CreatedAt       time.Time `json:"created_at"`
}
