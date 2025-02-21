-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    italian TEXT NOT NULL,
    english TEXT NOT NULL,
    parts JSON,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    words_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE words_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    word_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

CREATE TABLE study_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    thumbnail_url TEXT,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE study_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    study_activity_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (study_activity_id) REFERENCES study_activities(id) ON DELETE CASCADE
);

CREATE TABLE word_review_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    word_id INTEGER NOT NULL,
    study_session_id INTEGER NOT NULL,
    correct BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_words_groups_word_id ON words_groups(word_id);
CREATE INDEX idx_words_groups_group_id ON words_groups(group_id);
CREATE INDEX idx_study_sessions_group_id ON study_sessions(group_id);
CREATE INDEX idx_word_review_items_word_id ON word_review_items(word_id);
CREATE INDEX idx_word_review_items_study_session_id ON word_review_items(study_session_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS word_review_items;
DROP TABLE IF EXISTS study_sessions;
DROP TABLE IF EXISTS study_activities;
DROP TABLE IF EXISTS words_groups;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS words;
