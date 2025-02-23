package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/services"
)

type StudySessionHandler struct {
	service services.StudySessionServiceInterface
}

// GetStudySessionWords godoc
// @Summary Get words for a study session
// @Description Returns a paginated list of words reviewed in a specific study session
// @Tags study_sessions
// @Accept json
// @Produce json
// @Param id path int true "Study Session ID"
// @Param limit query int false "Number of items per page" default(100)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} models.StudySessionWordsResponse
// @Router /api/study_sessions/{id}/words [get]
func (h *StudySessionHandler) GetStudySessionWords(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid study session ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid study session ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	words, err := h.service.GetStudySessionWords(id, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get study session words")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, words)
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

// ReviewWord godoc
// @Summary Review a word in a study session
// @Description Records whether a word was correctly or incorrectly reviewed in a study session
// @Tags study_sessions
// @Accept json
// @Produce json
// @Param id path int true "Study Session ID"
// @Param word_id path int true "Word ID"
// @Param request body models.WordReviewRequest true "Review request"
// @Success 200 {object} models.WordReviewResponse
// @Router /api/study_sessions/{id}/words/{word_id}/review [post]
func (h *StudySessionHandler) ReviewWord(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid study session ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid study session ID"})
		return
	}

	wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid word ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	var req models.WordReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	resp, err := h.service.ReviewWord(sessionID, wordID, req.Correct)
	if err != nil {
		log.Error().Err(err).Msg("Failed to review word")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
