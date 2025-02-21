package models

import "time"

// DashboardLastStudySession represents the last study session information
type DashboardLastStudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int64     `json:"study_activity_id"`
	GroupName       string    `json:"group_name"`
}

// DashboardStudyProgress represents study progress statistics
type DashboardStudyProgress struct {
	TotalWordsStudied    int `json:"total_words_studied"`
	TotalAvailableWords int `json:"total_available_words"`
}

// DashboardQuickStats represents quick overview statistics
type DashboardQuickStats struct {
	SuccessRate        float64 `json:"success_rate"`
	TotalStudySessions int     `json:"total_study_sessions"`
	TotalActiveGroups  int     `json:"total_active_groups"`
	StudyStreakDays    int     `json:"study_streak_days"`
}
