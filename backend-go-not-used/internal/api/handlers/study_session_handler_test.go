package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/api/models"
	dbmodels "github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStudySessionService struct {
	mock.Mock
}

func (m *MockStudySessionService) CreateSession(ctx context.Context, session *dbmodels.StudySession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockStudySessionService) EndSession(id int64, endTime time.Time) error {
	args := m.Called(id, endTime)
	return args.Error(0)
}

func (m *MockStudySessionService) GetSession(id int64) (*dbmodels.StudySession, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dbmodels.StudySession), args.Error(1)
}

func (m *MockStudySessionService) GetSessionStats(ctx context.Context, groupID int64) (*dbmodels.StudyStats, error) {
	args := m.Called(ctx, groupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dbmodels.StudyStats), args.Error(1)
}

func (m *MockStudySessionService) ListSessions(page, limit int) ([]*dbmodels.StudySession, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]*dbmodels.StudySession), args.Get(1).(int), args.Error(2)
}

func (m *MockStudySessionService) ListGroupSessions(ctx context.Context, groupID int64, page, pageSize int) ([]*dbmodels.StudySession, error) {
	args := m.Called(ctx, groupID, page, pageSize)
	return args.Get(0).([]*dbmodels.StudySession), args.Error(1)
}

func TestCreateSession(t *testing.T) {
	mockService := new(MockStudySessionService)
	handler := NewStudySessionHandler(mockService)

	request := models.CreateStudySessionRequest{
		GroupID:         1,
		StudyActivityID: 2,
		TotalWords:      10,
		CorrectWords:    8,
		DurationSeconds: 300,
	}

	mockService.On("CreateSession", mock.Anything, mock.AnythingOfType("*dbmodels.StudySession")).Return(nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/sessions", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateSession(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.StudySession
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotZero(t, response.ID)
	assert.Equal(t, request.GroupID, response.GroupID)
	assert.Equal(t, request.StudyActivityID, response.StudyActivityID)

	mockService.AssertExpectations(t)
}

func TestGetGroupStats(t *testing.T) {
	mockService := new(MockStudySessionService)
	handler := NewStudySessionHandler(mockService)

	stats := models.SessionStats{
		TotalSessions: 5,
		TotalWords:    100,
		TotalCorrect:  80,
	}

	mockService.On("GetSessionStats", mock.Anything, int64(1)).Return(&stats, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "groupId", Value: "1"}}

	handler.GetGroupStats(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		TotalSessions   int64   `json:"total_sessions"`
		TotalWords      int     `json:"total_words"`
		CorrectWords    int     `json:"correct_words"`
		AverageAccuracy float64 `json:"average_accuracy"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), response.TotalSessions)
	assert.Equal(t, 100, response.TotalWords)
	assert.Equal(t, 80, response.CorrectWords)
	assert.Equal(t, float64(80), response.AverageAccuracy)

	mockService.AssertExpectations(t)
}

func TestListSessions(t *testing.T) {
	mockService := new(MockStudySessionService)
	handler := NewStudySessionHandler(mockService)

	sessions := []*dbmodels.StudySession{
		{
			ID:              1,
			GroupID:         1,
			TotalWords:      20,
			CorrectWords:    15,
			StartTime:       time.Now().Add(-1 * time.Hour),
			EndTime:         time.Now(),
			DurationSeconds: 3600,
		},
	}

	mockService.On("ListSessions", 1, 100).Return(sessions, 1, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/sessions?page=1&limit=100", nil)

	handler.ListSessions(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Items []models.StudySession `json:"items"`
		Total int                   `json:"total"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Items))
	assert.Equal(t, 20, response.Items[0].TotalWords)
	assert.Equal(t, 15, response.Items[0].CorrectWords)

	mockService.AssertExpectations(t)
}
