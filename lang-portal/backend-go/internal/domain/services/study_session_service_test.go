package services

import (
	"errors"
	"testing"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetStudySessionWords(sessionID int64, limit, offset int) ([]*models.WordResponse, int, error) {
	args := m.Called(sessionID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.WordResponse), args.Int(1), args.Error(2)
}

func (m *MockRepository) GetAllStudySessions(limit, offset int) ([]models.StudySession, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StudySession), args.Error(1)
}

func (m *MockRepository) GetTotalStudySessions() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetStudyActivity(id int64) (*models.StudyActivity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyActivity), args.Error(1)
}

func (m *MockRepository) GetGroupByID(id int64) (*models.GroupDetailResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupDetailResponse), args.Error(1)
}

func (m *MockRepository) GetWordReviewsBySessionID(sessionID int64) ([]models.WordReviewItem, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.WordReviewItem), args.Error(1)
}

func (m *MockRepository) CreateWordReview(sessionID, wordID int64, correct bool) error {
	args := m.Called(sessionID, wordID, correct)
	return args.Error(0)
}

func (m *MockRepository) CreateStudyActivitySession(activityID, groupID int64) (*models.LaunchStudyActivityResponse, error) {
	args := m.Called(activityID, groupID)
	return args.Get(0).(*models.LaunchStudyActivityResponse), args.Error(1)
}

func (m *MockRepository) CreateTables() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository) DropAllTables() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository) GetGroupStudySessions(groupID int64, limit, offset int) (*models.GroupStudySessionsResponse, error) {
	args := m.Called(groupID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupStudySessionsResponse), args.Error(1)
}

func (m *MockRepository) GetGroupWords(groupID int64, limit, offset int) (*models.GroupWordsResponse, error) {
	args := m.Called(groupID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupWordsResponse), args.Error(1)
}

func (m *MockRepository) GetGroups(limit, offset int) (*models.GroupListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupListResponse), args.Error(1)
}

func (m *MockRepository) Close() error {
	return nil
}

func (m *MockRepository) GetLastStudySession() (*models.DashboardLastStudySession, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardLastStudySession), args.Error(1)
}

func (m *MockRepository) GetQuickStats() (*models.DashboardQuickStats, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardQuickStats), args.Error(1)
}

func (m *MockRepository) GetStudyProgress() (*models.DashboardStudyProgress, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardStudyProgress), args.Error(1)
}

func (m *MockRepository) GetStudyActivities(limit, offset int) ([]models.StudyActivity, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StudyActivity), args.Error(1)
}

func (m *MockRepository) GetStudyActivitySessions(activityID int64, limit, offset int) ([]models.StudySession, error) {
	args := m.Called(activityID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StudySession), args.Error(1)
}

func (m *MockRepository) GetWordByID(id int64) (*models.WordResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WordResponse), args.Error(1)
}

func (m *MockRepository) GetWords(limit, offset int) (*models.WordListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WordListResponse), args.Error(1)
}

func (m *MockRepository) ResetHistory() error {
	args := m.Called()
	return args.Error(0)
}

func TestStudySessionService_GetStudySessionWords(t *testing.T) {
	mockRepo := new(MockRepository)
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
	mockRepo := new(MockRepository)
	service := NewStudySessionService(mockRepo)

	t.Run("successful retrieval", func(t *testing.T) {
		sessions := []*models.StudySession{
			{
				ID:              1,
				GroupID:         1,
				StudyActivityID: 1,
			},
		}
		activity := &models.StudyActivity{
			ID:           1,
			Name:         "Test Activity",
			Description:  "Test Description",
			ThumbnailURL: "http://example.com/thumb.jpg",
			LaunchURL:    models.StrPtr("http://example.com/launch"),
		}
		group := &models.Group{
			ID:         1,
			Name:       "Test Group",
			WordsCount: 10,
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

	t.Run("repository error - GetAllStudySessions", func(t *testing.T) {
		mockRepo.On("GetAllStudySessions", 10, 0).Return(nil, errors.New("repository error"))

		response, err := service.GetAllStudySessions(10, 0)

		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error - GetTotalStudySessions", func(t *testing.T) {
		sessions := []*models.StudySession{{ID: 1}}
		mockRepo.On("GetAllStudySessions", 10, 0).Return(sessions, nil)
		mockRepo.On("GetTotalStudySessions").Return(0, errors.New("repository error"))

		response, err := service.GetAllStudySessions(10, 0)

		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestStudySessionService_ReviewWord(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewStudySessionService(mockRepo)

	t.Run("successful review", func(t *testing.T) {
		review := &models.WordReviewResponse{
			Success: true,
			WordID:  1,
		}

		mockRepo.On("CreateWordReview", int64(1), int64(1), true).Return(nil)

		response, err := service.ReviewWord(1, 1, true)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, review.WordID, response.WordID)
		assert.Equal(t, review.Success, response.Success)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.On("CreateWordReview", int64(1), int64(1), true).Return(nil, errors.New("repository error"))

		response, err := service.ReviewWord(1, 1, true)

		assert.Error(t, err)
		assert.Nil(t, response)
		mockRepo.AssertExpectations(t)
	})
}
