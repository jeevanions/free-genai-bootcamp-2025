-- name: GetGroup :one
SELECT * FROM groups
WHERE id = ? LIMIT 1;

-- name: ListGroups :many
SELECT * FROM groups
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateGroup :execresult
INSERT INTO groups (
  name, description, difficulty_level, category
) VALUES (
  ?, ?, ?, ?
);

-- name: UpdateGroup :exec
UPDATE groups
SET name = ?,
    description = ?,
    difficulty_level = ?,
    category = ?
WHERE id = ?;

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = ?; 