package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

type MockStudySessionService struct {
	mock.Mock
}

func (m *MockStudySessionService) GetStudySessionWords(sessionID int64, limit, offset int) (*models.StudySessionWordsResponse, error) {
	args := m.Called(sessionID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudySessionWordsResponse), args.Error(1)
}

func (m *MockStudySessionService) GetAllStudySessions(limit, offset int) (*models.StudySessionListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudySessionListResponse), args.Error(1)
}

func (m *MockStudySessionService) ReviewWord(sessionID, wordID int64, correct bool) (*models.WordReviewResponse, error) {
	args := m.Called(sessionID, wordID, correct)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WordReviewResponse), args.Error(1)
}

func TestStudySessionHandler_GetStudySessionWords(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful retrieval", func(t *testing.T) {
		// Setup
		mockService := new(MockStudySessionService)
		handler := NewStudySessionHandler(mockService)
		// Arrange
		expectedResponse := &models.StudySessionWordsResponse{
			Items: []*models.WordResponse{
				{
					ID:      1,
					Italian: "ciao",
					English: "hello",
					Parts:   map[string]interface{}{"part": "greeting"},
				},
			},
			Pagination: models.PaginationResponse{
				CurrentPage:  1,
				TotalPages:   1,
				TotalItems:   1,
				ItemsPerPage: 100,
			},
		}

		mockService.On("GetStudySessionWords", int64(1), 100, 0).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/study_sessions/1/words", nil)

		// Act
		handler.GetStudySessionWords(c)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.StudySessionWordsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.Items[0].ID, response.Items[0].ID)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid session ID", func(t *testing.T) {
		// Setup
		mockService := new(MockStudySessionService)
		handler := NewStudySessionHandler(mockService)
		
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/study_sessions/invalid/words", nil)

		handler.GetStudySessionWords(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		// Setup
		mockService := new(MockStudySessionService)
		handler := NewStudySessionHandler(mockService)
		mockService.On("GetStudySessionWords", int64(1), 100, 0).Return(nil, errors.New("service error")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/study_sessions/1/words", nil)

		handler.GetStudySessionWords(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response["error"])
		mockService.AssertExpectations(t)
	})
}

func TestStudySessionHandler_GetAllStudySessions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful retrieval", func(t *testing.T) {
		mockService := new(MockStudySessionService)
		handler := NewStudySessionHandler(mockService)
		expectedResponse := &models.StudySessionListResponse{
			Items: []models.StudySessionDetailResponse{
				{
					ID:           1,
					ActivityName: "Test Activity",
					GroupName:    "Test Group",
					Stats: models.StudySessionStats{
						TotalWords:   10,
						CorrectWords: 5,
					},
					ReviewItems: []models.WordReviewItem{},
				},
			},
			Pagination: models.PaginationResponse{
				CurrentPage:  1,
				TotalPages:   1,
				TotalItems:   1,
				ItemsPerPage: 100,
			},
		}

		mockService.On("GetAllStudySessions", 100, 0).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/study_sessions", nil)

		handler.GetAllStudySessions(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.StudySessionListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.Items[0].ID, response.Items[0].ID)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(MockStudySessionService)
		handler := NewStudySessionHandler(mockService)
		mockService.On("GetAllStudySessions", 100, 0).Return(nil, errors.New("service error")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/study_sessions", nil)

		handler.GetAllStudySessions(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response["error"])
		mockService.AssertExpectations(t)
	})
}

func TestStudySessionHandler_ReviewWord(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful review", func(t *testing.T) {
		// Setup
		mockService := new(MockStudySessionService)
		handler := NewStudySessionHandler(mockService)
		expectedResponse := &models.WordReviewResponse{
			Success: true,
			WordID:  1,
		}

		mockService.On("ReviewWord", int64(1), int64(1), true).Return(expectedResponse, nil)

		reqBody := models.WordReviewRequest{Correct: true}
		reqBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{Key: "id", Value: "1"},
			{Key: "word_id", Value: "1"},
		}
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/study_sessions/1/words/1/review", bytes.NewBuffer(reqBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ReviewWord(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.WordReviewResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.Success, response.Success)
		assert.Equal(t, expectedResponse.WordID, response.WordID)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid session ID", func(t *testing.T) {
		// Setup
		mockService := new(MockStudySessionService)
		handler := NewStudySessionHandler(mockService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{Key: "id", Value: "invalid"},
			{Key: "word_id", Value: "1"},
		}
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/study_sessions/invalid/words/1/review", nil)

		handler.ReviewWord(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid word ID", func(t *testing.T) {
		// Setup
		mockService := new(MockStudySessionService)
		handler := NewStudySessionHandler(mockService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{Key: "id", Value: "1"},
			{Key: "word_id", Value: "invalid"},
		}
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/study_sessions/1/words/invalid/review", nil)

		handler.ReviewWord(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Setup
		mockService := new(MockStudySessionService)
		handler := NewStudySessionHandler(mockService)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{Key: "id", Value: "1"},
			{Key: "word_id", Value: "1"},
		}
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/study_sessions/1/words/1/review", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ReviewWord(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		// Setup
		mockService := new(MockStudySessionService)
		handler := NewStudySessionHandler(mockService)
		mockService.On("ReviewWord", int64(1), int64(1), true).Return(nil, errors.New("service error")).Once()

		reqBody := models.WordReviewRequest{Correct: true}
		reqBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{Key: "id", Value: "1"},
			{Key: "word_id", Value: "1"},
		}
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/study_sessions/1/words/1/review", bytes.NewBuffer(reqBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ReviewWord(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response["error"])
		mockService.AssertExpectations(t)
	})
}
