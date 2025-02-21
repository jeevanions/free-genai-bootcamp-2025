package services

import (
	"context"
	"errors"
	"time"

	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/db/repository"
)

var (
	ErrInvalidSession = errors.New("invalid session data")
)

type StudySessionService interface {
	CreateSession(ctx context.Context, session *models.StudySession) error
	GetSessionStats(ctx context.Context, groupID int64) (*models.StudyStats, error)
	ListGroupSessions(ctx context.Context, groupID int64, page, pageSize int) ([]*models.StudySession, error)
	ListSessions(page, limit int) ([]*models.StudySession, int, error)
}

type StudySessionServiceImpl struct {
	sessionRepo  repository.StudySessionRepository
	groupRepo    repository.GroupRepository
	activityRepo repository.StudyActivityRepository
}

func NewStudySessionService(
	sessionRepo repository.StudySessionRepository,
	groupRepo repository.GroupRepository,
	activityRepo repository.StudyActivityRepository,
) StudySessionService {
	return &StudySessionServiceImpl{
		sessionRepo:  sessionRepo,
		groupRepo:    groupRepo,
		activityRepo: activityRepo,
	}
}

func (s *StudySessionServiceImpl) CreateSession(ctx context.Context, session *models.StudySession) error {
	if err := validateSession(session); err != nil {
		return err
	}

	// Verify group exists
	if _, err := s.groupRepo.GetByID(ctx, session.GroupID); err != nil {
		return err
	}

	// Verify activity exists
	if _, err := s.activityRepo.GetByID(ctx, session.StudyActivityID); err != nil {
		return err
	}

	// Set start and end time if not already set
	if session.StartTime.IsZero() {
		session.StartTime = time.Now()
	}
	if session.EndTime.IsZero() {
		session.EndTime = session.StartTime.Add(time.Duration(session.DurationSeconds) * time.Second)
	}

	return s.sessionRepo.Create(ctx, session)
}

func (s *StudySessionServiceImpl) GetSessionStats(ctx context.Context, groupID int64) (*models.StudyStats, error) {
	// Verify group exists
	if _, err := s.groupRepo.GetByID(ctx, groupID); err != nil {
		return nil, err
	}

	return s.sessionRepo.GetStats(ctx, groupID)
}

func (s *StudySessionServiceImpl) ListGroupSessions(ctx context.Context, groupID int64, page, pageSize int) ([]*models.StudySession, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	return s.sessionRepo.ListByGroupID(ctx, groupID, offset, pageSize)
}

func (s *StudySessionServiceImpl) GetByID(ctx context.Context, id int64) (*models.StudySession, error) {
	return s.sessionRepo.GetByID(ctx, id)
}

func validateSession(session *models.StudySession) error {
	if session.GroupID <= 0 || session.StudyActivityID <= 0 {
		return ErrInvalidSession
	}
	if session.TotalWords < 0 || session.CorrectWords < 0 {
		return ErrInvalidSession
	}
	if session.CorrectWords > session.TotalWords {
		return ErrInvalidSession
	}
	if session.DurationSeconds <= 0 {
		return ErrInvalidSession
	}
	return nil
}

func (s *StudySessionServiceImpl) ListSessions(page, limit int) ([]*models.StudySession, int, error) {
	ctx := context.Background()
	return s.sessionRepo.List(ctx, page, limit)
}
