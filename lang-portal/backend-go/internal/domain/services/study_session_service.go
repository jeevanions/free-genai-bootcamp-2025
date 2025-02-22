package services

import (
	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

type StudySessionServiceInterface interface {
	GetAllStudySessions(limit, offset int) (*models.StudySessionListResponse, error)
}

type StudySessionService struct {
	repo repository.Repository
}

func NewStudySessionService(repo repository.Repository) *StudySessionService {
	return &StudySessionService{repo: repo}
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

	// Create response with detailed information
	response := &models.StudySessionListResponse{
		Items: make([]models.StudySessionDetailResponse, 0, len(sessions)),
		Pagination: models.PaginationResponse{
			CurrentPage:  (offset / limit) + 1,
			TotalPages:   (total + limit - 1) / limit,
			TotalItems:   total,
			ItemsPerPage: limit,
		},
	}

	// Populate detailed information for each session
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
