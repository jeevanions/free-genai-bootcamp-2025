package services

import (
	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

type WordServiceInterface interface {
	GetWords(limit, offset int) (*models.WordListResponse, error)
	GetWordByID(id int64) (*models.WordResponse, error)
}

type WordService struct {
	repo repository.Repository
}

func NewWordService(repo repository.Repository) *WordService {
	return &WordService{repo: repo}
}

func (s *WordService) GetWords(limit, offset int) (*models.WordListResponse, error) {
	return s.repo.GetWords(limit, offset)
}

func (s *WordService) GetWordByID(id int64) (*models.WordResponse, error) {
	return s.repo.GetWordByID(id)
}
