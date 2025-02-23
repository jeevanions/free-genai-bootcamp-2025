package services

import (
	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

type WordServiceInterface interface {
	GetWords(limit, offset int) (*models.WordListResponse, error)
	GetWordByID(id int64) (*models.WordResponse, error)
	ImportWords(groupID int64, words []string) (*models.ImportWordsResponse, error)
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

func (s *WordService) ImportWords(groupID int64, words []string) (*models.ImportWordsResponse, error) {
	importedCount := 0
	for _, word := range words {
		// Create a basic word response
		wordResp := &models.WordResponse{
			Italian: word,
			English: "", // You might want to add translation logic here
		}
		_, err := s.repo.CreateWord(wordResp)
		if err != nil {
			continue
		}
		importedCount++
	}
	return &models.ImportWordsResponse{ImportedCount: importedCount}, nil
}
