package router

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/db/repository"
	"github.com/jeevanions/italian-learning/internal/domain/services"
	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create new router
	r := gin.New()
	v1 := r.Group("/api/v1") // Create router group

	// Create mock services
	wordService := services.NewWordService(&mockWordRepo{})
	groupService := services.NewGroupService(&mockGroupRepo{})
	sessionService := services.NewStudySessionService(&mockSessionRepo{}, &mockGroupRepo{}, &mockActivityRepo{})
	reviewService := services.NewWordReviewService(&mockReviewRepo{}, &mockWordRepo{}, &mockSessionRepo{})
	activityService := services.NewStudyActivityService(&mockActivityRepo{})
	dashboardService := services.NewDashboardService(&mockDashboardRepo{}, &mockSessionRepo{}, &mockReviewRepo{})
	settingsService := services.NewSettingsService(&mockSettingsRepo{})
	analyticsService := services.NewAnalyticsService(&mockAnalyticsRepo{})

	// Setup routes
	SetupRoutes(v1, wordService, groupService, sessionService, reviewService, activityService, dashboardService, settingsService, analyticsService)

	// Get registered routes
	routes := r.Routes()

	// Assert health check route exists
	found := false
	for _, route := range routes {
		if route.Path == "/api/v1/health" && route.Method == "GET" {
			found = true
			break
		}
	}

	assert.True(t, found, "Health check route should be registered")
}

// Mock repositories for testing
type mockWordRepo struct {
	repository.WordRepository
}

type mockGroupRepo struct {
	repository.GroupRepository
}

type mockSessionRepo struct {
	repository.StudySessionRepository
}

type mockReviewRepo struct {
	repository.WordReviewRepository
}

type mockActivityRepo struct {
	repository.StudyActivityRepository
}

type mockDashboardRepo struct {
	repository.DashboardRepository
}

type mockSettingsRepo struct {
	repository.SettingsRepository
}

type mockAnalyticsRepo struct {
	repository.AnalyticsRepository
}
