package mocks

import (
	"database/sql"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
	"github.com/stretchr/testify/mock"
)

// MockRepository implements the repository.Repository interface for testing
type MockRepository struct {
	mock.Mock
}

// Dashboard operations
func (m *MockRepository) GetLastStudySession() (*models.DashboardLastStudySession, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardLastStudySession), args.Error(1)
}

func (m *MockRepository) GetStudyProgress() (*models.DashboardStudyProgress, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardStudyProgress), args.Error(1)
}

func (m *MockRepository) GetQuickStats() (*models.DashboardQuickStats, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardQuickStats), args.Error(1)
}

// Study activities
func (m *MockRepository) GetStudyActivities(limit, offset int) (*models.StudyActivityListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyActivityListResponse), args.Error(1)
}

func (m *MockRepository) GetStudyActivity(id int64) (*models.StudyActivityResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyActivityResponse), args.Error(1)
}

// Study sessions
func (m *MockRepository) GetAllStudySessions(limit, offset int) ([]models.StudySession, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StudySession), args.Error(1)
}

func (m *MockRepository) GetTotalStudySessions() (int, error) {
	args := m.Called()
	return args.Get(0).(int), args.Error(1)
}

func (m *MockRepository) GetStudySessionWords(sessionID int64, limit, offset int) ([]*models.WordResponse, int, error) {
	args := m.Called(sessionID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*models.WordResponse), args.Get(1).(int), args.Error(2)
}

func (m *MockRepository) GetStudyActivitySessions(activityID int64, limit, offset int) ([]models.StudySession, error) {
	args := m.Called(activityID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StudySession), args.Error(1)
}

func (m *MockRepository) CreateStudyActivitySession(activityID, groupID int64) (*models.LaunchStudyActivityResponse, error) {
	args := m.Called(activityID, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LaunchStudyActivityResponse), args.Error(1)
}

func (m *MockRepository) GetWordReviewsBySessionID(sessionID int64) ([]models.WordReviewItem, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.WordReviewItem), args.Error(1)
}

// Transaction operations
func (m *MockRepository) BeginTx() (*sql.Tx, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockRepository) CommitTx(tx *sql.Tx) error {
	return m.Called(tx).Error(0)
}

func (m *MockRepository) RollbackTx(tx *sql.Tx) error {
	return m.Called(tx).Error(0)
}

// Word operations
func (m *MockRepository) GetWords(limit, offset int) (*models.WordListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WordListResponse), args.Error(1)
}

func (m *MockRepository) GetWordByID(id int64) (*models.WordResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WordResponse), args.Error(1)
}

func (m *MockRepository) CreateWord(word *models.WordResponse) (int64, error) {
	args := m.Called(word)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) CreateWordReview(sessionID, wordID int64, correct bool) error {
	return m.Called(sessionID, wordID, correct).Error(0)
}

// Group operations
func (m *MockRepository) CreateGroup(name string) (int64, error) {
	args := m.Called(name)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) GetGroups(limit, offset int) (*models.GroupListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupListResponse), args.Error(1)
}

func (m *MockRepository) GetGroupByID(id int64) (*models.GroupDetailResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupDetailResponse), args.Error(1)
}

func (m *MockRepository) GetGroupIDByName(name string) (int64, error) {
	args := m.Called(name)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) GetGroupWords(groupID int64, limit, offset int) (*models.GroupWordsResponse, error) {
	args := m.Called(groupID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupWordsResponse), args.Error(1)
}

func (m *MockRepository) GetGroupStudySessions(groupID int64, limit, offset int) (*models.GroupStudySessionsResponse, error) {
	args := m.Called(groupID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupStudySessionsResponse), args.Error(1)
}

func (m *MockRepository) UpdateGroupWordsCount(groupID int64) error {
	return m.Called(groupID).Error(0)
}

// Other operations
func (m *MockRepository) AddWordToGroup(wordID, groupID int64) error {
	return m.Called(wordID, groupID).Error(0)
}

// Settings/Reset operations
func (m *MockRepository) ResetHistory() error {
	return m.Called().Error(0)
}

func (m *MockRepository) Close() error {
	return m.Called().Error(0)
}

func (m *MockRepository) DropAllTables() error {
	return m.Called().Error(0)
}

func (m *MockRepository) CreateTables() error {
	return m.Called().Error(0)
}
