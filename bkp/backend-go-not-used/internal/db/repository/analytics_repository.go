package repository

import (
	"context"
	"database/sql"

	"github.com/jeevanions/italian-learning/internal/db/models"
)

type SQLiteAnalyticsRepository struct {
	db *sql.DB
}

func NewSQLiteAnalyticsRepository(db *sql.DB) *SQLiteAnalyticsRepository {
	return &SQLiteAnalyticsRepository{db: db}
}

func (r *SQLiteAnalyticsRepository) GetSessionAnalytics(ctx context.Context, userID int64) (*models.SessionAnalytics, error) {
	// TODO: Implement
	return nil, nil
}

func (r *SQLiteAnalyticsRepository) GetSessionCalendar(ctx context.Context, userID int64) (*models.SessionCalendar, error) {
	// TODO: Implement
	return nil, nil
}
