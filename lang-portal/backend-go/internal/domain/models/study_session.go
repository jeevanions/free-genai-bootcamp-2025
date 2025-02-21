package models

// StudySessionStats represents statistics for a study session
type StudySessionStats struct {
	TotalWords    int     `json:"total_words"`
	CorrectWords  int     `json:"correct_words"`
	SuccessRate   float64 `json:"success_rate"`
	AverageTime   float64 `json:"average_time"`
	TotalDuration float64 `json:"total_duration"`
}
