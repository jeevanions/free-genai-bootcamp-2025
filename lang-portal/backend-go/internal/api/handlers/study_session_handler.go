package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/services"
)

type StudySessionHandler struct {
	service services.StudySessionServiceInterface
}

func NewStudySessionHandler(service services.StudySessionServiceInterface) *StudySessionHandler {
	return &StudySessionHandler{service: service}
}

// GetAllStudySessions godoc
// @Summary Get all study sessions
// @Description Returns a paginated list of all study sessions with activity name, group name, and review items
// @Tags study_sessions
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page" default(100)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} models.StudySessionListResponse
// @Router /api/study_sessions [get]
func (h *StudySessionHandler) GetAllStudySessions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	sessions, err := h.service.GetAllStudySessions(limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get study sessions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}
