package services

import (
	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

type StudySessionServiceInterface interface {
	GetAllStudySessions(limit, offset int) (*models.StudySessionListResponse, error)
	GetStudySessionWords(sessionID int64, limit, offset int) (*models.StudySessionWordsResponse, error)
	ReviewWord(sessionID, wordID int64, correct bool) (*models.WordReviewResponse, error)
}

type StudySessionService struct {
	repo repository.Repository
}

func NewStudySessionService(repo repository.Repository) *StudySessionService {
	return &StudySessionService{repo: repo}
}

// GetStudySessionWords returns a paginated list of words reviewed in a study session
func (s *StudySessionService) GetStudySessionWords(sessionID int64, limit, offset int) (*models.StudySessionWordsResponse, error) {
	words, total, err := s.repo.GetStudySessionWords(sessionID, limit, offset)
	if err != nil {
		return nil, err
	}

	totalPages := (total + limit - 1) / limit
	return &models.StudySessionWordsResponse{
		Items: words,
		Pagination: models.PaginationResponse{
			CurrentPage:   (offset / limit) + 1,
			TotalPages:    totalPages,
			TotalItems:    total,
			ItemsPerPage:  limit,
		},
	}, nil
}

func (s *StudySessionService) GetAllStudySessions(limit, offset int) (*models.StudySessionListResponse, error) {
	// Get study sessions
	sessions, err := s.repo.GetAllStudySessions(limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count for pagination
	total, err := s.repo.GetTotalStudySessions()
	if err != nil {
		return nil, err
	}

	// Build response with pagination
	totalPages := (total + limit - 1) / limit
	response := &models.StudySessionListResponse{
		Items: []models.StudySessionDetailResponse{},
		Pagination: models.PaginationResponse{
			CurrentPage:  (offset / limit) + 1,
			TotalPages:   totalPages,
			TotalItems:   total,
			ItemsPerPage: limit,
		},
	}

	// Process each study session
	for _, session := range sessions {
		// Get activity details
		activity, err := s.repo.GetStudyActivity(session.StudyActivityID)
		if err != nil {
			return nil, err
		}

		// Get group details
		group, err := s.repo.GetGroupByID(session.GroupID)
		if err != nil {
			return nil, err
		}

		// Get word reviews
		reviews, err := s.repo.GetWordReviewsBySessionID(session.ID)
		if err != nil {
			return nil, err
		}

		// Calculate stats
		var stats models.StudySessionStats
		if len(reviews) > 0 {
			var correctCount int
			for _, review := range reviews {
				if review.Correct {
					correctCount++
				}
			}
			stats.TotalWords = len(reviews)
			stats.CorrectWords = correctCount
			stats.SuccessRate = float64(correctCount) / float64(len(reviews)) * 100
		}

		// Create detailed response
		detailedSession := models.StudySessionDetailResponse{
			ID:           session.ID,
			ActivityName: activity.Name,
			GroupName:    group.Name,
			CreatedAt:    session.CreatedAt,
			Stats:        stats,
			ReviewItems:  reviews,
		}

		response.Items = append(response.Items, detailedSession)
	}

	return response, nil
}

// ReviewWord records a word review in a study session
func (s *StudySessionService) ReviewWord(sessionID, wordID int64, correct bool) (*models.WordReviewResponse, error) {
	// Create the word review
	err := s.repo.CreateWordReview(sessionID, wordID, correct)
	if err != nil {
		return nil, err
	}

	return &models.WordReviewResponse{
		Success: true,
		WordID:  wordID,
	}, nil
}
