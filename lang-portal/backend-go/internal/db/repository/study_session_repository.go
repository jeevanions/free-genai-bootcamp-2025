package repository

import (
	"context"
	"database/sql"

	"github.com/jeevanions/italian-learning/internal/db/models"
)

type SQLiteStudySessionRepository struct {
	db *sql.DB
}

func NewSQLiteStudySessionRepository(db *sql.DB) *SQLiteStudySessionRepository {
	return &SQLiteStudySessionRepository{db: db}
}

func (r *SQLiteStudySessionRepository) Create(ctx context.Context, session *models.StudySession) error {
	query := `
		INSERT INTO study_sessions (
			group_id, study_activity_id, total_words, correct_words, duration_seconds
		) VALUES (?, ?, ?, ?, ?)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		session.GroupID,
		session.StudyActivityID,
		session.TotalWords,
		session.CorrectWords,
		session.DurationSeconds,
	).Scan(&session.ID, &session.CreatedAt)

	return err
}

func (r *SQLiteStudySessionRepository) GetByID(ctx context.Context, id int64) (*models.StudySession, error) {
	session := &models.StudySession{}
	query := `
		SELECT id, group_id, study_activity_id, total_words, 
			   correct_words, duration_seconds, created_at
		FROM study_sessions
		WHERE id = ?`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.GroupID,
		&session.StudyActivityID,
		&session.TotalWords,
		&session.CorrectWords,
		&session.DurationSeconds,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SQLiteStudySessionRepository) ListByGroupID(ctx context.Context, groupID int64, offset, limit int) ([]*models.StudySession, error) {
	query := `
		SELECT id, group_id, study_activity_id, total_words, 
			   correct_words, duration_seconds, created_at
		FROM study_sessions
		WHERE group_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*models.StudySession
	for rows.Next() {
		session := &models.StudySession{}
		err := rows.Scan(
			&session.ID,
			&session.GroupID,
			&session.StudyActivityID,
			&session.TotalWords,
			&session.CorrectWords,
			&session.DurationSeconds,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

func (r *SQLiteStudySessionRepository) GetStats(ctx context.Context, groupID int64) (*models.StudyStats, error) {
	stats := &models.StudyStats{}
	query := `
		SELECT 
			COUNT(*) as total_sessions,
			SUM(total_words) as total_words,
			SUM(correct_words) as total_correct,
			SUM(duration_seconds) as total_duration
		FROM study_sessions
		WHERE group_id = ?`

	err := r.db.QueryRowContext(ctx, query, groupID).Scan(
		&stats.TotalSessions,
		&stats.TotalWords,
		&stats.TotalCorrect,
		&stats.TotalDuration,
	)

	if err != nil {
		return nil, err
	}

	stats.CalculateAccuracy()
	return stats, nil
}
