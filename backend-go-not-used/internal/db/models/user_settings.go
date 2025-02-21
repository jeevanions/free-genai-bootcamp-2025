package models

import "time"

// Valid values for preferred study time
var ValidStudyTimes = []string{"Morning", "Afternoon", "Evening"}

// UserSettings represents user preferences and settings in the database
type UserSettings struct {
	ID                  int64     `db:"id" json:"id"`
	UserID             int64     `db:"user_id" json:"user_id"`
	DailyWordGoal      int       `db:"daily_word_goal" json:"daily_word_goal"`
	PreferredStudyTime string    `db:"preferred_study_time" json:"preferred_study_time"`
	NotificationEnabled bool      `db:"notification_enabled" json:"notification_enabled"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the name of the table for the UserSettings model
func (UserSettings) TableName() string {
	return "user_settings"
}

// DefaultSettings returns default user settings for a new user
func DefaultSettings(userID int64) UserSettings {
	return UserSettings{
		UserID:             userID,
		DailyWordGoal:      10,
		PreferredStudyTime: "Morning",
		NotificationEnabled: true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}
