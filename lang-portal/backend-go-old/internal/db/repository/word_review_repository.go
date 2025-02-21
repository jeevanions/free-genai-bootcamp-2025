package repository

import (
	"context"
	"database/sql"

	"github.com/jeevanions/italian-learning/internal/db/models"
)

type SQLiteWordReviewRepository struct {
	db *sql.DB
}

func NewSQLiteWordReviewRepository(db *sql.DB) *SQLiteWordReviewRepository {
	return &SQLiteWordReviewRepository{db: db}
}

func (r *SQLiteWordReviewRepository) Create(ctx context.Context, review *models.WordReviewItem) error {
	query := `
		INSERT INTO word_review_items (
			word_id, study_session_id, correct
		) VALUES (?, ?, ?)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		review.WordID,
		review.StudySessionID,
		review.Correct,
	).Scan(&review.ID, &review.CreatedAt)

	return err
}

func (r *SQLiteWordReviewRepository) ListBySession(ctx context.Context, sessionID int64) ([]*models.WordReviewItem, error) {
	query := `
		SELECT wri.id, wri.word_id, wri.study_session_id, wri.correct, wri.created_at,
			   w.italian_word, w.english_word, w.category, w.difficulty
		FROM word_review_items wri
		JOIN words w ON w.id = wri.word_id
		WHERE wri.study_session_id = ?
		ORDER BY wri.created_at`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*models.WordReviewItem
	for rows.Next() {
		review := &models.WordReviewItem{
			WordDetails: models.WordReviewDetails{},
		}
		err := rows.Scan(
			&review.ID,
			&review.WordID,
			&review.StudySessionID,
			&review.Correct,
			&review.CreatedAt,
			&review.WordDetails.ItalianWord,
			&review.WordDetails.EnglishWord,
			&review.WordDetails.WordCategory,
			&review.WordDetails.Difficulty,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, rows.Err()
}

func (r *SQLiteWordReviewRepository) GetWordStats(ctx context.Context, wordID int64) (*models.WordStats, error) {
	stats := &models.WordStats{WordID: wordID}
	query := `
		SELECT 
			COUNT(*) as total_attempts,
			SUM(CASE WHEN correct THEN 1 ELSE 0 END) as correct_attempts
		FROM word_review_items
		WHERE word_id = ?`

	err := r.db.QueryRowContext(ctx, query, wordID).Scan(
		&stats.TotalAttempts,
		&stats.CorrectAttempts,
	)

	if err != nil {
		return nil, err
	}

	stats.SuccessRate = float64(stats.CorrectAttempts) / float64(stats.TotalAttempts) * 100
	return stats, nil
}
