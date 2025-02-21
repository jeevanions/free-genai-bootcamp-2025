-- name: GetWordReviewsBySessionID :many
SELECT * FROM word_review_items
WHERE study_session_id = ?;
