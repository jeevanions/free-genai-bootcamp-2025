package services

import "github.com/jeevanions/lang-portal/backend-go/internal/domain/models"

// DashboardServiceInterface defines the interface for dashboard service
type DashboardServiceInterface interface {
	GetLastStudySession() (*models.DashboardLastStudySession, error)
	GetStudyProgress() (*models.DashboardStudyProgress, error)
	GetQuickStats() (*models.DashboardQuickStats, error)
}
