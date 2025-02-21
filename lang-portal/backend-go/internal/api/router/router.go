package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jeevanions/lang-portal/backend-go/internal/api/handlers"
	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/services"
)

// Setup initializes the router with all routes and middleware
func Setup(db *repository.SQLiteRepository) *gin.Engine {
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
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	studyActivityService := services.NewStudyActivityService(db)
	studyActivityHandler := handlers.NewStudyActivityHandler(studyActivityService)

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
			studyActivities.GET("/:id", studyActivityHandler.GetStudyActivity)
			studyActivities.GET("/:id/study_sessions", studyActivityHandler.GetStudyActivitySessions)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"), // The url pointing to API definition
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	return r
}
