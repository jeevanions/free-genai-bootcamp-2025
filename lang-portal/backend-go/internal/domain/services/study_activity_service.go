package services

import (
	"fmt"

	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

type StudyActivityServiceInterface interface {
	GetStudyActivities(limit, offset int) (*models.StudyActivityListResponse, error)
	GetStudyActivity(id int64) (*models.StudyActivityResponse, error)
	GetStudyActivitySessions(activityID int64) (*models.StudySessionsListResponse, error)
	LaunchStudyActivity(activityID, groupID int64) (*models.LaunchStudyActivityResponse, error)
}

type StudyActivityService struct {
	repo repository.Repository
}

func NewStudyActivityService(repo repository.Repository) StudyActivityServiceInterface {
	return &StudyActivityService{repo: repo}
}

func (s *StudyActivityService) GetStudyActivities(limit, offset int) (*models.StudyActivityListResponse, error) {
	activities, err := s.repo.GetStudyActivities(limit, offset)
	if err != nil {
		return nil, err
	}

	response := &models.StudyActivityListResponse{
		Items: make([]models.StudyActivityResponse, 0, len(activities)),
		Pagination: models.PaginationResponse{
			CurrentPage: offset/limit + 1,
			ItemsPerPage: limit,
		},
	}

	for _, activity := range activities {
		response.Items = append(response.Items, models.StudyActivityResponse{
			ID:           activity.ID,
			Name:         activity.Name,
			ThumbnailURL: activity.ThumbnailURL,
			Description:  activity.Description,
			LaunchURL:    activity.LaunchURL,
			CreatedAt:    activity.CreatedAt,
		})
	}

	return response, nil
}

func (s *StudyActivityService) GetStudyActivity(id int64) (*models.StudyActivityResponse, error) {
	activity, err := s.repo.GetStudyActivity(id)
	if err != nil {
		return nil, err
	}

	if activity == nil {
		return nil, nil
	}

	return &models.StudyActivityResponse{
		ID:           activity.ID,
		Name:         activity.Name,
		ThumbnailURL: activity.ThumbnailURL,
		Description:  activity.Description,
		CreatedAt:    activity.CreatedAt,
	}, nil
}

func (s *StudyActivityService) LaunchStudyActivity(activityID, groupID int64) (*models.LaunchStudyActivityResponse, error) {
	// Verify activity exists
	activity, err := s.repo.GetStudyActivity(activityID)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		return nil, fmt.Errorf("study activity not found")
	}

	// Create study session
	return s.repo.CreateStudyActivitySession(activityID, groupID)
}

func (s *StudyActivityService) GetStudyActivitySessions(activityID int64) (*models.StudySessionsListResponse, error) {
	sessions, err := s.repo.GetStudyActivitySessions(activityID, 100, 0) // Using limit=100 as per spec
	if err != nil {
		return nil, err
	}

	response := &models.StudySessionsListResponse{
		Items: make([]models.StudySessionResponse, 0, len(sessions)),
	}

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

		// Get session stats
		wordReviews, err := s.repo.GetWordReviewsBySessionID(session.ID)
		if err != nil {
			return nil, err
		}

		// Calculate stats
		wordsCount := len(wordReviews)
		correctCount := 0
		for _, review := range wordReviews {
			if review.Correct {
				correctCount++
			}
		}

		response.Items = append(response.Items, models.StudySessionResponse{
			ID:           session.ID,
			ActivityName: activity.Name,
			GroupID:      session.GroupID,
			GroupName:    group.Name,
			CreatedAt:    session.CreatedAt,
			WordsCount:   wordsCount,
			CorrectCount: correctCount,
		})
	}

	return response, nil
}
