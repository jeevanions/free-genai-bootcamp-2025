package models

import "time"

// StudySessionStats represents statistics for a study session
type StudySessionStats struct {
	TotalWords    int     `json:"total_words"`
	CorrectWords  int     `json:"correct_words"`
	SuccessRate   float64 `json:"success_rate"`
	AverageTime   float64 `json:"average_time"`
	TotalDuration float64 `json:"total_duration"`
}

// StudySessionListResponse represents the response for listing study sessions
type StudySessionListResponse struct {
	Items []StudySessionDetailResponse `json:"items"`
	Pagination PaginationResponse     `json:"pagination"`
}

// StudySessionDetailResponse represents a detailed study session response
type StudySessionDetailResponse struct {
	ID            int64            `json:"id"`
	ActivityName  string           `json:"activity_name"`
	GroupName     string           `json:"group_name"`
	CreatedAt     time.Time        `json:"created_at"`
	Stats         StudySessionStats `json:"stats"`
	ReviewItems   []WordReviewItem  `json:"review_items"`
}
