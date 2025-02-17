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
) {
	// Health check
	r.GET("/health", handlers.HealthCheck)

	// Word routes
	wordHandler := handlers.NewWordHandler(wordService)
	words := r.Group("/words")
	{
		words.GET("/", wordHandler.ListWords)
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
		activities.POST("/", activityHandler.CreateActivity)
		activities.GET("/:id", activityHandler.GetActivity)
		activities.GET("/", activityHandler.ListActivities)
	}
}
