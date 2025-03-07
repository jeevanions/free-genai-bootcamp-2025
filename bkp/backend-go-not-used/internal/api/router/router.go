package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/api/handlers"
	"github.com/jeevanions/italian-learning/internal/domain/services"
)

func SetupRoutes(
	r *gin.RouterGroup,
	wordService services.WordService,
	groupService services.GroupService,
	sessionService services.StudySessionService,
	reviewService services.WordReviewService,
	activityService services.StudyActivityService,
	dashboardService services.DashboardService,
	settingsService services.SettingsService,
	analyticsService services.AnalyticsService,
) {
	// Health check
	r.GET("/health", handlers.HealthCheck)

	// Dashboard routes
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	dashboard := r.Group("/dashboard")
	{
		dashboard.GET("/last_study_session", dashboardHandler.GetLastStudySession)
		dashboard.GET("/study_progress", dashboardHandler.GetStudyProgress)
		dashboard.GET("/quick-stats", dashboardHandler.GetQuickStats)
		dashboard.GET("/streak", dashboardHandler.GetStreak)
		dashboard.GET("/mastery", dashboardHandler.GetMasteryMetrics)
	}

	// Word routes
	wordHandler := handlers.NewWordHandler(wordService)
	words := r.Group("/words")
	{
		words.GET("/", wordHandler.ListWords)
		words.GET("/search", wordHandler.SearchWords)
		words.GET("/filters", wordHandler.GetFilters)
		words.POST("/", wordHandler.CreateWord)
		words.GET("/:id", wordHandler.GetWord)
	}

	// Group endpoints
	groupHandler := handlers.NewGroupHandler(groupService)
	groups := r.Group("/groups")
	{
		groups.POST("/", groupHandler.CreateGroup)
		groups.GET("/:id", groupHandler.GetGroup)
		groups.GET("/", groupHandler.ListGroups)
		groups.GET("/:id/progress", groupHandler.GetGroupProgress)
		groups.POST("/:id/words", groupHandler.AddWordsToGroup)
	}

	// Study session endpoints
	sessionHandler := handlers.NewStudySessionHandler(sessionService)
	sessions := r.Group("/sessions")
	{
		sessions.POST("/", sessionHandler.CreateSession)
		sessions.GET("/groups/:groupId/stats/", sessionHandler.GetGroupStats)
		sessions.GET("/groups/:groupId/", sessionHandler.ListGroupSessions)
	}

	// Word review endpoints
	reviewHandler := handlers.NewWordReviewHandler(reviewService)
	reviews := r.Group("/reviews")
	{
		reviews.POST("/", reviewHandler.CreateReview)
		reviews.GET("/sessions/:sessionId/", reviewHandler.ListSessionReviews)
		reviews.GET("/words/:wordId/stats/", reviewHandler.GetWordStats)
	}

	// Study activity endpoints
	activityHandler := handlers.NewStudyActivityHandler(activityService)
	activities := r.Group("/activities")
	{
		activities.GET("/", activityHandler.ListActivities)
		activities.GET("/categories", activityHandler.GetCategories)
		activities.GET("/recommended", activityHandler.GetRecommended)
		activities.POST("/", activityHandler.CreateActivity)
		activities.GET("/:id", activityHandler.GetActivity)
	}

	// Settings endpoints
	settingsHandler := handlers.NewSettingsHandler(settingsService)
	settings := r.Group("/settings")
	{
		settings.GET("/", settingsHandler.GetSettings)
		settings.PUT("/", settingsHandler.UpdateSettings)
		settings.GET("/preferences", settingsHandler.GetPreferences)
		settings.PUT("/preferences", settingsHandler.UpdatePreferences)
	}

	// Analytics endpoints
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	analytics := r.Group("/sessions")
	{
		analytics.GET("/analytics", analyticsHandler.GetSessionAnalytics)
		analytics.GET("/calendar", analyticsHandler.GetSessionCalendar)
	}
}
