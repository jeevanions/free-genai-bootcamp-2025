package repository

import (
	"context"
	"database/sql"

	"github.com/jeevanions/italian-learning/internal/db/models"
)

type SQLiteSettingsRepository struct {
	db *sql.DB
}

func NewSQLiteSettingsRepository(db *sql.DB) *SQLiteSettingsRepository {
	return &SQLiteSettingsRepository{db: db}
}

func (r *SQLiteSettingsRepository) GetSettings(ctx context.Context, userID int64) (*models.Settings, error) {
	// TODO: Implement
	return nil, nil
}

func (r *SQLiteSettingsRepository) UpdateSettings(ctx context.Context, settings *models.Settings) error {
	// TODO: Implement
	return nil
}

func (r *SQLiteSettingsRepository) GetPreferences(ctx context.Context, userID int64) (*models.Preferences, error) {
	// TODO: Implement
	return nil, nil
}

func (r *SQLiteSettingsRepository) UpdatePreferences(ctx context.Context, preferences *models.Preferences) error {
	// TODO: Implement
	return nil
}
