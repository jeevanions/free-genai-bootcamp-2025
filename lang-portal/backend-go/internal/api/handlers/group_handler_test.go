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

type MockGroupService struct {
	mock.Mock
}

func (m *MockGroupService) GetGroups(limit, offset int) (*models.GroupListResponse, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupListResponse), args.Error(1)
}

func (m *MockGroupService) GetGroupByID(id int64) (*models.GroupDetailResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupDetailResponse), args.Error(1)
}

func (m *MockGroupService) GetGroupWords(groupID int64, limit, offset int) (*models.GroupWordsResponse, error) {
	args := m.Called(groupID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupWordsResponse), args.Error(1)
}

func (m *MockGroupService) GetGroupStudySessions(groupID int64, limit, offset int) (*models.GroupStudySessionsResponse, error) {
	args := m.Called(groupID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupStudySessionsResponse), args.Error(1)
}

func (m *MockGroupService) CreateGroup(name string) (*models.GroupResponse, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GroupResponse), args.Error(1)
}

func TestGroupHandler_GetGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		limit      string
		offset     string
		mockSetup  func(*MockGroupService)
		wantStatus int
		wantBody   *models.GroupListResponse
	}{
		{
			name:   "successful retrieval",
			limit:  "10",
			offset: "0",
			mockSetup: func(m *MockGroupService) {
				m.On("GetGroups", 10, 0).Return(&models.GroupListResponse{
					Items: []models.GroupResponse{
						{ID: 1, Name: "Basic Words", WordCount: 10},
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
			wantBody: &models.GroupListResponse{
				Items: []models.GroupResponse{
					{ID: 1, Name: "Basic Words", WordCount: 10},
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
			mockSetup: func(m *MockGroupService) {
				m.On("GetGroups", 10, 0).Return(nil, errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockGroupService)
			tt.mockSetup(mockService)
			handler := NewGroupHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/?limit="+tt.limit+"&offset="+tt.offset, nil)

			handler.GetGroups(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != nil {
				var got models.GroupListResponse
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &got)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestGroupHandler_GetGroupByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		groupID    string
		mockSetup  func(*MockGroupService)
		wantStatus int
		wantBody   *models.GroupDetailResponse
	}{
		{
			name:    "successful retrieval",
			groupID: "1",
			mockSetup: func(m *MockGroupService) {
				m.On("GetGroupByID", int64(1)).Return(&models.GroupDetailResponse{
					ID:   1,
					Name: "Basic Words",
					Stats: models.GroupStats{
						TotalWordCount: 10,
					},
				}, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: &models.GroupDetailResponse{
				ID:   1,
				Name: "Basic Words",
				Stats: models.GroupStats{
					TotalWordCount: 10,
				},
			},
		},
		{
			name:    "group not found",
			groupID: "999",
			mockSetup: func(m *MockGroupService) {
				m.On("GetGroupByID", int64(999)).Return(nil, nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   &models.GroupDetailResponse{},
		},
		{
			name:    "invalid id",
			groupID: "invalid",
			mockSetup: func(m *MockGroupService) {
				// No mock setup needed as it won't reach the service
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockGroupService)
			tt.mockSetup(mockService)
			handler := NewGroupHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.groupID}}
			c.Request = httptest.NewRequest("GET", "/groups/"+tt.groupID, nil)

			handler.GetGroupByID(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != nil {
				var got models.GroupDetailResponse
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &got)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestGroupHandler_GetGroupWords(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		groupID    string
		limit      string
		offset     string
		mockSetup  func(*MockGroupService)
		wantStatus int
		wantBody   *models.GroupWordsResponse
	}{
		{
			name:    "successful retrieval",
			groupID: "1",
			limit:   "10",
			offset:  "0",
			mockSetup: func(m *MockGroupService) {
				m.On("GetGroupWords", int64(1), 10, 0).Return(&models.GroupWordsResponse{
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
			wantBody: &models.GroupWordsResponse{
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
			name:    "invalid group id",
			groupID: "invalid",
			limit:   "10",
			offset:  "0",
			mockSetup: func(m *MockGroupService) {
				// No mock setup needed as it won't reach the service
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "service error",
			groupID: "1",
			limit:   "10",
			offset:  "0",
			mockSetup: func(m *MockGroupService) {
				m.On("GetGroupWords", int64(1), 10, 0).Return(nil, errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockGroupService)
			tt.mockSetup(mockService)
			handler := NewGroupHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.groupID}}
			c.Request = httptest.NewRequest("GET", "/groups/"+tt.groupID+"/words?limit="+tt.limit+"&offset="+tt.offset, nil)

			handler.GetGroupWords(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != nil {
				var got models.GroupWordsResponse
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &got)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestGroupHandler_GetGroupStudySessions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		groupID    string
		limit      string
		offset     string
		mockSetup  func(*MockGroupService)
		wantStatus int
		wantBody   *models.GroupStudySessionsResponse
	}{
		{
			name:    "successful retrieval",
			groupID: "1",
			limit:   "10",
			offset:  "0",
			mockSetup: func(m *MockGroupService) {
				m.On("GetGroupStudySessions", int64(1), 10, 0).Return(&models.GroupStudySessionsResponse{
					Items: []models.StudySessionResponse{
						{
							ID:           1,
							ActivityName: "Vocabulary Review",
							GroupName:    "Basic Words",
							WordsCount:   10,
							CorrectCount: 8,
						},
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
			wantBody: &models.GroupStudySessionsResponse{
				Items: []models.StudySessionResponse{
					{
						ID:           1,
						ActivityName: "Vocabulary Review",
						GroupName:    "Basic Words",
						WordsCount:   10,
						CorrectCount: 8,
					},
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
			name:    "invalid group id",
			groupID: "invalid",
			limit:   "10",
			offset:  "0",
			mockSetup: func(m *MockGroupService) {
				// No mock setup needed as it won't reach the service
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "service error",
			groupID: "1",
			limit:   "10",
			offset:  "0",
			mockSetup: func(m *MockGroupService) {
				m.On("GetGroupStudySessions", int64(1), 10, 0).Return(nil, errors.New("service error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockGroupService)
			tt.mockSetup(mockService)
			handler := NewGroupHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.groupID}}
			c.Request = httptest.NewRequest("GET", "/groups/"+tt.groupID+"/study_sessions?limit="+tt.limit+"&offset="+tt.offset, nil)

			handler.GetGroupStudySessions(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != nil {
				var got models.GroupStudySessionsResponse
				err := json.Unmarshal(w.Body.Bytes(), &got)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBody, &got)
			}
			mockService.AssertExpectations(t)
		})
	}
}
