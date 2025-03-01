package models

import "time"

type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	StudyActivityID int64     `json:"study_activity_id"`
	TotalWords      int       `json:"total_words"`
	CorrectWords    int       `json:"correct_words"`
	DurationSeconds int       `json:"duration_seconds"`
	StartTime       time.Time `json:"start_time" db:"start_time"`
	EndTime         time.Time `json:"end_time" db:"end_time"`
	CreatedAt       time.Time `json:"created_at"`
}

type SessionAnalytics struct {
	PerformanceTrend []DailyPerformance `json:"performance_trend"`
	TotalStudyTime   int                `json:"total_study_time"`
	AvgSessionLength int                `json:"average_session_length"`
}

type DailyPerformance struct {
	Date         string  `json:"date"`
	SuccessRate  float64 `json:"success_rate"`
	WordsStudied int     `json:"words_studied"`
}

type CalendarSession struct {
	Date         string  `json:"date"`
	SessionCount int     `json:"session_count"`
	TotalWords   int     `json:"total_words"`
	SuccessRate  float64 `json:"success_rate"`
}
