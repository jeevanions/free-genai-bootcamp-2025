package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/services"
)

type StudyActivityHandler struct {
	service services.StudyActivityServiceInterface
}

func NewStudyActivityHandler(service services.StudyActivityServiceInterface) *StudyActivityHandler {
	return &StudyActivityHandler{service: service}
}

// GetStudyActivities godoc
// @Summary Get all study activities
// @Description Returns a list of available study activities
// @Tags study_activities
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page" default(100)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} models.StudyActivityListResponse
// @Router /api/study_activities [get]
func (h *StudyActivityHandler) GetStudyActivities(c *gin.Context) {
	limit := 100 // Default limit as per spec
	offset := 0  // Default offset

	// Parse pagination parameters if provided
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	activities, err := h.service.GetStudyActivities(limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get study activities")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if activities == nil {
		c.JSON(http.StatusOK, models.StudyActivityListResponse{
			Items: []models.StudyActivityResponse{},
			Pagination: models.PaginationResponse{
				CurrentPage:  1,
				TotalPages:   1,
				TotalItems:   0,
				ItemsPerPage: limit,
			},
		})
		return
	}

	c.JSON(http.StatusOK, activities)
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
// LaunchStudyActivity godoc
// @Summary Launch a new study activity session
// @Description Launches a new study activity session for a specific group
// @Tags study_activities
// @Accept json
// @Produce json
// @Param id path int true "Study Activity ID"
// @Param request body models.LaunchStudyActivityRequest true "Launch request"
// @Success 200 {object} models.LaunchStudyActivityResponse
// @Router /api/study_activities/{id}/launch [post]
func (h *StudyActivityHandler) LaunchStudyActivity(c *gin.Context) {
	// Parse activity ID
	activityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	// Parse request body
	var request models.LaunchStudyActivityRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Launch study activity
	response, err := h.service.LaunchStudyActivity(activityID, request.GroupID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to launch study activity")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to launch study activity"})
		return
	}

	c.JSON(http.StatusOK, response)
}

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
