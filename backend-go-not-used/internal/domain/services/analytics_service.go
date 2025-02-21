package services

import (
	"context"

	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/db/repository"
)

type AnalyticsService interface {
	GetSessionAnalytics(ctx context.Context, userID int64) (*models.SessionAnalytics, error)
	GetSessionCalendar(ctx context.Context, userID int64) (*models.SessionCalendar, error)
}

type analyticsServiceImpl struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsService(repo repository.AnalyticsRepository) AnalyticsService {
	return &analyticsServiceImpl{repo: repo}
}

func (s *analyticsServiceImpl) GetSessionAnalytics(ctx context.Context, userID int64) (*models.SessionAnalytics, error) {
	return s.repo.GetSessionAnalytics(ctx, userID)
}

func (s *analyticsServiceImpl) GetSessionCalendar(ctx context.Context, userID int64) (*models.SessionCalendar, error) {
	return s.repo.GetSessionCalendar(ctx, userID)
}
