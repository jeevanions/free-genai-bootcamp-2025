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

type MockSettingsService struct {
    mock.Mock
}

func (m *MockSettingsService) GetSettings(ctx context.Context, userID int64) (*models.Settings, error) {
    args := m.Called(ctx, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Settings), args.Error(1)
}

func (m *MockSettingsService) UpdateSettings(ctx context.Context, settings *models.Settings) error {
    args := m.Called(ctx, settings)
    return args.Error(0)
}

func (m *MockSettingsService) GetPreferences(ctx context.Context, userID int64) (*models.Preferences, error) {
    args := m.Called(ctx, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.Preferences), args.Error(1)
}

func (m *MockSettingsService) UpdatePreferences(ctx context.Context, preferences *models.Preferences) error {
    args := m.Called(ctx, preferences)
    return args.Error(0)
}

func TestGetSettings(t *testing.T) {
    mockService := new(MockSettingsService)
    handler := NewSettingsHandler(mockService)

    settings := &models.Settings{
        ID:       1,
        UserID:   1,
        Theme:    "light",
        Language: "en",
    }

    mockService.On("GetSettings", mock.Anything, int64(1)).Return(settings, nil)

    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    handler.GetSettings(c)

    assert.Equal(t, http.StatusOK, w.Code)

    var response models.Settings
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, settings.Theme, response.Theme)
    assert.Equal(t, settings.Language, response.Language)

    mockService.AssertExpectations(t)
}

func TestUpdateSettings(t *testing.T) {
    mockService := new(MockSettingsService)
    handler := NewSettingsHandler(mockService)

    settings := &models.Settings{
        ID:       1,
        UserID:   1,
        Theme:    "dark",
        Language: "it",
    }

    mockService.On("UpdateSettings", mock.Anything, mock.AnythingOfType("*models.Settings")).Return(nil)

    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)

    jsonData, _ := json.Marshal(settings)
    c.Request = httptest.NewRequest("PUT", "/settings", bytes.NewBuffer(jsonData))
    c.Request.Header.Set("Content-Type", "application/json")

    handler.UpdateSettings(c)

    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]bool
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.True(t, response["success"])

    mockService.AssertExpectations(t)
}
