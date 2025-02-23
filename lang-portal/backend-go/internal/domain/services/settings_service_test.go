package services

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *MockSettingsRepository) AddWordToGroup(wordID, groupID int64) error {
	args := m.Called(wordID, groupID)
	return args.Error(0)
}

func (m *MockSettingsRepository) BeginTx() (*sql.Tx, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockSettingsRepository) CreateGroup(name string) (int64, error) {
	args := m.Called(name)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSettingsRepository) GetGroupIDByName(name string) (int64, error) {
	args := m.Called(name)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSettingsRepository) CreateWord(word *models.Word) (int64, error) {
	args := m.Called(word)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSettingsRepository) UpdateGroupWordsCount(groupID int64) error {
	args := m.Called(groupID)
	return args.Error(0)
}

func (m *MockSettingsRepository) CreateStudyActivity(activity *models.StudyActivity) (int64, error) {
	args := m.Called(activity)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSettingsRepository) UpdateStudyActivity(activity *models.StudyActivity) error {
	args := m.Called(activity)
	return args.Error(0)
}

func (m *MockSettingsRepository) DeleteStudyActivity(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSettingsRepository) DeleteGroup(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSettingsRepository) UpdateGroup(group *models.Group) error {
	args := m.Called(group)
	return args.Error(0)
}

func (m *MockSettingsRepository) RemoveWordFromGroup(wordID, groupID int64) error {
	args := m.Called(wordID, groupID)
	return args.Error(0)
}

func (m *MockSettingsRepository) UpdateWord(word *models.Word) error {
	args := m.Called(word)
	return args.Error(0)
}

func (m *MockSettingsRepository) DeleteWord(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSettingsRepository) CommitTx(tx *sql.Tx) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockSettingsRepository) RollbackTx(tx *sql.Tx) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockSettingsRepository) GetStudyActivityByID(id int64) (*models.StudyActivity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyActivity), args.Error(1)
}

func (m *MockSettingsRepository) GetStudyActivityList(limit, offset int) (*models.StudyActivityListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyActivityListResponse), args.Error(1)
}

func (m *MockSettingsRepository) GetWordsByGroupID(groupID int64, limit, offset int) ([]*models.Word, int, error) {
	args := m.Called(groupID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.Word), args.Int(1), args.Error(2)
}

func (m *MockSettingsRepository) GetWordList(limit, offset int) (*models.WordListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WordListResponse), args.Error(1)
}

func (m *MockSettingsRepository) GetGroupList(limit, offset int) (*models.GroupListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupListResponse), args.Error(1)
}

func (m *MockSettingsRepository) GetGroupByID(id int64) (*models.Group, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Group), args.Error(1)
}

