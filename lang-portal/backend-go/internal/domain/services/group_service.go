package services

import (
	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

type GroupServiceInterface interface {
	GetGroups(limit, offset int) (*models.GroupListResponse, error)
	GetGroupByID(id int64) (*models.GroupDetailResponse, error)
	GetGroupWords(groupID int64, limit, offset int) (*models.GroupWordsResponse, error)
	GetGroupStudySessions(groupID int64, limit, offset int) (*models.GroupStudySessionsResponse, error)
}

type GroupService struct {
	repo repository.Repository
}

func NewGroupService(repo repository.Repository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) GetGroups(limit, offset int) (*models.GroupListResponse, error) {
	return s.repo.GetGroups(limit, offset)
}

func (s *GroupService) GetGroupByID(id int64) (*models.GroupDetailResponse, error) {
	return s.repo.GetGroupByID(id)
}

func (s *GroupService) GetGroupWords(groupID int64, limit, offset int) (*models.GroupWordsResponse, error) {
	return s.repo.GetGroupWords(groupID, limit, offset)
}

func (s *GroupService) GetGroupStudySessions(groupID int64, limit, offset int) (*models.GroupStudySessionsResponse, error) {
	return s.repo.GetGroupStudySessions(groupID, limit, offset)
}
