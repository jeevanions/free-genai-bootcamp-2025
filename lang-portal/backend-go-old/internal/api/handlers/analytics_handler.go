package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/domain/services"
)

type AnalyticsHandler struct {
	analyticsService services.AnalyticsService
}

func NewAnalyticsHandler(as services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: as}
}

func (h *AnalyticsHandler) GetSessionAnalytics(c *gin.Context) {
	// TODO: Get userID from authenticated session
	userID := int64(1) // Temporary hardcoded value

	analytics, err := h.analyticsService.GetSessionAnalytics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, analytics)
}

func (h *AnalyticsHandler) GetSessionCalendar(c *gin.Context) {
	// TODO: Get userID from authenticated session
	userID := int64(1) // Temporary hardcoded value

	calendar, err := h.analyticsService.GetSessionCalendar(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, calendar)
}
