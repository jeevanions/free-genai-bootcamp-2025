package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	apimodels "github.com/jeevanions/italian-learning/internal/api/models"
	dbmodels "github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/domain/services"
)

type StudySessionHandler struct {
	service services.StudySessionService
}

func NewStudySessionHandler(service services.StudySessionService) *StudySessionHandler {
	return &StudySessionHandler{service: service}
}

// CreateSession godoc
// @Summary Create a new study session
// @Description Record a completed study session
// @Tags sessions
// @Accept json
// @Produce json
// @Param session body apimodels.CreateStudySessionRequest true "Study session details"
// @Success 201 {object} apimodels.StudySessionResponse
// @Failure 400 {object} map[string]string
// @Router /sessions [post]
func (h *StudySessionHandler) CreateSession(c *gin.Context) {
	var request apimodels.CreateStudySessionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbSession := &dbmodels.StudySession{
		GroupID:         request.GroupID,
		StudyActivityID: request.StudyActivityID,
		TotalWords:      request.TotalWords,
		CorrectWords:    request.CorrectWords,
		DurationSeconds: request.DurationSeconds,
	}

	if err := h.service.CreateSession(c.Request.Context(), dbSession); err != nil {
		if err == services.ErrInvalidSession {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	response := apimodels.StudySessionResponse{
		ID:              dbSession.ID,
		GroupID:         dbSession.GroupID,
		StudyActivityID: dbSession.StudyActivityID,
		TotalWords:      dbSession.TotalWords,
		CorrectWords:    dbSession.CorrectWords,
		DurationSeconds: dbSession.DurationSeconds,
		CreatedAt:       dbSession.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetGroupStats godoc
// @Summary Get study statistics for a group
// @Description Retrieve aggregated study statistics for a specific group
// @Tags sessions
// @Produce json
// @Param groupId path int true "Group ID"
// @Success 200 {object} apimodels.GroupStats
// @Failure 404 {object} map[string]string
// @Router /sessions/groups/{groupId}/stats [get]
func (h *StudySessionHandler) GetGroupStats(c *gin.Context) {
	// Add proper error handling and database query
	groupID, _ := strconv.ParseInt(c.Param("groupId"), 10, 64)

	stats, err := h.service.GetSessionStats(c.Request.Context(), groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	groupStats := apimodels.GroupStats{
		TotalSessions:   stats.TotalSessions,
		TotalWords:      int(stats.TotalWords),
		CorrectWords:    int(stats.TotalCorrect),
		AverageAccuracy: float64(stats.TotalCorrect) / float64(stats.TotalWords) * 100,
	}

	c.JSON(http.StatusOK, groupStats)
}

// ListGroupSessions godoc
// @Summary List study sessions for a group
// @Description Get a paginated list of study sessions for a specific group
// @Tags sessions
// @Produce json
// @Param groupId path int true "Group ID"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {array} apimodels.StudySessionResponse
// @Router /sessions/groups/{groupId} [get]
func (h *StudySessionHandler) ListGroupSessions(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("groupId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	// Validate and set pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	sessions, err := h.service.ListGroupSessions(c.Request.Context(), groupID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sessions"})
		return
	}

	// Convert to response format
	responses := make([]apimodels.StudySessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = apimodels.StudySessionResponse{
			ID:              session.ID,
			GroupID:         session.GroupID,
			StudyActivityID: session.StudyActivityID,
			TotalWords:      session.TotalWords,
			CorrectWords:    session.CorrectWords,
			DurationSeconds: session.DurationSeconds,
			StartTime:       session.StartTime,
			EndTime:         session.EndTime,
			CreatedAt:       session.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, responses)
}

func (h *StudySessionHandler) ListSessions(c *gin.Context) {
	page, limit := parsePaginationParams(c)

	sessions, total, err := h.service.ListSessions(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sessions"})
		return
	}

	response := make([]apimodels.StudySessionResponse, len(sessions))
	for i, s := range sessions {
		response[i] = apimodels.StudySessionResponse{
			ID:              s.ID,
			GroupID:         s.GroupID,
			StudyActivityID: s.StudyActivityID,
			TotalWords:      s.TotalWords,
			CorrectWords:    s.CorrectWords,
			DurationSeconds: s.DurationSeconds,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items": response,
		"total": total,
	})
}
