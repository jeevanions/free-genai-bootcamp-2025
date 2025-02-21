package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/services"
)

type GroupHandler struct {
	service services.GroupServiceInterface
}

func NewGroupHandler(service services.GroupServiceInterface) *GroupHandler {
	return &GroupHandler{service: service}
}

// GetGroups godoc
// @Summary Get all groups
// @Description Returns a paginated list of groups
// @Tags groups
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page" default(100)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} models.GroupListResponse
// @Router /api/groups [get]
func (h *GroupHandler) GetGroups(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	groups, err := h.service.GetGroups(limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get groups")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetGroupByID godoc
// @Summary Get group by ID
// @Description Returns details about a specific group
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} models.GroupDetailResponse
// @Router /api/groups/{id} [get]
func (h *GroupHandler) GetGroupByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid group ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	group, err := h.service.GetGroupByID(id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get group")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if group == nil {
		c.JSON(http.StatusOK, models.GroupDetailResponse{})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetGroupWords godoc
// @Summary Get words in a group
// @Description Returns a paginated list of words in a specific group
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param limit query int false "Number of items per page" default(100)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} models.GroupWordsResponse
// @Router /api/groups/{id}/words [get]
func (h *GroupHandler) GetGroupWords(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid group ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	words, err := h.service.GetGroupWords(groupID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get group words")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, words)
}

// GetGroupStudySessions godoc
// @Summary Get study sessions for a group
// @Description Returns a paginated list of study sessions for a specific group
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param limit query int false "Number of items per page" default(100)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} models.GroupStudySessionsResponse
// @Router /api/groups/{id}/study_sessions [get]
func (h *GroupHandler) GetGroupStudySessions(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid group ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	sessions, err := h.service.GetGroupStudySessions(groupID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get group study sessions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}
