package services

import (
	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

type DashboardService struct {
	repo repository.Repository
}

func NewDashboardService(repo repository.Repository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetLastStudySession() (*models.DashboardLastStudySession, error) {
	return s.repo.GetLastStudySession()
}

func (s *DashboardService) GetStudyProgress() (*models.DashboardStudyProgress, error) {
	return s.repo.GetStudyProgress()
}

func (s *DashboardService) GetQuickStats() (*models.DashboardQuickStats, error) {
	return s.repo.GetQuickStats()
}
