package router

import (
	"net/http"

	"github.com/jeevanions/lang-portal/backend-go/internal/db/seeder"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jeevanions/lang-portal/backend-go/internal/api/handlers"
	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/services"
)

// Setup initializes the router with all routes and middleware
func Setup(db *repository.SQLiteRepository, seeder *seeder.Seeder) *gin.Engine {
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		AllowWildcard:    true,
		MaxAge:           12 * 3600,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Initialize services and handlers
	dashboardService := services.NewDashboardService(db)
	settingsService := services.NewSettingsService(db, seeder)
	studySessionService := services.NewStudySessionService(db)

	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	settingsHandler := handlers.NewSettingsHandler(settingsService)
	studySessionHandler := handlers.NewStudySessionHandler(studySessionService)

	studyActivityService := services.NewStudyActivityService(db)
	studyActivityHandler := handlers.NewStudyActivityHandler(studyActivityService)

	wordService := services.NewWordService(db)
	wordHandler := handlers.NewWordHandler(wordService)

	groupService := services.NewGroupService(db)
	groupHandler := handlers.NewGroupHandler(groupService)

	// API routes
	api := r.Group("/api")
	{
		// Dashboard routes
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/last_study_session", dashboardHandler.GetLastStudySession)
			dashboard.GET("/study_progress", dashboardHandler.GetStudyProgress)
			dashboard.GET("/quick-stats", dashboardHandler.GetQuickStats)
		}

		// Study Activity routes
		studyActivities := api.Group("/study_activities")
		{
			studyActivities.GET("", studyActivityHandler.GetStudyActivities)
			studyActivities.GET("/:id", studyActivityHandler.GetStudyActivity)
			studyActivities.GET("/:id/study_sessions", studyActivityHandler.GetStudyActivitySessions)
			studyActivities.POST("/:id/launch", studyActivityHandler.LaunchStudyActivity)
		}

		// Word routes
		words := api.Group("/words")
		{
			words.GET("", wordHandler.GetWords)
			words.GET("/:id", wordHandler.GetWordByID)
		}

		// Settings routes
		api.POST("/reset_history", settingsHandler.ResetHistory)
		api.POST("/full_reset", settingsHandler.FullReset)

		// Study Session routes
		api.GET("/study_sessions", studySessionHandler.GetAllStudySessions)

		// Group routes
		groups := api.Group("/groups")
		{
			groups.GET("", groupHandler.GetGroups)
			groups.GET("/:id", groupHandler.GetGroupByID)
			groups.GET("/:id/words", groupHandler.GetGroupWords)
			groups.GET("/:id/study_sessions", groupHandler.GetGroupStudySessions)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"), // The url pointing to API definition
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	return r
}
