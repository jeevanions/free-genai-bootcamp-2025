-- name: GetStudyActivityByID :one
SELECT * FROM study_activities
WHERE id = ? LIMIT 1;

-- name: GetStudySessionsByActivityID :many
SELECT 
    s.id,
    s.created_at,
    s.group_id,
    sa.name as activity_name
FROM study_sessions s
JOIN study_activities sa ON s.study_activity_id = sa.id
WHERE sa.id = ?
ORDER BY s.created_at DESC
LIMIT 100;

-- name: GetGroupByID :one
SELECT * FROM groups
WHERE id = ? LIMIT 1;

-- name: GetSessionStats :one
SELECT 
    COUNT(*) as total_words,
    SUM(CASE WHEN correct THEN 1 ELSE 0 END) as correct_words
FROM word_review_items
WHERE study_session_id = ?;
