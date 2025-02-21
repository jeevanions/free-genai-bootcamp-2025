package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStudyActivityService struct {
	mock.Mock
}

func (m *MockStudyActivityService) ListActivities(ctx context.Context, page, limit int) ([]*models.StudyActivity, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]*models.StudyActivity), args.Error(1)
}

func (m *MockStudyActivityService) GetActivity(ctx context.Context, id int64) (*models.StudyActivity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StudyActivity), args.Error(1)
}

func (m *MockStudyActivityService) CreateActivity(ctx context.Context, activity *models.StudyActivity) error {
	args := m.Called(ctx, activity)
	return args.Error(0)
}

func (m *MockStudyActivityService) GetCategories(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStudyActivityService) GetRecommended(ctx context.Context, limit int) ([]*models.StudyActivity, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*models.StudyActivity), args.Error(1)
}

func TestListActivities(t *testing.T) {
	mockService := new(MockStudyActivityService)
	handler := NewStudyActivityHandler(mockService)

	activities := []*models.StudyActivity{
		{
			ID:              1,
			Name:            "Vocabulary Quiz",
			Type:            "quiz",
			RequiresAudio:   false,
			DifficultyLevel: 2,
			Instructions:    "Match the Italian words with their English translations",
		},
	}

	mockService.On("ListActivities", mock.Anything, 1, 100).Return(activities, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/activities?page=1&limit=100", nil)

	handler.ListActivities(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Items []*models.StudyActivity `json:"items"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Items))
	assert.Equal(t, "Vocabulary Quiz", response.Items[0].Name)
	assert.Equal(t, "quiz", response.Items[0].Type)

	mockService.AssertExpectations(t)
}

func TestGetActivity(t *testing.T) {
	mockService := new(MockStudyActivityService)
	handler := NewStudyActivityHandler(mockService)

	activity := &models.StudyActivity{
		ID:              1,
		Name:            "Pronunciation Practice",
		Type:            "practice",
		RequiresAudio:   true,
		DifficultyLevel: 3,
		Instructions:    "Listen to the audio and repeat the words",
	}

	mockService.On("GetActivity", mock.Anything, int64(1)).Return(activity, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("GET", "/activities/1", nil)

	handler.GetActivity(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.StudyActivity
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Pronunciation Practice", response.Name)
	assert.True(t, response.RequiresAudio)
	assert.Equal(t, 3, response.DifficultyLevel)

	mockService.AssertExpectations(t)
}
