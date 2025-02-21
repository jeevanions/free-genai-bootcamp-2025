package models

import "fmt"

// UpdateUserSettingsRequest represents the request to update user settings
type UpdateUserSettingsRequest struct {
	DailyWordGoal      *int    `json:"daily_word_goal,omitempty" binding:"omitempty,min=1,max=100" example:"15"`
	PreferredStudyTime *string `json:"preferred_study_time,omitempty" example:"Afternoon"`
	NotificationEnabled *bool   `json:"notification_enabled,omitempty" example:"true"`
}

// Validate validates the UpdateUserSettingsRequest
func (r *UpdateUserSettingsRequest) Validate() error {
	if r.PreferredStudyTime != nil {
		validTime := false
		for _, t := range []string{"Morning", "Afternoon", "Evening"} {
			if *r.PreferredStudyTime == t {
				validTime = true
				break
			}
		}
		if !validTime {
			return fmt.Errorf("invalid preferred_study_time: must be one of Morning, Afternoon, Evening")
		}
	}
	return nil
}

// UserSettingsResponse represents user settings in API responses
type UserSettingsResponse struct {
	ID                  int64  `json:"id"`
	UserID             int64  `json:"user_id"`
	DailyWordGoal      int    `json:"daily_word_goal"`
	PreferredStudyTime string `json:"preferred_study_time"`
	NotificationEnabled bool   `json:"notification_enabled"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

// UserSettingsStats represents statistics about user's study settings and progress
type UserSettingsStats struct {
	DailyWordGoal     int     `json:"daily_word_goal"`
	WordsLearnedToday int     `json:"words_learned_today"`
	GoalProgress      float64 `json:"goal_progress"`
	StudyStreak       int     `json:"study_streak"`
	TotalStudyDays    int     `json:"total_study_days"`
}
