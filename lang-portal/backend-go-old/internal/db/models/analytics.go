package models

import "time"

type SessionAnalytics struct {
	TotalSessions     int           `json:"total_sessions"`
	TotalStudyTime    time.Duration `json:"total_study_time"`
	AverageSessionLen time.Duration `json:"average_session_length"`
	WordsLearned      int           `json:"words_learned"`
	WordsReviewed     int           `json:"words_reviewed"`
}

type SessionCalendar struct {
	Sessions []CalendarSession `json:"sessions"`
}

type CalendarSession struct {
	Date      time.Time     `json:"date"`
	Duration  time.Duration `json:"duration"`
	WordCount int           `json:"word_count"`
}
