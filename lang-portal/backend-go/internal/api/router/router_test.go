package router

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create new router
	r := gin.New()

	// Setup routes
	SetupRoutes(r)

	// Get registered routes
	routes := r.Routes()

	// Assert health check route exists
	found := false
	for _, route := range routes {
		if route.Path == "/health" && route.Method == "GET" {
			found = true
			break
		}
	}

	assert.True(t, found, "Health check route should be registered")
}
