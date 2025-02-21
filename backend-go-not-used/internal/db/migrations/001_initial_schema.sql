-- +goose Up
-- Core Vocabulary with Italian-specific features
CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    italian TEXT NOT NULL,
    english TEXT NOT NULL,
    parts_of_speech TEXT NOT NULL,
    gender TEXT,
    number TEXT,
    difficulty_level INTEGER NOT NULL,
    verb_conjugation JSON,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(italian),
    CHECK (gender IS NULL OR gender IN ('masculine', 'feminine', 'neuter')),
    CHECK (number IS NULL OR number IN ('singular', 'plural')),
    CHECK (difficulty_level BETWEEN 1 AND 5)
);

CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    difficulty_level INTEGER NOT NULL CHECK (difficulty_level BETWEEN 1 AND 5),
    category TEXT NOT NULL,               -- grammar/thematic/situational
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE words_groups (
    id INTEGER PRIMARY KEY,
    word_id INTEGER REFERENCES words(id),
    group_id INTEGER REFERENCES groups(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS study_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,                   -- vocabulary/grammar/pronunciation
    requires_audio BOOLEAN NOT NULL,      -- For pronunciation exercises
    difficulty_level INTEGER NOT NULL,    -- Progressive difficulty
    instructions TEXT NOT NULL,           -- Activity guidelines
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE study_sessions (
    id INTEGER PRIMARY KEY,
    group_id INTEGER REFERENCES groups(id),
    study_activity_id INTEGER REFERENCES study_activities(id),
    total_words INTEGER NOT NULL,         -- Words attempted
    correct_words INTEGER NOT NULL,       -- Successful attempts
    duration_seconds INTEGER NOT NULL,    -- Session length
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE word_review_items (
    id INTEGER PRIMARY KEY,
    word_id INTEGER REFERENCES words(id),
    study_session_id INTEGER REFERENCES study_sessions(id),
    correct BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create all indexes
CREATE INDEX idx_words_italian ON words(italian);
CREATE INDEX idx_words_difficulty ON words(difficulty_level);
CREATE INDEX idx_words_parts_speech ON words(parts_of_speech);

CREATE INDEX idx_words_groups_word ON words_groups(word_id);
CREATE INDEX idx_words_groups_group ON words_groups(group_id);

CREATE INDEX idx_groups_category ON groups(category);
CREATE INDEX idx_groups_difficulty ON groups(difficulty_level);

CREATE INDEX idx_study_sessions_group ON study_sessions(group_id);
CREATE INDEX idx_study_sessions_created ON study_sessions(created_at);
CREATE INDEX idx_study_sessions_metrics ON study_sessions(correct_words, total_words);

CREATE INDEX idx_study_activities_type ON study_activities(type);
CREATE INDEX idx_study_activities_difficulty ON study_activities(difficulty_level);

CREATE INDEX idx_word_review_performance ON word_review_items(word_id, correct);
CREATE INDEX idx_word_review_session ON word_review_items(study_session_id);
CREATE INDEX idx_word_review_created ON word_review_items(created_at);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS word_review_items;
DROP TABLE IF EXISTS study_sessions;
DROP TABLE IF EXISTS study_activities;
DROP TABLE IF EXISTS words_groups;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS words; 