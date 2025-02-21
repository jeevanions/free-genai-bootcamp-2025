package services

import (
	"context"
	"fmt"

	apimodels "github.com/jeevanions/italian-learning/internal/api/models"
	dbmodels "github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/db/repository"
)

type DashboardService interface {
	GetLastStudySession(ctx context.Context, userID int64) (*apimodels.LastStudySessionResponse, error)
	GetStudyProgress(ctx context.Context, userID int64) (*dbmodels.StudyProgress, error)
	GetQuickStats(ctx context.Context, userID int64) (*dbmodels.QuickStats, error)
	GetStreak(ctx context.Context, userID int64) (*dbmodels.Streak, error)
	GetMasteryMetrics(ctx context.Context, userID int64) (*dbmodels.MasteryMetrics, error)
}

type dashboardServiceImpl struct {
	dashboardRepo repository.DashboardRepository
	sessionRepo  repository.StudySessionRepository
	reviewRepo   repository.WordReviewRepository
	groupRepo    repository.GroupRepository
	activityRepo repository.StudyActivityRepository
}

func NewDashboardService(
	dashboardRepo repository.DashboardRepository,
	sessionRepo repository.StudySessionRepository,
	reviewRepo repository.WordReviewRepository,
	groupRepo repository.GroupRepository,
	activityRepo repository.StudyActivityRepository,
) DashboardService {
	return &dashboardServiceImpl{
		dashboardRepo: dashboardRepo,
		sessionRepo:  sessionRepo,
		reviewRepo:   reviewRepo,
		groupRepo:    groupRepo,
		activityRepo: activityRepo,
	}
}

func (s *dashboardServiceImpl) GetLastStudySession(ctx context.Context, userID int64) (*apimodels.LastStudySessionResponse, error) {
	// Get the last study session
	session, err := s.dashboardRepo.GetLastStudySession(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get the group name
	group, err := s.groupRepo.GetByID(ctx, session.GroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %v", err)
	}

	// Get the activity name
	activity, err := s.activityRepo.GetByID(ctx, session.StudyActivityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity: %v", err)
	}

	// Build the response
	return &apimodels.LastStudySessionResponse{
		ID:              session.ID,
		GroupID:         session.GroupID,
		GroupName:       group.Name,
		ActivityID:      session.StudyActivityID,
		ActivityName:    activity.Name,
		TotalWords:      session.TotalWords,
		CorrectWords:    session.CorrectWords,
		DurationSeconds: session.DurationSeconds,
		StartTime:       session.StartTime,
		EndTime:         session.EndTime,
	}, nil
}

func (s *dashboardServiceImpl) GetStudyProgress(ctx context.Context, userID int64) (*dbmodels.StudyProgress, error) {
	return s.dashboardRepo.GetStudyProgress(ctx, userID)
}

func (s *dashboardServiceImpl) GetQuickStats(ctx context.Context, userID int64) (*dbmodels.QuickStats, error) {
	return s.dashboardRepo.GetQuickStats(ctx, userID)
}

func (s *dashboardServiceImpl) GetStreak(ctx context.Context, userID int64) (*dbmodels.Streak, error) {
	return s.dashboardRepo.GetStreak(ctx, userID)
}

func (s *dashboardServiceImpl) GetMasteryMetrics(ctx context.Context, userID int64) (*dbmodels.MasteryMetrics, error) {
	return s.dashboardRepo.GetMasteryMetrics(ctx, userID)
}
