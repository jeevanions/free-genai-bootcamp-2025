package services

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

type WordServiceInterface interface {
	GetWords(limit, offset int) (*models.WordListResponse, error)
	GetWordByID(id int64) (*models.WordResponse, error)
	ImportWords(groupID int64, words []models.WordResponse) (*models.ImportWordsResponse, error)
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

func (s *WordService) ImportWords(groupID int64, words []models.WordResponse) (*models.ImportWordsResponse, error) {
	// Verify group exists
	_, err := s.repo.GetGroupByID(groupID)
	if err != nil {
		return nil, fmt.Errorf("group not found: %v", err)
	}

	importedCount := 0
	for _, word := range words {
		// Create the word
		wordID, err := s.repo.CreateWord(&word)
		if err != nil {
			continue
		}

		// Associate word with group
		if err := s.repo.AddWordToGroup(wordID, groupID); err != nil {
			continue
		}

		importedCount++
	}

	// Update group words count
	if err := s.repo.UpdateGroupWordsCount(groupID); err != nil {
		log.Error().Err(err).Msg("Failed to update group words count")
	}

	return &models.ImportWordsResponse{ImportedCount: importedCount}, nil
}
