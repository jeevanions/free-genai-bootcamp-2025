package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

func strPtr(s string) *string {
	return &s
}

type MockStudyActivityService struct {
	mock.Mock
}

func (m *MockStudyActivityService) GetStudyActivities(limit, offset int) (*models.StudyActivityListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyActivityListResponse), args.Error(1)
}

func (m *MockStudyActivityService) GetStudyActivity(id int64) (*models.StudyActivityResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyActivityResponse), args.Error(1)
}

func (m *MockStudyActivityService) LaunchStudyActivity(activityID, groupID int64) (*models.LaunchStudyActivityResponse, error) {
	args := m.Called(activityID, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LaunchStudyActivityResponse), args.Error(1)
}

func (m *MockStudyActivityService) GetStudyActivitySessions(activityID int64) (*models.StudySessionsListResponse, error) {
	args := m.Called(activityID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudySessionsListResponse), args.Error(1)
}

func TestStudyActivityHandler_GetStudyActivities(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockStudyActivityService)
	handler := NewStudyActivityHandler(mockService)

	t.Run("successful retrieval", func(t *testing.T) {
		// Arrange
		expectedActivities := &models.StudyActivityListResponse{
			Items: []models.StudyActivityResponse{
				{
					ID:           1,
					Name:         "Vocabulary Quiz",
					ThumbnailURL: "https://example.com/thumbnail.jpg",
					Description:  "Practice your vocabulary with flashcards",
					LaunchURL:    strPtr("https://example.com/quiz/launch"),
					CreatedAt:    time.Now(),
				},
			},
			Pagination: models.PaginationResponse{
				CurrentPage:  1,
				TotalPages:   1,
				TotalItems:   1,
				ItemsPerPage: 100,
			},
		}

		mockService.On("GetStudyActivities", 100, 0).Return(expectedActivities, nil)

		// Create request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/study_activities", nil)

		// Act
		handler.GetStudyActivities(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.StudyActivityListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedActivities.Items[0].ID, response.Items[0].ID)
		assert.Equal(t, expectedActivities.Items[0].Name, response.Items[0].Name)
		assert.Equal(t, expectedActivities.Pagination.CurrentPage, response.Pagination.CurrentPage)
	})

	t.Run("empty list", func(t *testing.T) {
		// Arrange
		emptyResponse := &models.StudyActivityListResponse{
			Items: []models.StudyActivityResponse{},
			Pagination: models.PaginationResponse{
				CurrentPage:  1,
				TotalPages:   1,
				TotalItems:   0,
				ItemsPerPage: 100,
			},
		}

		mockService.On("GetStudyActivities", 100, 0).Return(emptyResponse, nil)

		// Create request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/study_activities", nil)

		// Act
		handler.GetStudyActivities(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.StudyActivityListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Empty(t, response.Items)
		assert.Equal(t, 1, response.Pagination.CurrentPage)
	})
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

func TestStudyActivityHandler_LaunchStudyActivity(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockStudyActivityService)
	handler := NewStudyActivityHandler(mockService)

	t.Run("successful launch", func(t *testing.T) {
		// Arrange
		expectedResponse := &models.LaunchStudyActivityResponse{
			StudySessionID:  456,
			StudyActivityID: 789,
			GroupID:        123,
			CreatedAt:      time.Now(),
		}

		mockService.On("LaunchStudyActivity", int64(1), int64(123)).Return(expectedResponse, nil)

		// Create request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Create request body
		reqBody := models.LaunchStudyActivityRequest{GroupID: 123}
		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/study_activities/1/launch", bytes.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "1")

		// Act
		handler.LaunchStudyActivity(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.LaunchStudyActivityResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.StudySessionID, response.StudySessionID)
		assert.Equal(t, expectedResponse.StudyActivityID, response.StudyActivityID)
		assert.Equal(t, expectedResponse.GroupID, response.GroupID)
	})

	t.Run("invalid activity ID", func(t *testing.T) {
		// Create request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/study_activities/invalid/launch", nil)
		c.AddParam("id", "invalid")

		// Act
		handler.LaunchStudyActivity(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Create request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/study_activities/1/launch", strings.NewReader(`{"invalid":"json"}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.AddParam("id", "1")

		// Act
		handler.LaunchStudyActivity(c)

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
