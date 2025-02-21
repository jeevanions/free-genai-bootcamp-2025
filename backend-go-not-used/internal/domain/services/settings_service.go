package services

import (
	"context"

	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/db/repository"
)

type SettingsService interface {
	GetSettings(ctx context.Context, userID int64) (*models.Settings, error)
	UpdateSettings(ctx context.Context, settings *models.Settings) error
	GetPreferences(ctx context.Context, userID int64) (*models.Preferences, error)
	UpdatePreferences(ctx context.Context, preferences *models.Preferences) error
}

type settingsServiceImpl struct {
	repo repository.SettingsRepository
}

func NewSettingsService(repo repository.SettingsRepository) SettingsService {
	return &settingsServiceImpl{repo: repo}
}

func (s *settingsServiceImpl) GetSettings(ctx context.Context, userID int64) (*models.Settings, error) {
	return s.repo.GetSettings(ctx, userID)
}

func (s *settingsServiceImpl) UpdateSettings(ctx context.Context, settings *models.Settings) error {
	return s.repo.UpdateSettings(ctx, settings)
}

func (s *settingsServiceImpl) GetPreferences(ctx context.Context, userID int64) (*models.Preferences, error) {
	return s.repo.GetPreferences(ctx, userID)
}

func (s *settingsServiceImpl) UpdatePreferences(ctx context.Context, preferences *models.Preferences) error {
	return s.repo.UpdatePreferences(ctx, preferences)
}