func (m *MockSettingsRepository) GetWordByID(id int64) (*models.Word, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Word), args.Error(1)
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
		mockTx := &sql.Tx{}
		mockRepo.On("BeginTx").Return(mockTx, nil)
		mockRepo.On("DropAllTables").Return(nil)
		mockRepo.On("CreateTables").Return(nil)
		mockRepo.On("CreateGroup", mock.AnythingOfType("string")).Return(int64(1), nil)
		mockRepo.On("GetGroupIDByName", mock.AnythingOfType("string")).Return(int64(1), nil)
		mockSeeder.On("SeedFromJSON", "internal/db/seeds").Return(nil)

		err := service.FullReset()

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockSeeder.AssertExpectations(t)
	})

	t.Run("begin transaction error", func(t *testing.T) {
		mockRepo.On("BeginTx").Return(nil, errors.New("transaction error"))

		err := service.FullReset()

		assert.Error(t, err)
		assert.Equal(t, "transaction error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("drop tables error", func(t *testing.T) {
		mockTx := &sql.Tx{}
		mockRepo.On("BeginTx").Return(mockTx, nil)
		mockRepo.On("DropAllTables").Return(errors.New("drop error"))

		err := service.FullReset()

		assert.Error(t, err)
		assert.Equal(t, "drop error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("create tables error", func(t *testing.T) {
		mockTx := &sql.Tx{}
		mockRepo.On("BeginTx").Return(mockTx, nil)
		mockRepo.On("DropAllTables").Return(nil)
		mockRepo.On("CreateTables").Return(errors.New("create error"))

		err := service.FullReset()

		assert.Error(t, err)
		assert.Equal(t, "create error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("create group error", func(t *testing.T) {
		mockTx := &sql.Tx{}
		mockRepo.On("BeginTx").Return(mockTx, nil)
		mockRepo.On("DropAllTables").Return(nil)
		mockRepo.On("CreateTables").Return(nil)
		mockRepo.On("CreateGroup", mock.AnythingOfType("string")).Return(int64(0), errors.New("group error"))

		err := service.FullReset()

		assert.Error(t, err)
		assert.Equal(t, "group error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("seed error", func(t *testing.T) {
		mockTx := &sql.Tx{}
		mockRepo.On("BeginTx").Return(mockTx, nil)
		mockRepo.On("DropAllTables").Return(nil)
		mockRepo.On("CreateTables").Return(nil)
		mockRepo.On("CreateGroup", mock.AnythingOfType("string")).Return(int64(1), nil)
		mockRepo.On("GetGroupIDByName", mock.AnythingOfType("string")).Return(int64(1), nil)
		mockSeeder.On("SeedFromJSON", "internal/db/seeds").Return(errors.New("seed error"))

		err := service.FullReset()

		assert.Error(t, err)
		assert.Equal(t, "seed error", err.Error())
		mockRepo.AssertExpectations(t)
		mockSeeder.AssertExpectations(t)
	})
}

func TestSettingsService_CreateGroup(t *testing.T) {
	mockRepo := &MockSettingsRepository{}
	mockSeeder := &MockSeeder{}
	service := NewSettingsService(mockRepo, mockSeeder)

	t.Run("successful group creation", func(t *testing.T) {
		groupName := "Test Group"
		expectedID := int64(1)

		mockRepo.On("CreateGroup", groupName).Return(expectedID, nil)
		mockRepo.On("UpdateGroupWordsCount", expectedID).Return(nil)

		id, err := service.CreateGroup(groupName)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("create group error", func(t *testing.T) {
		groupName := "Test Group"

		mockRepo.On("CreateGroup", groupName).Return(int64(0), errors.New("creation error"))

		id, err := service.CreateGroup(groupName)

		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.Equal(t, "creation error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("update words count error", func(t *testing.T) {
		groupName := "Test Group"
		expectedID := int64(1)

		mockRepo.On("CreateGroup", groupName).Return(expectedID, nil)
		mockRepo.On("UpdateGroupWordsCount", expectedID).Return(errors.New("update error"))

		id, err := service.CreateGroup(groupName)

		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.Equal(t, "update error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestSettingsService_GetStudyActivities(t *testing.T) {
	mockRepo := &MockSettingsRepository{}
	mockSeeder := &MockSeeder{}
	service := NewSettingsService(mockRepo, mockSeeder)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedResponse := &models.StudyActivityListResponse{
			Items: []models.StudyActivityResponse{
				{
					ID:           1,
					Name:         "Test Activity",
					Description:  "Test Description",
					ThumbnailURL: "http://example.com/thumb.jpg",
				},
			},
			Pagination: models.PaginationResponse{
				CurrentPage:  1,
				TotalPages:   1,
				TotalItems:   1,
				ItemsPerPage: 10,
			},
		}

		mockRepo.On("GetStudyActivities", 10, 0).Return(expectedResponse, nil)

		response, err := service.GetStudyActivities(10, 0)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.On("GetStudyActivities", 10, 0).Return(nil, errors.New("repository error"))

		response, err := service.GetStudyActivities(10, 0)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "repository error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
