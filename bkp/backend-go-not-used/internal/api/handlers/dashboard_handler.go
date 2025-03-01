package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/domain/services"
)

type DashboardHandler struct {
	dashboardService services.DashboardService
}

func NewDashboardHandler(ds services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: ds}
}

// GetLastStudySession godoc
// @Summary Get last study session
// @Description Get details of the user's most recent study session
// @Tags dashboard
// @Produce json
// @Success 200 {object} apimodels.LastStudySessionResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /dashboard/last_study_session [get]
func (h *DashboardHandler) GetLastStudySession(c *gin.Context) {
	// TODO: Get userID from authenticated session
	userID := int64(1) // Temporary hardcoded value

	session, err := h.dashboardService.GetLastStudySession(c.Request.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No study sessions found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *DashboardHandler) GetStudyProgress(c *gin.Context) {
	// TODO: Get userID from authenticated session
	userID := int64(1) // Temporary hardcoded value

	progress, err := h.dashboardService.GetStudyProgress(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, progress)
}

func (h *DashboardHandler) GetQuickStats(c *gin.Context) {
	// TODO: Get userID from authenticated session
	userID := int64(1) // Temporary hardcoded value

	stats, err := h.dashboardService.GetQuickStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *DashboardHandler) GetStreak(c *gin.Context) {
	// TODO: Get userID from authenticated session
	userID := int64(1) // Temporary hardcoded value

	streak, err := h.dashboardService.GetStreak(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, streak)
}

func (h *DashboardHandler) GetMasteryMetrics(c *gin.Context) {
	// TODO: Get userID from authenticated session
	userID := int64(1) // Temporary hardcoded value

	metrics, err := h.dashboardService.GetMasteryMetrics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}
