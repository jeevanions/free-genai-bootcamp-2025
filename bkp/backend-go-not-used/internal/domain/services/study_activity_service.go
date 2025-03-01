package services

import (
	"context"
	"errors"

	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/db/repository"
)

var (
	ErrInvalidActivity = errors.New("invalid activity data")
)

type StudyActivityService interface {
	CreateActivity(ctx context.Context, activity *models.StudyActivity) error
	GetActivity(ctx context.Context, id int64) (*models.StudyActivity, error)
	ListActivities(ctx context.Context, page, pageSize int) ([]*models.StudyActivity, error)
	GetCategories(ctx context.Context) ([]string, error)
	GetRecommended(ctx context.Context, count int) ([]*models.StudyActivity, error)
}

type studyActivityServiceImpl struct {
	repo repository.StudyActivityRepository
}

func NewStudyActivityService(repo repository.StudyActivityRepository) StudyActivityService {
	return &studyActivityServiceImpl{
		repo: repo,
	}
}

func (s *studyActivityServiceImpl) CreateActivity(ctx context.Context, activity *models.StudyActivity) error {
	if err := validateActivity(activity); err != nil {
		return err
	}
	return s.repo.Create(ctx, activity)
}

func (s *studyActivityServiceImpl) GetActivity(ctx context.Context, id int64) (*models.StudyActivity, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *studyActivityServiceImpl) ListActivities(ctx context.Context, page, pageSize int) ([]*models.StudyActivity, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, offset, pageSize)
}

func (s *studyActivityServiceImpl) GetCategories(ctx context.Context) ([]string, error) {
	return s.repo.GetCategories(ctx)
}

func (s *studyActivityServiceImpl) GetRecommended(ctx context.Context, count int) ([]*models.StudyActivity, error) {
	if count < 1 {
		count = 5
	}
	return s.repo.GetRecommended(ctx, count)
}

func validateActivity(activity *models.StudyActivity) error {
	if activity.Name == "" || activity.Type == "" || activity.Instructions == "" {
		return ErrInvalidActivity
	}
	if activity.DifficultyLevel < 1 || activity.DifficultyLevel > 5 {
		return ErrInvalidActivity
	}
	return nil
}
