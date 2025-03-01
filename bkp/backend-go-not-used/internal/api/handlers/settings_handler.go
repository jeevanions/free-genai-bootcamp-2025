package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/domain/services"
)

type SettingsHandler struct {
	settingsService services.SettingsService
}

func NewSettingsHandler(ss services.SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: ss}
}

func (h *SettingsHandler) GetSettings(c *gin.Context) {
	// TODO: Get userID from authenticated session
	userID := int64(1) // Temporary hardcoded value

	settings, err := h.settingsService.GetSettings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	// TODO: Get userID from authenticated session
	userID := int64(1) // Temporary hardcoded value

	var settings models.Settings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings.UserID = userID
	if err := h.settingsService.UpdateSettings(c.Request.Context(), &settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *SettingsHandler) GetPreferences(c *gin.Context) {
	// TODO: Get userID from authenticated session
	userID := int64(1) // Temporary hardcoded value

	prefs, err := h.settingsService.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prefs)
}

func (h *SettingsHandler) UpdatePreferences(c *gin.Context) {
	// TODO: Get userID from authenticated session
	userID := int64(1) // Temporary hardcoded value

	var prefs models.Preferences
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prefs.UserID = userID
	if err := h.settingsService.UpdatePreferences(c.Request.Context(), &prefs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
