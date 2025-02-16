package integration

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yourusername/italian-learning/internal/api/router"
	"github.com/yourusername/italian-learning/internal/config"
)

func TestApplicationStartup(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load configuration with test values
	os.Setenv("SERVER_ADDRESS", ":8081") // Use different port for testing
	cfg, err := config.Load()
	assert.NoError(t, err)

	// Create a new gin engine
	r := gin.New()
	router.SetupRoutes(r)

	// Create server
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Run tests
	t.Run("Health Check", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8081/health")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// Shutdown server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	assert.NoError(t, err, "Server shutdown failed")
}

func TestGracefulShutdown(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load configuration
	os.Setenv("SERVER_ADDRESS", ":8082") // Use different port for second test
	cfg, err := config.Load()
	assert.NoError(t, err)

	// Create a new gin engine
	r := gin.New()
	router.SetupRoutes(r)

	// Create server
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Make a request
	resp, err := http.Get("http://localhost:8082/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Shutdown server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	assert.NoError(t, err, "Server shutdown failed")
}
