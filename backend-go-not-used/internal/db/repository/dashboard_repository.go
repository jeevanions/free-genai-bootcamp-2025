package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jeevanions/italian-learning/internal/db/models"
)

type SQLiteDashboardRepository struct {
	db *sql.DB
}

func NewSQLiteDashboardRepository(db *sql.DB) *SQLiteDashboardRepository {
	return &SQLiteDashboardRepository{db: db}
}

func (r *SQLiteDashboardRepository) GetLastStudySession(ctx context.Context, userID int64) (*models.StudySession, error) {
	query := `
		SELECT id, group_id, study_activity_id, total_words, correct_words, 
		       duration_seconds, start_time, end_time, created_at
		FROM study_sessions
		ORDER BY created_at DESC
		LIMIT 1`

	row := r.db.QueryRowContext(ctx, query)

	var session models.StudySession
	err := row.Scan(
		&session.ID,
		&session.GroupID,
		&session.StudyActivityID,
		&session.TotalWords,
		&session.CorrectWords,
		&session.DurationSeconds,
		&session.StartTime,
		&session.EndTime,
		&session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("error getting last study session: %v", err)
	}

	return &session, nil
}

func (r *SQLiteDashboardRepository) GetStudyProgress(ctx context.Context, userID int64) (*models.StudyProgress, error) {
	// TODO: Implement
	return nil, nil
}

func (r *SQLiteDashboardRepository) GetQuickStats(ctx context.Context, userID int64) (*models.QuickStats, error) {
	// TODO: Implement
	return nil, nil
}

func (r *SQLiteDashboardRepository) GetStreak(ctx context.Context, userID int64) (*models.Streak, error) {
	// TODO: Implement
	return nil, nil
}

func (r *SQLiteDashboardRepository) GetMasteryMetrics(ctx context.Context, userID int64) (*models.MasteryMetrics, error) {
	// TODO: Implement
	return nil, nil
}
