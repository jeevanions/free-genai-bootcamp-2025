-- name: GetWord :one
SELECT * FROM words
WHERE id = ? LIMIT 1;

-- name: ListWords :many
SELECT * FROM words
ORDER BY id;

-- name: CreateWord :one
INSERT INTO words (
  italian, english, parts
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: UpdateWord :one
UPDATE words
SET italian = ?, english = ?, parts = ?
WHERE id = ?
RETURNING *;

-- name: DeleteWord :exec
DELETE FROM words
WHERE id = ?;
