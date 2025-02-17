package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/docs"
	"github.com/jeevanions/italian-learning/internal/api/router"
	"github.com/jeevanions/italian-learning/internal/config"
	"github.com/jeevanions/italian-learning/internal/db/repository"
	"github.com/jeevanions/italian-learning/internal/domain/services"
	"github.com/jeevanions/italian-learning/internal/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Italian Learning API
// @version         1.0
// @description     API for Italian language learning platform
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func main() {
	// Initialize logger
	logger.Setup()
	log := logger.GetLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize database
	db, err := config.NewDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Initialize repositories
	wordRepo := repository.NewSQLiteWordRepository(db)
	groupRepo := repository.NewSQLiteGroupRepository(db)
	sessionRepo := repository.NewSQLiteStudySessionRepository(db)
	reviewRepo := repository.NewSQLiteWordReviewRepository(db)
	activityRepo := repository.NewSQLiteStudyActivityRepository(db)

	// Initialize services
	wordService := services.NewWordService(wordRepo)
	groupService := services.NewGroupService(groupRepo)
	sessionService := services.NewStudySessionService(sessionRepo, groupRepo, activityRepo)
	reviewService := services.NewWordReviewService(reviewRepo, wordRepo, sessionRepo)
	activityService := services.NewStudyActivityService(activityRepo)

	// Initialize Gin router
	r := gin.Default()

	// Setup API v1 routes
	v1 := r.Group("/api/v1")
	router.SetupRoutes(v1, wordService, groupService, sessionService, reviewService, activityService)

	// Swagger documentation
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	log.Info().Str("address", cfg.ServerAddress).Msg("Starting server")
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
