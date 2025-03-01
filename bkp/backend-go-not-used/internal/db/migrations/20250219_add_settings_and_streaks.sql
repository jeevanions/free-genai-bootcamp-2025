-- +goose Up
-- Add new columns to study_activities
ALTER TABLE study_activities ADD COLUMN thumbnail_url TEXT;
ALTER TABLE study_activities ADD COLUMN category TEXT NOT NULL DEFAULT 'general';

-- Add new columns to study_sessions
ALTER TABLE study_sessions ADD COLUMN start_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE study_sessions ADD COLUMN end_time DATETIME;

-- Create user_settings table
CREATE TABLE user_settings (
    id INTEGER PRIMARY KEY,
    notification_enabled BOOLEAN DEFAULT true,
    study_reminder_time TEXT,
    difficulty_preference INTEGER CHECK (difficulty_preference BETWEEN 1 AND 5),
    ui_theme TEXT DEFAULT 'light',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create study_streaks table
CREATE TABLE study_streaks (
    id INTEGER PRIMARY KEY,
    current_streak INTEGER DEFAULT 0,
    longest_streak INTEGER DEFAULT 0,
    last_study_date DATE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for new columns
CREATE INDEX idx_study_activities_category ON study_activities(category);
CREATE INDEX idx_study_sessions_times ON study_sessions(start_time, end_time);
CREATE INDEX idx_study_streaks_last_date ON study_streaks(last_study_date);

-- +goose Down
DROP INDEX idx_study_streaks_last_date;
DROP INDEX idx_study_sessions_times;
DROP INDEX idx_study_activities_category;

DROP TABLE study_streaks;
DROP TABLE user_settings;

ALTER TABLE study_sessions DROP COLUMN end_time;
ALTER TABLE study_sessions DROP COLUMN start_time;

ALTER TABLE study_activities DROP COLUMN category;
ALTER TABLE study_activities DROP COLUMN thumbnail_url;
