package services

import (
	"context"
	"errors"

	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/db/repository"
)

var (
	ErrInvalidGroup = errors.New("invalid group data")
)

// GroupService defines the interface for group-related operations
type GroupService interface {
	CreateGroup(ctx context.Context, group *models.Group) error
	GetGroup(ctx context.Context, id int64) (*models.Group, error)
	ListGroups(ctx context.Context, page, pageSize int) ([]*models.Group, error)
}

type groupServiceImpl struct {
	repo repository.GroupRepository
}

func NewGroupService(repo repository.GroupRepository) GroupService {
	return &groupServiceImpl{
		repo: repo,
	}
}

func (s *groupServiceImpl) CreateGroup(ctx context.Context, group *models.Group) error {
	if err := validateGroup(group); err != nil {
		return err
	}
	return s.repo.Create(ctx, group)
}

func (s *groupServiceImpl) GetGroup(ctx context.Context, id int64) (*models.Group, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *groupServiceImpl) ListGroups(ctx context.Context, page, pageSize int) ([]*models.Group, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, offset, pageSize)
}

func validateGroup(group *models.Group) error {
	if group.Name == "" || group.Category == "" {
		return ErrInvalidGroup
	}
	if group.DifficultyLevel < 1 || group.DifficultyLevel > 5 {
		return ErrInvalidGroup
	}
	return nil
}
