package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAnalyticsService struct {
	mock.Mock
}

func (m *MockAnalyticsService) GetSessionAnalytics(ctx context.Context, userID int64) (*models.SessionAnalytics, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SessionAnalytics), args.Error(1)
}

func (m *MockAnalyticsService) GetSessionCalendar(ctx context.Context, userID int64) (*models.SessionCalendar, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SessionCalendar), args.Error(1)
}

func TestGetSessionAnalytics(t *testing.T) {
	mockService := new(MockAnalyticsService)
	handler := NewAnalyticsHandler(mockService)

	analytics := &models.SessionAnalytics{
		TotalSessions:     10,
		TotalStudyTime:    2 * time.Hour,
		AverageSessionLen: 12 * time.Minute,
		WordsLearned:      50,
		WordsReviewed:     100,
	}

	mockService.On("GetSessionAnalytics", mock.Anything, int64(1)).Return(analytics, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler.GetSessionAnalytics(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SessionAnalytics
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, analytics.TotalSessions, response.TotalSessions)
	assert.Equal(t, analytics.WordsLearned, response.WordsLearned)

	mockService.AssertExpectations(t)
}

func TestGetSessionCalendar(t *testing.T) {
	mockService := new(MockAnalyticsService)
	handler := NewAnalyticsHandler(mockService)

	testDate := time.Date(2025, 2, 18, 0, 0, 0, 0, time.UTC)
	calendar := &models.SessionCalendar{
		Sessions: []models.CalendarSession{
			{
				Date:      testDate,
				Duration:  30 * time.Minute,
				WordCount: 25,
			},
		},
	}

	mockService.On("GetSessionCalendar", mock.Anything, int64(1)).Return(calendar, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/analytics/calendar?start=2025-02-01&end=2025-02-28", nil)

	handler.GetSessionCalendar(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.SessionCalendar
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Sessions))

	mockService.AssertExpectations(t)
    assert.NoError(t, err)
    assert.Equal(t, 1, len(response.Sessions))
    assert.Equal(t, calendar.Sessions[0].WordCount, response.Sessions[0].WordCount)

    mockService.AssertExpectations(t)
}
