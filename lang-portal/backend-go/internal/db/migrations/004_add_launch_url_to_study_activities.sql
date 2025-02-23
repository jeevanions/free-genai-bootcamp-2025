-- +goose Up
ALTER TABLE study_activities ADD COLUMN launch_url TEXT;

-- +goose Down
ALTER TABLE study_activities DROP COLUMN launch_url;
