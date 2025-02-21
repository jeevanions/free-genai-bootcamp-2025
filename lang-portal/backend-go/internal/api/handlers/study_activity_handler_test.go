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
)

type MockStudyActivityService struct {
	mock.Mock
}

func (m *MockStudyActivityService) GetStudyActivity(id int64) (*models.StudyActivityResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyActivityResponse), args.Error(1)
}

func (m *MockStudyActivityService) GetStudyActivitySessions(activityID int64) (*models.StudySessionsListResponse, error) {
	args := m.Called(activityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudySessionsListResponse), args.Error(1)
}

func TestStudyActivityHandler_GetStudyActivity(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockStudyActivityService)
	handler := NewStudyActivityHandler(mockService)

	t.Run("successful retrieval", func(t *testing.T) {
		// Arrange
		expectedActivity := &models.StudyActivityResponse{
			ID:           1,
			Name:         "Vocabulary Quiz",
			ThumbnailURL: "https://example.com/thumbnail.jpg",
			Description:  "Practice your vocabulary with flashcards",
			CreatedAt:    time.Now(),
		}
		mockService.On("GetStudyActivity", int64(1)).Return(expectedActivity, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		// Act
		handler.GetStudyActivity(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.StudyActivityResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedActivity.ID, response.ID)
		assert.Equal(t, expectedActivity.Name, response.Name)
		
		mockService.AssertExpectations(t)
	})

	t.Run("activity not found", func(t *testing.T) {
		// Arrange
		mockService.On("GetStudyActivity", int64(999)).Return(nil, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "999"}}

		// Act
		handler.GetStudyActivity(c)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		// Act
		handler.GetStudyActivity(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestStudyActivityHandler_GetStudyActivitySessions(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockStudyActivityService)
	handler := NewStudyActivityHandler(mockService)

	t.Run("successful retrieval", func(t *testing.T) {
		// Arrange
		expectedSessions := &models.StudySessionsListResponse{
			Items: []models.StudySessionResponse{
				{
					ID:           1,
					ActivityName: "Vocabulary Quiz",
					GroupID:      2,
					GroupName:    "Basic Greetings",
					CreatedAt:    time.Now(),
					WordsCount:   10,
					CorrectCount: 8,
				},
			},
		}
		mockService.On("GetStudyActivitySessions", int64(1)).Return(expectedSessions, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		// Act
		handler.GetStudyActivitySessions(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response models.StudySessionsListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Items, 1)
		assert.Equal(t, expectedSessions.Items[0].ID, response.Items[0].ID)
		
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		// Act
		handler.GetStudyActivitySessions(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		// Arrange
		mockService.On("GetStudyActivitySessions", int64(1)).Return(nil, assert.AnError).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		// Act
		handler.GetStudyActivitySessions(c)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
