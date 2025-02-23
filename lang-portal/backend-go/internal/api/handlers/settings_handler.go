package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/services"
)

type SettingsHandler struct {
	service services.SettingsServiceInterface
}

func NewSettingsHandler(service services.SettingsServiceInterface) *SettingsHandler {
	return &SettingsHandler{service: service}
}

// ResetHistory godoc
// @Summary Reset study history
// @Description Deletes all study sessions and word review items
// @Tags settings
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/reset_history [post]
func (h *SettingsHandler) ResetHistory(c *gin.Context) {
	if err := h.service.ResetHistory(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.Status(http.StatusOK)
}

// FullReset godoc
// @Summary Full system reset
// @Description Drops all tables and recreates them with seed data
// @Tags settings
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/full_reset [post]
func (h *SettingsHandler) FullReset(c *gin.Context) {
	if err := h.service.FullReset(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.Status(http.StatusOK)
}
