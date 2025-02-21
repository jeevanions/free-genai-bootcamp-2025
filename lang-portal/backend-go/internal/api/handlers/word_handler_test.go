package handlers

import (
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

type MockWordService struct {
	mock.Mock
}

func (m *MockWordService) GetWords(limit, offset int) (*models.WordListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WordListResponse), args.Error(1)
}

func (m *MockWordService) GetWordByID(id int64) (*models.WordResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WordResponse), args.Error(1)
}

func TestWordHandler_GetWords(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		limit      string
		offset     string
		mockSetup  func(*MockWordService)
		wantStatus int
		wantBody   *models.WordListResponse
	}{
		{
			name:   "successful retrieval",
			limit:  "10",
			offset: "0",
			mockSetup: func(m *MockWordService) {
				m.On("GetWords", 10, 0).Return(&models.WordListResponse{
					Items: []models.WordResponse{
						{ID: 1, Italian: "ciao", English: "hello"},
					},
					Pagination: models.PaginationResponse{
						CurrentPage:  1,
						TotalPages:   1,
						TotalItems:   1,
						ItemsPerPage: 10,
					},
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: &models.WordListResponse{
				Items: []models.WordResponse{
					{ID: 1, Italian: "ciao", English: "hello"},
				},
				Pagination: models.PaginationResponse{
					CurrentPage:  1,
					TotalPages:   1,
					TotalItems:   1,
					ItemsPerPage: 10,
				},
			},
		},
		{
			name:   "service error",
			limit:  "10",
			offset: "0",
			mockSetup: func(m *MockWordService) {
				m.On("GetWords", 10, 0).Return(nil, errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWordService)
			tt.mockSetup(mockService)
			handler := NewWordHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?limit="+tt.limit+"&offset="+tt.offset, nil)

			handler.GetWords(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != nil {
				var got models.WordListResponse
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &got)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestWordHandler_GetWordByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		wordID     string
		mockSetup  func(*MockWordService)
		wantStatus int
		wantBody   *models.WordResponse
	}{
		{
			name:   "successful retrieval",
			wordID: "1",
			mockSetup: func(m *MockWordService) {
				m.On("GetWordByID", int64(1)).Return(&models.WordResponse{
					ID:       1,
					Italian:  "ciao",
					English: "hello",
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: &models.WordResponse{
				ID:       1,
				Italian:  "ciao",
				English: "hello",
			},
		},
		{
			name:   "word not found",
			wordID: "999",
			mockSetup: func(m *MockWordService) {
				m.On("GetWordByID", int64(999)).Return(nil, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: &models.WordResponse{},
		},
		{
			name:   "invalid id",
			wordID: "invalid",
			mockSetup: func(m *MockWordService) {
				// No mock setup needed as it won't reach the service
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockWordService)
			tt.mockSetup(mockService)
			handler := NewWordHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.wordID}}
			c.Request = httptest.NewRequest("GET", "/words/"+tt.wordID, nil)

			handler.GetWordByID(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != nil {
				var got models.WordResponse
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &got)
			}
			mockService.AssertExpectations(t)
		})
	}
}
