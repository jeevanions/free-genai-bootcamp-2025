package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/api/models"
	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetLastStudySession(ctx context.Context, userID int64) (*apimodels.LastStudySessionResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*apimodels.LastStudySessionResponse), args.Error(1)
}

func (m *MockDashboardService) GetStudyProgress(ctx context.Context, userID int64) (*models.StudyProgress, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyProgress), args.Error(1)
}

func (m *MockDashboardService) GetQuickStats(ctx context.Context, userID int64) (*models.QuickStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.QuickStats), args.Error(1)
}

func (m *MockDashboardService) GetStreak(ctx context.Context, userID int64) (*models.Streak, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Streak), args.Error(1)
}

func (m *MockDashboardService) GetMasteryMetrics(ctx context.Context, userID int64) (*models.MasteryMetrics, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MasteryMetrics), args.Error(1)
}

func TestGetLastStudySession(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *MockDashboardService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			setupMock: func(m *MockDashboardService) {
				response := &apimodels.LastStudySessionResponse{
					ID:              1,
					GroupID:         1,
					GroupName:       "Basic Italian",
					ActivityID:      1,
					ActivityName:    "Vocabulary Quiz",
					TotalWords:      10,
					CorrectWords:    8,
					DurationSeconds: 300,
					StartTime:       time.Now(),
					EndTime:         time.Now().Add(5 * time.Minute),
				}
				m.On("GetLastStudySession", mock.Anything, int64(1)).Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"group_id":1,"group_name":"Basic Italian","activity_id":1,"activity_name":"Vocabulary Quiz","total_words":10,"correct_words":8,"duration_seconds":300}`,
		},
		{
			name: "no sessions found",
			setupMock: func(m *MockDashboardService) {
				m.On("GetLastStudySession", mock.Anything, int64(1)).Return(nil, sql.ErrNoRows)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"No study sessions found"}`,
		},
		{
			name: "internal error",
			setupMock: func(m *MockDashboardService) {
				m.On("GetLastStudySession", mock.Anything, int64(1)).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"database error"}`,
		},
	}

	for _, tt := range tests {
		t := tt
		t.Run(t.Name(), func(t *testing.T) {
			// Setup
			mockService := new(MockDashboardService)
			tt.setupMock(mockService)
			handler := NewDashboardHandler(mockService)

			// Create request
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Test
			handler.GetLastStudySession(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetStudyProgress(t *testing.T) {
	mockService := new(MockDashboardService)
	handler := NewDashboardHandler(mockService)

	progress := &models.StudyProgress{
		WordsLearned:   150,
		WordsToReview:  50,
		CompletionRate: 35.0,
	}

	mockService.On("GetStudyProgress", mock.Anything, int64(1)).Return(progress, nil)

    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    handler.GetStudyProgress(c)

    assert.Equal(t, http.StatusOK, w.Code)

    var response models.StudyProgress
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, progress.WordsLearned, response.WordsLearned)
    assert.Equal(t, progress.CompletionRate, response.CompletionRate)

    mockService.AssertExpectations(t)
}
