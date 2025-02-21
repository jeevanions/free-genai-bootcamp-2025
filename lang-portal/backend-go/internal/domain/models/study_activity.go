package models

import "time"

type StudyActivityResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

type StudyActivityListResponse struct {
	Items      []StudyActivityResponse `json:"items"`
	Pagination PaginationResponse     `json:"pagination"`
}

type LaunchStudyActivityRequest struct {
	GroupID int64 `json:"group_id" binding:"required"`
}

type LaunchStudyActivityResponse struct {
	StudySessionID   int64     `json:"study_session_id"`
	StudyActivityID  int64     `json:"study_activity_id"`
	GroupID         int64     `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type StudySessionResponse struct {
	ID           int64     `json:"id"`
	ActivityName string    `json:"activity_name"`
	GroupID      int64     `json:"group_id"`
	GroupName    string    `json:"group_name"`
	CreatedAt    time.Time `json:"created_at"`
	WordsCount   int       `json:"words_count"`
	CorrectCount int       `json:"correct_count"`
}

type StudySessionsListResponse struct {
	Items []StudySessionResponse `json:"items"`
}
