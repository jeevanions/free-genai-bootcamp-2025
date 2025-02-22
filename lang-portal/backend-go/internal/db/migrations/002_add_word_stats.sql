-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE words ADD COLUMN correct_count INTEGER DEFAULT 0;
ALTER TABLE words ADD COLUMN wrong_count INTEGER DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
CREATE TABLE words_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    italian TEXT NOT NULL,
    english TEXT NOT NULL,
    parts JSON,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO words_temp SELECT id, italian, english, parts, created_at FROM words;
DROP TABLE words;
ALTER TABLE words_temp RENAME TO words;
