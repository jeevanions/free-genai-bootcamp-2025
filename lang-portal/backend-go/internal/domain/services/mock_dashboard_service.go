package services

import (
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetLastStudySession() (*models.DashboardLastStudySession, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardLastStudySession), args.Error(1)
}

func (m *MockDashboardService) GetStudyProgress() (*models.DashboardStudyProgress, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardStudyProgress), args.Error(1)
}

func (m *MockDashboardService) GetQuickStats() (*models.DashboardQuickStats, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardQuickStats), args.Error(1)
}
