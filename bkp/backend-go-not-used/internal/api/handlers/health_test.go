package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name:         "GET request returns 200",
			method:       "GET",
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status": "healthy",
			},
		},
		{
			name:         "POST request returns 405",
			method:       "POST",
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: map[string]interface{}{
				"error": "Method not allowed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Gin in test mode
			gin.SetMode(gin.TestMode)
			router := gin.New()

			// Setup route with all methods to handle 405
			router.GET("/health", HealthCheck)
			router.Handle("POST", "/health", func(c *gin.Context) {
				c.JSON(http.StatusMethodNotAllowed, gin.H{
					"error": "Method not allowed",
				})
			})

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, "/health", nil)

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Parse and assert response body
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err == nil {
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

// Test concurrent requests
func TestHealthCheckConcurrent(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", HealthCheck)

	// Number of concurrent requests
	concurrentRequests := 100
	done := make(chan bool)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/health", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < concurrentRequests; i++ {
		<-done
	}
}
