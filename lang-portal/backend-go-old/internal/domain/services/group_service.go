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
// Error definitions
var (
	ErrGroupNotFound = fmt.Errorf("group not found")
	ErrInvalidGroup  = fmt.Errorf("invalid group data")
)

type GroupService interface {
	CreateGroup(ctx context.Context, group *models.Group) error
	GetGroup(ctx context.Context, id int64) (*models.Group, error)
	GetGroupDetails(ctx context.Context, id int64) (*apimodels.GroupDetailResponse, error)
	GetGroupProgress(ctx context.Context, id int64) (*apimodels.GroupProgressResponse, error)
	ListGroups(ctx context.Context, page, pageSize int) ([]*models.Group, error)
	AddWordsToGroup(ctx context.Context, groupID int64, wordIDs []int64) error
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
	group, err := s.repo.GetByID(ctx, id)
	if err == repository.ErrNotFound {
		return nil, ErrGroupNotFound
	}
	return group, err
}

func (s *groupServiceImpl) GetGroupProgress(ctx context.Context, id int64) (*apimodels.GroupProgressResponse, error) {
	// Verify group exists
	_, err := s.GetGroup(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get word count
	wordCount, err := s.repo.GetWordCount(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get study statistics
	stats, err := s.repo.GetGroupStatistics(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get progress metrics
	progress, err := s.repo.GetGroupProgress(ctx, id)
	if err != nil {
		return nil, err
	}

	// Construct response
	response := &apimodels.GroupProgressResponse{
		TotalWords:         wordCount,
		StudiedWords:       stats.StudiedWords,
		MasteryPercentage:  progress.MasteryPercentage,
		LastStudyDate:      progress.LastStudyDate,
		StudyStreak:        progress.Streak,
		SuccessRate:        stats.SuccessRate,
		TotalSessions:      stats.TotalSessions,
		AverageDuration:    stats.AverageDuration,
		RecentActivity:     nil, // TODO: Add recent activity if needed
	}

	return response, nil
}

func (s *groupServiceImpl) AddWordsToGroup(ctx context.Context, groupID int64, wordIDs []int64) (*apimodels.AddWordsToGroupResponse, error) {
	// Verify group exists
	_, err := s.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Add words to group
	wordsAdded, err := s.repo.AddWordsToGroup(ctx, groupID, wordIDs)
	if err != nil {
		return nil, err
	}

	// Get total word count after adding
	totalWords, err := s.repo.GetWordCount(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return &apimodels.AddWordsToGroupResponse{
		GroupID:    groupID,
		WordsAdded: wordsAdded,
		TotalWords: totalWords,
	}, nil
}

func (s *groupServiceImpl) GetGroupDetails(ctx context.Context, id int64) (*apimodels.GroupDetailResponse, error) {
	// Get base group info
	group, err := s.GetGroup(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get word count and statistics
	wordCount, err := s.repo.GetWordCount(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get study statistics
	stats, err := s.repo.GetGroupStatistics(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get progress metrics
	progress, err := s.repo.GetGroupProgress(ctx, id)
	if err != nil {
		return nil, err
	}

	// Construct response
	response := &apimodels.GroupDetailResponse{
		GroupResponse: apimodels.GroupResponse{
			ID:              group.ID,
			Name:            group.Name,
			Description:     group.Description,
			DifficultyLevel: group.DifficultyLevel,
			Category:        group.Category,
			CreatedAt:       group.CreatedAt,
			WordCount:       wordCount,
		},
		Statistics: stats,
		Progress:   progress,
	}

	return response, nil
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
