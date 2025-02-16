package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/italian-learning/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine) {
	// Health check endpoint
	r.GET("/health", handlers.HealthCheck)
}
