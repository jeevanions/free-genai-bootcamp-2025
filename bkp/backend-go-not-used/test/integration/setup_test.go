package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplicationStartup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Use test server instead of real server
	db := setupTestDB(t)
	r := setupTestServer(t, db)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Health Check", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestGracefulShutdown(t *testing.T) {
	// Setup test server
	gin.SetMode(gin.TestMode)
	logger.Setup()

	db := setupTestDB(t)
	r := setupTestServer(t, db)
	srv := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Make a test request
	_, err = http.Get("http://localhost:8082/health")
	require.NoError(t, err)

	// Initiate shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	assert.NoError(t, err, "Server shutdown failed")

	// Verify server has shut down
	_, err = http.Get("http://localhost:8082/health")
	assert.Error(t, err, "Expected error after shutdown")
}



}



}


