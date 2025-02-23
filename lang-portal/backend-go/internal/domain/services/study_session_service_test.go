package services

import (
	"errors"
	"testing"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
	"github.com/jeevanions/lang-portal/backend-go/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestStudySessionService_GetStudySessionWords(t *testing.T) {
	mockRepo := &mocks.MockRepository{}
	service := NewStudySessionService(mockRepo)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedWords := []*models.WordResponse{
			{
				ID:      1,
				Italian: "ciao",
				English: "hello",
				Parts:   map[string]interface{}{"part": "greeting"},
			},
		}
		totalWords := 1

		mockRepo.On("GetStudySessionWords", int64(1), 10, 0).Return(expectedWords, totalWords, nil)

		response, err := service.GetStudySessionWords(1, 10, 0)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedWords[0].ID, response.Items[0].ID)
		assert.Equal(t, 1, response.Pagination.TotalPages)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.On("GetStudySessionWords", int64(1), 10, 0).Return(nil, 0, errors.New("repository error"))

		response, err := service.GetStudySessionWords(1, 10, 0)

		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestStudySessionService_GetAllStudySessions(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		// Create new mock for each test
		mockRepo := new(mocks.MockRepository)
		service := NewStudySessionService(mockRepo)

		sessions := []models.StudySession{
			{
				ID:              1,
				GroupID:         1,
				StudyActivityID: 1,
			},
		}
		activity := &models.StudyActivityResponse{
			ID:           1,
			Name:         "Test Activity",
			Description:  models.StrPtr("Test Description"),
			ThumbnailURL: models.StrPtr("http://example.com/thumb.jpg"),
			LaunchURL:    models.StrPtr("http://example.com/launch"),
		}
		group := &models.GroupDetailResponse{
			ID:    1,
			Name:  "Test Group",
			Stats: models.GroupStats{TotalWordCount: 10},
		}
		reviews := []models.WordReviewItem{{WordID: 1, Correct: true}}

		mockRepo.On("GetAllStudySessions", 10, 0).Return(sessions, nil)
		mockRepo.On("GetTotalStudySessions").Return(1, nil)
		mockRepo.On("GetStudyActivity", int64(1)).Return(activity, nil)
		mockRepo.On("GetGroupByID", int64(1)).Return(group, nil)
		mockRepo.On("GetWordReviewsBySessionID", int64(1)).Return(reviews, nil)

		response, err := service.GetAllStudySessions(10, 0)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, len(response.Items))
		assert.Equal(t, sessions[0].ID, response.Items[0].ID)
		assert.Equal(t, activity.Name, response.Items[0].ActivityName)
		assert.Equal(t, group.Name, response.Items[0].GroupName)
		assert.Equal(t, 100.0, response.Items[0].Stats.SuccessRate)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty study sessions list", func(t *testing.T) {
		// Create new mock for each test
		mockRepo := new(mocks.MockRepository)
		service := NewStudySessionService(mockRepo)

		// Only set up expectations needed for this test
		mockRepo.On("GetAllStudySessions", 10, 0).Return([]models.StudySession{}, nil)
		mockRepo.On("GetTotalStudySessions").Return(0, nil)

		response, err := service.GetAllStudySessions(10, 0)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Empty(t, response.Items)
		assert.Equal(t, models.PaginationResponse{
			CurrentPage:  1,
			TotalPages:   0,
			TotalItems:   0,
			ItemsPerPage: 10,
		}, response.Pagination)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error - GetTotalStudySessions", func(t *testing.T) {
		// Create new mock for each test
		mockRepo := new(mocks.MockRepository)
		service := NewStudySessionService(mockRepo)

		// Only set up expectations needed for this test
		mockRepo.On("GetAllStudySessions", 10, 0).Return([]models.StudySession{}, nil)
		mockRepo.On("GetTotalStudySessions").Return(0, errors.New("repository error"))

		response, err := service.GetAllStudySessions(10, 0)

		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestStudySessionService_ReviewWord(t *testing.T) {
	mockRepo := &mocks.MockRepository{}
	service := NewStudySessionService(mockRepo)

	t.Run("successful review", func(t *testing.T) {
		mockRepo.On("CreateWordReview", int64(1), int64(1), true).Return(nil)

		response, err := service.ReviewWord(1, 1, true)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, int64(1), response.WordID)
		assert.True(t, response.Success)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.On("CreateWordReview", int64(1), int64(1), true).Return(errors.New("repository error"))

		response, err := service.ReviewWord(1, 1, true)

		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}
