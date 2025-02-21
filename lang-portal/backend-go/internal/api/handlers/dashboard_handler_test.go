package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/services"
)

// MockDashboardService is a mock implementation of the dashboard service
// MockDashboardService is a mock implementation of services.DashboardServiceInterface
type MockDashboardService struct {
	mock.Mock
}

// Verify that MockDashboardService implements services.DashboardServiceInterface
var _ services.DashboardServiceInterface = (*MockDashboardService)(nil)

func (m *MockDashboardService) GetLastStudySession() (*models.DashboardLastStudySession, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardLastStudySession), args.Error(1)
}

func (m *MockDashboardService) GetStudyProgress() (*models.DashboardStudyProgress, error) {
	args := m.Called()
	return args.Get(0).(*models.DashboardStudyProgress), args.Error(1)
}

func (m *MockDashboardService) GetQuickStats() (*models.DashboardQuickStats, error) {
	args := m.Called()
	return args.Get(0).(*models.DashboardQuickStats), args.Error(1)
}

func TestDashboardHandler_GetLastStudySession(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockDashboardService)
	handler := NewDashboardHandler(mockService)

	t.Run("successful retrieval", func(t *testing.T) {
		// Arrange
		expectedSession := &models.DashboardLastStudySession{
			ID:              1,
			GroupID:         2,
			CreatedAt:       time.Now(),
			StudyActivityID: 3,
			GroupName:       "Test Group",
		}
		mockService.On("GetLastStudySession").Return(expectedSession, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Act
		handler.GetLastStudySession(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.DashboardLastStudySession
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedSession.ID, response.ID)
		assert.Equal(t, expectedSession.GroupName, response.GroupName)
		
		mockService.AssertExpectations(t)
	})

	t.Run("no sessions found", func(t *testing.T) {
		// Arrange
		mockService.On("GetLastStudySession").Return(nil, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Act
		handler.GetLastStudySession(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDashboardHandler_GetStudyProgress(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockDashboardService)
	handler := NewDashboardHandler(mockService)

	t.Run("successful retrieval", func(t *testing.T) {
		// Arrange
		expectedProgress := &models.DashboardStudyProgress{
			TotalWordsStudied:    10,
			TotalAvailableWords: 100,
		}
		mockService.On("GetStudyProgress").Return(expectedProgress, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Act
		handler.GetStudyProgress(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.DashboardStudyProgress
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedProgress.TotalWordsStudied, response.TotalWordsStudied)
		assert.Equal(t, expectedProgress.TotalAvailableWords, response.TotalAvailableWords)
		
		mockService.AssertExpectations(t)
	})
}

func TestDashboardHandler_GetQuickStats(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockDashboardService)
	handler := NewDashboardHandler(mockService)

	t.Run("successful retrieval", func(t *testing.T) {
		// Arrange
		expectedStats := &models.DashboardQuickStats{
			SuccessRate:        80.0,
			TotalStudySessions: 10,
			TotalActiveGroups:  3,
			StudyStreakDays:    5,
		}
		mockService.On("GetQuickStats").Return(expectedStats, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Act
		handler.GetQuickStats(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.DashboardQuickStats
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedStats.SuccessRate, response.SuccessRate)
		assert.Equal(t, expectedStats.TotalStudySessions, response.TotalStudySessions)
		assert.Equal(t, expectedStats.TotalActiveGroups, response.TotalActiveGroups)
		assert.Equal(t, expectedStats.StudyStreakDays, response.StudyStreakDays)
		
		mockService.AssertExpectations(t)
	})
}
