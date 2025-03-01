package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/api/models"
	dbmodels "github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWordService struct {
	mock.Mock
}

func (m *MockWordService) GetWord(ctx context.Context, id int64) (*models.Word, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Word), args.Error(1)
}

func (m *MockWordService) ListWords(ctx context.Context, page, limit int) ([]models.Word, int, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]models.Word), args.Int(1), args.Error(2)
}

func (m *MockWordService) CreateWord(ctx context.Context, req *models.CreateWordRequest) (*models.Word, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Word), args.Error(1)
}

func (m *MockWordService) SearchWords(ctx context.Context, query string, page, pageSize int) ([]models.Word, int, error) {
	args := m.Called(ctx, query, page, pageSize)
	return args.Get(0).([]models.Word), args.Int(1), args.Error(2)
}

func (m *MockWordService) GetFilters(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestListWords(t *testing.T) {
	mockService := new(MockWordService)
	handler := NewWordHandler(mockService)

	words := []dbmodels.Word{
		{
			ID:              1,
			Italian:         "ciao",
			English:         "hello",
			PartsOfSpeech:   "interjection",
			DifficultyLevel: 1,
		},
	}

	mockService.On("ListWords", mock.Anything, 1, 100).Return(words, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/words?page=1&limit=100", nil)

	handler.ListWords(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Items []dbmodels.Word `json:"items"`
		Total int             `json:"total"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Items))
	assert.Equal(t, "ciao", response.Items[0].Italian)

	mockService.AssertExpectations(t)
}

func TestCreateWord(t *testing.T) {
	mockService := new(MockWordService)
	handler := NewWordHandler(mockService)

	request := &models.CreateWordRequest{
		Italian:         "ciao",
		English:         "hello",
		PartsOfSpeech:   "interjection",
		DifficultyLevel: 1,
	}

	mockService.On("CreateWord", mock.Anything, mock.AnythingOfType("*models.CreateWordRequest")).Return(&models.Word{
		ID:              1,
		Italian:         request.Italian,
		English:         request.English,
		PartsOfSpeech:   request.PartsOfSpeech,
		DifficultyLevel: request.DifficultyLevel,
	}, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request = httptest.NewRequest("POST", "/words", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateWord(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	mockService.AssertExpectations(t)
}

func TestSearchWords(t *testing.T) {
	mockService := new(MockWordService)
	handler := NewWordHandler(mockService)

	words := []models.Word{
		{
			ID:              1,
			Italian:         "uomo",
			English:         "man",
			PartsOfSpeech:   "noun",
			DifficultyLevel: 1,
		},
	}

	mockService.On("SearchWords", mock.Anything, "man", 1, 10).Return(words, 1, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/words/search?q=man", nil)
	c.Request.URL.RawQuery = "q=man"

	handler.SearchWords(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.Total)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.Limit)

	items, ok := response.Items.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1, len(items))

	mockService.AssertExpectations(t)

	// Test empty query
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/words/search", nil)

	handler.SearchWords(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test too long query
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/words/search?q="+strings.Repeat("a", 101), nil)
	c.Request.URL.RawQuery = "q=" + strings.Repeat("a", 101)

	handler.SearchWords(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetFilters(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockService := new(MockWordService)
		handler := NewWordHandler(mockService)

		filters := map[string]interface{}{
			"parts_of_speech":   []string{"noun", "verb", "adjective"},
			"difficulty_levels": []int{1, 2, 3, 4, 5},
			"genders":          []string{"masculine", "feminine", "neuter"},
			"numbers":          []string{"singular", "plural"},
		}

		mockService.On("GetFilters", mock.Anything).Return(filters, nil)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/words/filters", nil)

		handler.GetFilters(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check all filter keys exist
		assert.Contains(t, response, "parts_of_speech")
		assert.Contains(t, response, "difficulty_levels")
		assert.Contains(t, response, "genders")
		assert.Contains(t, response, "numbers")

		// Check filter values
		pos, ok := response["parts_of_speech"].([]interface{})
		assert.True(t, ok)
		assert.Contains(t, pos, "noun")

		levels, ok := response["difficulty_levels"].([]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(1), levels[0])

		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(MockWordService)
		handler := NewWordHandler(mockService)

		mockService.On("GetFilters", mock.Anything).Return(nil, assert.AnError)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/words/filters", nil)

		handler.GetFilters(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, assert.AnError.Error(), response.Error)

		mockService.AssertExpectations(t)
	})
}
