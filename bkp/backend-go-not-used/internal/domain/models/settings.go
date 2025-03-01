package models

import "time"

type UserSettings struct {
    ID                  int64     `json:"id"`
    NotificationEnabled bool      `json:"notification_enabled"`
    StudyReminderTime   string    `json:"study_reminder_time"`
    DifficultyPref      int       `json:"difficulty_preference"`
    UITheme             string    `json:"ui_theme"`
    CreatedAt           time.Time `json:"created_at"`
    UpdatedAt           time.Time `json:"updated_at"`
}

type StudyStreak struct {
    ID             int64     `json:"id"`
    CurrentStreak  int       `json:"current_streak"`
    LongestStreak  int       `json:"longest_streak"`
    LastStudyDate  time.Time `json:"last_study_date"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}
