package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/services"
)

type StudyActivityHandler struct {
	service services.StudyActivityServiceInterface
}

func NewStudyActivityHandler(service services.StudyActivityServiceInterface) *StudyActivityHandler {
	return &StudyActivityHandler{service: service}
}

// GetStudyActivity godoc
// @Summary Get study activity details
// @Description Returns details about a specific study activity
// @Tags study_activities
// @Accept json
// @Produce json
// @Param id path int true "Study Activity ID"
// @Success 200 {object} models.StudyActivityResponse
// @Router /api/study_activities/{id} [get]
func (h *StudyActivityHandler) GetStudyActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid study activity ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid study activity ID"})
		return
	}

	activity, err := h.service.GetStudyActivity(id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get study activity")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if activity == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study activity not found"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

// GetStudyActivitySessions godoc
// @Summary Get study sessions for an activity
// @Description Returns a list of study sessions for a specific activity
// @Tags study_activities
// @Accept json
// @Produce json
// @Param id path int true "Study Activity ID"
// @Success 200 {object} models.StudySessionsListResponse
// @Router /api/study_activities/{id}/study_sessions [get]
func (h *StudyActivityHandler) GetStudyActivitySessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid study activity ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid study activity ID"})
		return
	}

	sessions, err := h.service.GetStudyActivitySessions(id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get study activity sessions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}
