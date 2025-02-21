-- +goose Up
-- Add constraint to enforce valid parts of speech by recreating the table
PRAGMA foreign_keys=off;

-- Create new table with all constraints
CREATE TABLE words_new (
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
    CHECK (difficulty_level BETWEEN 1 AND 5),
    CHECK (parts_of_speech IN ('noun', 'verb', 'adjective', 'adverb', 'preposition', 'conjunction', 'interjection'))
);

-- Copy data
INSERT INTO words_new SELECT * FROM words;

-- Drop old table and rename new one
DROP TABLE words;
ALTER TABLE words_new RENAME TO words;

-- Recreate indexes
CREATE INDEX idx_words_italian ON words(italian);
CREATE INDEX idx_words_difficulty ON words(difficulty_level);
CREATE INDEX idx_words_parts_speech ON words(parts_of_speech);

PRAGMA foreign_keys=on;

-- +goose Down
-- Remove the constraint by recreating the table without it
PRAGMA foreign_keys=off;

CREATE TABLE words_new (
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

INSERT INTO words_new SELECT * FROM words;
DROP TABLE words;
ALTER TABLE words_new RENAME TO words;

-- Recreate indexes
CREATE INDEX idx_words_italian ON words(italian);
CREATE INDEX idx_words_difficulty ON words(difficulty_level);
CREATE INDEX idx_words_parts_speech ON words(parts_of_speech);

PRAGMA foreign_keys=on;
