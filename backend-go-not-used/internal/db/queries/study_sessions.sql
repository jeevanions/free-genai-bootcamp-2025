-- name: GetStudySession :one
SELECT * FROM study_sessions
WHERE id = ? LIMIT 1;

-- name: ListGroupSessions :many
SELECT * FROM study_sessions
WHERE group_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CreateStudySession :execresult
INSERT INTO study_sessions (
  group_id, study_activity_id, total_words,
  correct_words, duration_seconds
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: GetGroupStats :one
SELECT 
  COUNT(*) as total_sessions,
  SUM(total_words) as total_words,
  SUM(correct_words) as total_correct,
  SUM(duration_seconds) as total_duration
FROM study_sessions
WHERE group_id = ?; 