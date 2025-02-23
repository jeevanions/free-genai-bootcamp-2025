package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/services"
)

type DashboardHandler struct {
	service services.DashboardServiceInterface
}

func NewDashboardHandler(service services.DashboardServiceInterface) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetLastStudySession godoc
// @Summary Get last study session
// @Description Returns information about the most recent study session
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} models.DashboardLastStudySession
// @Router /api/dashboard/last_study_session [get]
func (h *DashboardHandler) GetLastStudySession(c *gin.Context) {
	session, err := h.service.GetLastStudySession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No study sessions found"})
		return
	}
	c.JSON(http.StatusOK, session)
}

// GetStudyProgress godoc
// @Summary Get study progress
// @Description Returns study progress statistics
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} models.DashboardStudyProgress
// @Router /api/dashboard/study_progress [get]
func (h *DashboardHandler) GetStudyProgress(c *gin.Context) {
	progress, err := h.service.GetStudyProgress()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get study progress")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetQuickStats godoc
// @Summary Get quick stats
// @Description Returns quick overview statistics
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} models.DashboardQuickStats
// @Router /api/dashboard/quick-stats [get]
func (h *DashboardHandler) GetQuickStats(c *gin.Context) {
	stats, err := h.service.GetQuickStats()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get quick stats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
