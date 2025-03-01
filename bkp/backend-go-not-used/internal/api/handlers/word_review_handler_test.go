package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/api/models"
	dbmodels "github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWordReviewService struct {
	mock.Mock
}

func (m *MockWordReviewService) RecordReview(ctx context.Context, userID int64, review *dbmodels.WordReviewItem) error {
	args := m.Called(ctx, userID, review)
	return args.Error(0)
}

func (m *MockWordReviewService) GetWordStats(ctx context.Context, wordID int64) (*dbmodels.WordStats, error) {
	args := m.Called(ctx, wordID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dbmodels.WordStats), args.Error(1)
}

func (m *MockWordReviewService) GetReviewHistory(ctx context.Context, userID, wordID int64, limit int) ([]dbmodels.WordReviewItem, error) {
	args := m.Called(ctx, userID, wordID, limit)
	return args.Get(0).([]dbmodels.WordReviewItem), args.Error(1)
}

func (m *MockWordReviewService) CreateReview(ctx context.Context, review *dbmodels.WordReviewItem) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

func (m *MockWordReviewService) ListSessionReviews(ctx context.Context, sessionID int64) ([]*dbmodels.WordReviewItem, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]*dbmodels.WordReviewItem), args.Error(1)
}

// TestRecordReview tests the RecordReview handler with various scenarios including:
// - Successful review recording
// - Invalid review data
// - Service layer errors
func TestRecordReview(t *testing.T) {
	tests := []struct {
		name             string
		review           *models.WordReviewItem
		setupMock        func(*MockWordReviewService)
		expectedCode     int
		expectedResponse map[string]interface{}
	}{
		{
			name: "successful review recording",
			review: &models.WordReviewItem{
				UserID:         1,
				WordID:         1,
				StudySessionID: 2,
				Correct:        true,
			},
			setupMock: func(m *MockWordReviewService) {
				m.On("RecordReview", mock.Anything, int64(1), mock.AnythingOfType("*models.WordReviewItem")).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"success": true,
			},
		},
		{
			name: "invalid review data",
			review: &models.WordReviewItem{
				UserID: 1,
				// Missing required fields
			},
			setupMock: func(m *MockWordReviewService) {
				// No mock setup needed for invalid data
			},
			expectedCode: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "Invalid review data",
			},
		},
		{
			name: "service error",
			review: &models.WordReviewItem{
				UserID:         1,
				WordID:         1,
				StudySessionID: 2,
				Correct:        true,
			},
			setupMock: func(m *MockWordReviewService) {
				m.On("RecordReview", mock.Anything, int64(1), mock.AnythingOfType("*models.WordReviewItem")).Return(errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "Failed to record review",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWordReviewService)
			handler := NewWordReviewHandler(mockService)

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonData, _ := json.Marshal(tt.review)
			c.Request = httptest.NewRequest("POST", "/reviews", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.RecordReview(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResponse, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetWordStats(t *testing.T) {
	mockService := new(MockWordReviewService)
	handler := NewWordReviewHandler(mockService)

	stats := &dbmodels.WordStats{
		UserID:          1,
		WordID:          1,
		TotalAttempts:   10,
		CorrectAttempts: 8,
		SuccessRate:     80.0,
		LastReviewedAt:  "2025-02-18T10:00:00Z",
		MasteryLevel:    4,
	}

	mockService.On("GetWordStats", mock.Anything, int64(1)).Return(stats, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("GET", "/reviews/words/1/stats", nil)

	handler.GetWordStats(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dbmodels.WordStats
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 10, response.TotalAttempts)
	assert.Equal(t, 8, response.CorrectAttempts)
	assert.Equal(t, 80.0, response.SuccessRate)

	mockService.AssertExpectations(t)
}

func TestGetReviewHistory(t *testing.T) {
	mockService := new(MockWordReviewService)
	handler := NewWordReviewHandler(mockService)

	history := []dbmodels.WordReviewItem{
		{
			ID:             1,
			UserID:         1,
			WordID:         1,
			StudySessionID: 1,
			Correct:        true,
		},
		{
			ID:             2,
			UserID:         1,
			WordID:         1,
			StudySessionID: 2,
			Correct:        false,
		},
	}

	mockService.On("GetReviewHistory", mock.Anything, int64(1), int64(1), 10).Return(history, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("GET", "/reviews/words/1/history?limit=10", nil)

	handler.GetReviewHistory(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		History []dbmodels.WordReviewItem `json:"history"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response.History))
	assert.True(t, response.History[0].Correct)
	assert.False(t, response.History[1].Correct)

	mockService.AssertExpectations(t)
}
