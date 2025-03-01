package models

import "time"

type StudyActivity struct {
    ID              int64     `json:"id"`
    Name            string    `json:"name"`
    Type            string    `json:"type"`
    RequiresAudio   bool      `json:"requires_audio"`
    DifficultyLevel int       `json:"difficulty_level"`
    Instructions    string    `json:"instructions"`
    ThumbnailURL    string    `json:"thumbnail_url"`
    Category        string    `json:"category"`
    CreatedAt       time.Time `json:"created_at"`
}

type ActivityCategory struct {
    Name  string `json:"name"`
    Count int    `json:"count"`
}

type ActivityRecommendation struct {
    ActivityID int    `json:"activity_id"`
    Name       string `json:"name"`
    Reason     string `json:"reason"`
}
