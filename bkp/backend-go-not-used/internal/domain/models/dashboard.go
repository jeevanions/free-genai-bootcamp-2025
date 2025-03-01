package models

type StudyProgress struct {
    TotalWordsStudied    int     `json:"total_words_studied"`
    TotalAvailableWords  int     `json:"total_available_words"`
    MasteryPercentage    float64 `json:"mastery_percentage"`
}

type QuickStats struct {
    SuccessRate    float64 `json:"success_rate"`
    TotalSessions  int     `json:"total_sessions"`
    ActiveGroups   int     `json:"active_groups"`
}

type MasteryMetrics struct {
    OverallMastery float64            `json:"overall_mastery"`
    ByCategory     map[string]float64 `json:"by_category"`
}
