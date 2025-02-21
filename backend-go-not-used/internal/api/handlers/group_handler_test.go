package handlers

import (
	"bytes"
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

type MockGroupService struct {
	mock.Mock
}

func (m *MockGroupService) GetGroup(ctx context.Context, id int64) (*models.Group, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Group), args.Error(1)
}

func (m *MockGroupService) ListGroups(ctx context.Context, page, limit int) ([]*models.Group, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]*models.Group), args.Error(1)
}

func (m *MockGroupService) GetGroupWords(ctx context.Context, groupID int64, page, limit int) ([]models.Word, int, error) {
	args := m.Called(ctx, groupID, page, limit)
	return args.Get(0).([]models.Word), args.Int(1), args.Error(2)
}

func (m *MockGroupService) CreateGroup(ctx context.Context, group *models.Group) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func TestListGroups(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	groups := []*models.Group{
		{
			ID:              1,
			Name:            "Food and Dining",
			Description:     "Common food and restaurant vocabulary",
			DifficultyLevel: 2,
			Category:        "thematic",
		},
	}

	mockService.On("ListGroups", mock.Anything, 1, 100).Return(groups, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/groups?page=1&limit=100", nil)

	handler.ListGroups(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Items []*models.Group `json:"items"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Items))
	assert.Equal(t, "Food and Dining", response.Items[0].Name)

	mockService.AssertExpectations(t)
}

func TestCreateGroup(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	group := &models.Group{
		Name:            "Food and Dining",
		Description:     "Common food and restaurant vocabulary",
		DifficultyLevel: 2,
		Category:        "thematic",
	}

	mockService.On("CreateGroup", mock.Anything, mock.AnythingOfType("*models.Group")).Return(nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(group)
	c.Request = httptest.NewRequest("POST", "/groups", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateGroup(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	mockService.AssertExpectations(t)
}

func TestGetGroup(t *testing.T) {
	mockService := new(MockGroupService)
	handler := NewGroupHandler(mockService)

	group := &models.Group{
		ID:              1,
		Name:            "Food and Dining",
		Description:     "Common food and restaurant vocabulary",
		DifficultyLevel: 2,
		Category:        "thematic",
	}

	mockService.On("GetGroup", mock.Anything, int64(1)).Return(group, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("GET", "/groups/1", nil)

	handler.GetGroup(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Group
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, group.ID, response.ID)
	assert.Equal(t, group.Name, response.Name)

	mockService.AssertExpectations(t)
}
