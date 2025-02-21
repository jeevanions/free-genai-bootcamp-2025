package models

import "time"

type StudyProgress struct {
	WordsLearned   int     `json:"words_learned"`
	WordsToReview  int     `json:"words_to_review"`
	CompletionRate float64 `json:"completion_rate"`
}

type QuickStats struct {
	TodayStudyTime time.Duration `json:"today_study_time"`
	WeekStudyTime  time.Duration `json:"week_study_time"`
	MonthStudyTime time.Duration `json:"month_study_time"`
}

type Streak struct {
	CurrentStreak int `json:"current_streak"`
	LongestStreak int `json:"longest_streak"`
	TotalDays     int `json:"total_days"`
}

type MasteryMetrics struct {
	Beginner     int `json:"beginner"`
	Intermediate int `json:"intermediate"`
	Advanced     int `json:"advanced"`
	Mastered     int `json:"mastered"`
}
