-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Update words_count for existing groups
UPDATE groups
SET words_count = (
    SELECT COUNT(*)
    FROM words_groups
    WHERE words_groups.group_id = groups.id
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
-- No rollback needed since we're just updating data
