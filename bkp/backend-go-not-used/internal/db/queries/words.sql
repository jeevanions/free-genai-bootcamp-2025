-- name: GetWord :one
SELECT * FROM words
WHERE id = ? LIMIT 1;

-- name: ListWords :many
SELECT * FROM words
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateWord :execresult
INSERT INTO words (
  italian, english, parts_of_speech, gender,
  number, difficulty_level, verb_conjugation, notes
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateWord :exec
UPDATE words
SET italian = ?,
    english = ?,
    parts_of_speech = ?,
    gender = ?,
    number = ?,
    difficulty_level = ?,
    verb_conjugation = ?,
    notes = ?
WHERE id = ?;

-- name: DeleteWord :exec
DELETE FROM words
WHERE id = ?; 