package models

import (
	"fmt"
	"time"
)

// Valid values for group categories
var ValidGroupCategories = []string{"grammar", "thematic", "situational"}

// CreateGroupRequest represents the request to create a new group
type CreateGroupRequest struct {
	Name            string `json:"name" binding:"required" example:"Food Vocabulary"`
	Description     string `json:"description" binding:"required" example:"Common Italian words related to food and dining"`
	DifficultyLevel int    `json:"difficulty_level" binding:"required,min=1,max=5" example:"1"`
	Category        string `json:"category" binding:"required" example:"thematic"`
}

// Validate validates the CreateGroupRequest
func (r *CreateGroupRequest) Validate() error {
	// Validate difficulty level
	if r.DifficultyLevel < 1 || r.DifficultyLevel > 5 {
		return fmt.Errorf("difficulty_level must be between 1 and 5")
	}

	// Validate category
	validCategory := false
	for _, c := range ValidGroupCategories {
		if r.Category == c {
			validCategory = true
			break
		}
	}
	if !validCategory {
		return fmt.Errorf("invalid category: must be one of %v", ValidGroupCategories)
	}

	return nil
}

// GroupResponse represents a group in the API response
type GroupResponse struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	DifficultyLevel int       `json:"difficulty_level"`
	Category        string    `json:"category"`
	CreatedAt       time.Time `json:"created_at"`
	WordCount       int       `json:"word_count,omitempty"`
}

// GroupStatistics represents statistics for a group
type GroupStatistics struct {
	TotalWords      int     `json:"total_words"`
	StudiedWords    int     `json:"studied_words"`
	SuccessRate     float64 `json:"success_rate"`
	TotalSessions   int     `json:"total_sessions"`
	AverageDuration int     `json:"average_duration_seconds"`
}

// GroupProgress represents learning progress for a group
type GroupProgress struct {
	MasteryPercentage float64 `json:"mastery_percentage"`
	LastStudyDate     string  `json:"last_study_date,omitempty"`
	Streak            int     `json:"study_streak"`
}

// GroupDetailResponse represents the detailed response for a group
type GroupDetailResponse struct {
	GroupResponse
	Statistics GroupStatistics `json:"statistics"`
	Progress   GroupProgress   `json:"progress"`
}

// RecentActivityItem represents a recent study activity for a group
type RecentActivityItem struct {
	SessionID       int64     `json:"session_id"`
	ActivityType    string    `json:"activity_type"`
	WordsStudied    int       `json:"words_studied"`
	CorrectWords    int       `json:"correct_words"`
	CompletedAt     time.Time `json:"completed_at"`
}

// GroupProgressResponse represents the response for group progress endpoint
type GroupProgressResponse struct {
	TotalWords        int                  `json:"total_words"`
	StudiedWords      int                  `json:"studied_words"`
	MasteryPercentage float64              `json:"mastery_percentage"`
	LastStudyDate     string               `json:"last_study_date,omitempty"`
	StudyStreak       int                  `json:"study_streak"`
	SuccessRate       float64              `json:"success_rate"`
	TotalSessions     int                  `json:"total_sessions"`
	AverageDuration   int                  `json:"average_duration_seconds"`
	RecentActivity    []RecentActivityItem `json:"recent_activity,omitempty"`
}

// AddWordsToGroupRequest represents the request to add words to a group
type AddWordsToGroupRequest struct {
	WordIDs []int64 `json:"word_ids" binding:"required,min=1"`
}

// AddWordsToGroupResponse represents the response after adding words to a group
type AddWordsToGroupResponse struct {
	GroupID      int64   `json:"group_id"`
	WordsAdded   int     `json:"words_added"`
	TotalWords   int     `json:"total_words"`
}
