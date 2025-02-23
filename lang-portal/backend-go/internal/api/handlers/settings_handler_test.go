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
)

type MockSettingsService struct {
	mock.Mock
}

func (m *MockSettingsService) ResetHistory() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSettingsService) FullReset() error {
	args := m.Called()
	return args.Error(0)
}

func TestSettingsHandler_ResetHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful reset", func(t *testing.T) {
		// Setup
		mockService := new(MockSettingsService)
		handler := NewSettingsHandler(mockService)
		mockService.On("ResetHistory").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/reset_history", nil)

		handler.ResetHistory(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "History reset successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		// Setup
		mockService := new(MockSettingsService)
		handler := NewSettingsHandler(mockService)
		mockService.On("ResetHistory").Return(errors.New("service error")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/reset_history", nil)

		handler.ResetHistory(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response["error"])
		mockService.AssertExpectations(t)
	})
}

func TestSettingsHandler_FullReset(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful reset", func(t *testing.T) {
		// Setup
		mockService := new(MockSettingsService)
		handler := NewSettingsHandler(mockService)
		mockService.On("FullReset").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/full_reset", nil)

		handler.FullReset(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Full reset completed successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		// Setup
		mockService := new(MockSettingsService)
		handler := NewSettingsHandler(mockService)
		mockService.On("FullReset").Return(errors.New("service error")).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/full_reset", nil)

		handler.FullReset(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response["error"])
		mockService.AssertExpectations(t)
	})
}
