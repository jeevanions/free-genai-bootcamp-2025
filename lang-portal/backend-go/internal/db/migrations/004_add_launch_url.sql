-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE study_activities ADD COLUMN launch_url TEXT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
CREATE TABLE study_activities_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    thumbnail_url TEXT,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO study_activities_temp 
SELECT id, name, thumbnail_url, description, created_at
FROM study_activities;

DROP TABLE study_activities;
ALTER TABLE study_activities_temp RENAME TO study_activities;
