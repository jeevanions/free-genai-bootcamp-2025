package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/yourusername/italian-learning/docs"
	"github.com/yourusername/italian-learning/internal/api/router"
	"github.com/yourusername/italian-learning/internal/config"
	"github.com/yourusername/italian-learning/internal/pkg/logger"
)

// @title Italian Learning API
// @version 1.0
// @description API for Italian language learning platform
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize logger
	logger.Setup()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	router.SetupRoutes(r)

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
