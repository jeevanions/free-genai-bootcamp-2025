package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

type MockSettingsRepository struct {
	mock.Mock
}

func (m *MockSettingsRepository) ResetHistory() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSettingsRepository) DropAllTables() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSettingsRepository) CreateTables() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSettingsRepository) Close() error {
	return nil
}

func (m *MockSettingsRepository) CreateStudyActivitySession(activityID, groupID int64) (*models.LaunchStudyActivityResponse, error) {
	args := m.Called(activityID, groupID)
	return args.Get(0).(*models.LaunchStudyActivityResponse), args.Error(1)
}

func (m *MockSettingsRepository) CreateWordReview(sessionID, wordID int64, correct bool) error {
	args := m.Called(sessionID, wordID, correct)
	return args.Error(0)
}

func (m *MockSettingsRepository) GetAllStudySessions(limit, offset int) ([]models.StudySession, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StudySession), args.Error(1)
}

func (m *MockSettingsRepository) GetGroupByID(id int64) (*models.GroupDetailResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupDetailResponse), args.Error(1)
}

func (m *MockSettingsRepository) GetGroupStudySessions(groupID int64, limit, offset int) (*models.GroupStudySessionsResponse, error) {
	args := m.Called(groupID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupStudySessionsResponse), args.Error(1)
}

func (m *MockSettingsRepository) GetGroupWords(groupID int64, limit, offset int) (*models.GroupWordsResponse, error) {
	args := m.Called(groupID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupWordsResponse), args.Error(1)
}

func (m *MockSettingsRepository) GetGroups(limit, offset int) (*models.GroupListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupListResponse), args.Error(1)
}

func (m *MockSettingsRepository) GetLastStudySession() (*models.DashboardLastStudySession, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardLastStudySession), args.Error(1)
}

func (m *MockSettingsRepository) GetQuickStats() (*models.DashboardQuickStats, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardQuickStats), args.Error(1)
}

func (m *MockSettingsRepository) GetStudyProgress() (*models.DashboardStudyProgress, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardStudyProgress), args.Error(1)
}

func (m *MockSettingsRepository) GetStudyActivities(limit, offset int) ([]models.StudyActivity, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StudyActivity), args.Error(1)
}

func (m *MockSettingsRepository) GetStudyActivity(id int64) (*models.StudyActivity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyActivity), args.Error(1)
}

func (m *MockSettingsRepository) GetStudyActivitySessions(activityID int64, limit, offset int) ([]models.StudySession, error) {
	args := m.Called(activityID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StudySession), args.Error(1)
}

func (m *MockSettingsRepository) GetStudySessionWords(sessionID int64, limit, offset int) ([]*models.WordResponse, int, error) {
	args := m.Called(sessionID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.WordResponse), args.Int(1), args.Error(2)
}

func (m *MockSettingsRepository) GetTotalStudySessions() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockSettingsRepository) GetWordByID(id int64) (*models.WordResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WordResponse), args.Error(1)
}

func (m *MockSettingsRepository) GetWords(limit, offset int) (*models.WordListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WordListResponse), args.Error(1)
}

func (m *MockSettingsRepository) GetWordReviewsBySessionID(sessionID int64) ([]models.WordReviewItem, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.WordReviewItem), args.Error(1)
}

func (m *MockSettingsRepository) GetStudySessionDetail(sessionID int64) (*models.StudySession, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudySession), args.Error(1)
}

func (m *MockSettingsRepository) GetStudySessionStats(sessionID int64) (*models.StudySessionStats, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudySessionStats), args.Error(1)
}

func (m *MockSettingsRepository) GetTotalWordReviews() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockSettingsRepository) GetTotalWords() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockSettingsRepository) GetTotalGroups() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockSettingsRepository) GetTotalStudyActivities() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

type MockSeeder struct {
	mock.Mock
}

func (m *MockSeeder) SeedFromJSON(seedDir string) error {
	args := m.Called(seedDir)
	return args.Error(0)
}

func TestSettingsService_ResetHistory(t *testing.T) {
	t.Run("successful reset", func(t *testing.T) {
		mockRepo := new(MockSettingsRepository)
		mockSeeder := new(MockSeeder)
		service := NewSettingsService(mockRepo, mockSeeder)

		mockRepo.On("ResetHistory").Return(nil)
		err := service.ResetHistory()
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockSettingsRepository)
		mockSeeder := new(MockSeeder)
		service := NewSettingsService(mockRepo, mockSeeder)

		expectedErr := errors.New("database error")
		mockRepo.On("ResetHistory").Return(expectedErr)

		err := service.ResetHistory()
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSettingsService_FullReset(t *testing.T) {
	mockRepo := &MockSettingsRepository{}
	mockSeeder := &MockSeeder{}
	service := NewSettingsService(mockRepo, mockSeeder)

	t.Run("successful reset", func(t *testing.T) {
		mockRepo.On("DropAllTables").Return(nil)
		mockRepo.On("CreateTables").Return(nil)
		mockSeeder.On("SeedFromJSON", "internal/db/seeds").Return(nil)

		err := service.FullReset()

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockSeeder.AssertExpectations(t)
	})

	t.Run("drop tables error", func(t *testing.T) {
		mockRepo.On("DropAllTables").Return(errors.New("drop error"))

		err := service.FullReset()

		assert.Error(t, err)
		assert.Equal(t, "drop error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("create tables error", func(t *testing.T) {
		mockRepo.On("DropAllTables").Return(nil)
		mockRepo.On("CreateTables").Return(errors.New("create error"))

		err := service.FullReset()

		assert.Error(t, err)
		assert.Equal(t, "create error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("seed error", func(t *testing.T) {
		mockRepo.On("DropAllTables").Return(nil)
		mockRepo.On("CreateTables").Return(nil)
		mockSeeder.On("SeedFromJSON", "internal/db/seeds").Return(errors.New("seed error"))

		err := service.FullReset()

		assert.Error(t, err)
		assert.Equal(t, "seed error", err.Error())
		mockRepo.AssertExpectations(t)
		mockSeeder.AssertExpectations(t)
	})
}
