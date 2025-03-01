package models

type StudyStats struct {
	TotalSessions int64   `json:"total_sessions"`
	TotalWords    int64   `json:"total_words"`
	TotalCorrect  int64   `json:"total_correct"`
	TotalDuration int64   `json:"total_duration"`
	Accuracy      float64 `json:"accuracy"`
}

func (s *StudyStats) CalculateAccuracy() {
	if s.TotalWords > 0 {
		s.Accuracy = float64(s.TotalCorrect) / float64(s.TotalWords) * 100
	}
}

func (s *StudyStats) AverageSessionDuration() float64 {
	if s.TotalSessions == 0 {
		return 0
	}
	return float64(s.TotalDuration) / float64(s.TotalSessions)
}
