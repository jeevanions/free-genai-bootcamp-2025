package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/domain/services"
)

type GroupHandler struct {
	service services.GroupService
}

func NewGroupHandler(service services.GroupService) *GroupHandler {
	return &GroupHandler{service: service}
}

// GetGroup godoc
// @Summary Get a group by ID
// @Description Retrieve a study group by its ID
// @Tags groups
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} models.Group
// @Failure 404 {object} map[string]string
// @Router /groups/{id} [get]
func (h *GroupHandler) GetGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	group, err := h.service.GetGroup(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// CreateGroup godoc
// @Summary Create a new group
// @Description Create a new study group
// @Tags groups
// @Accept json
// @Produce json
// @Param group body models.Group true "Group object"
// @Success 201 {object} models.Group
// @Failure 400 {object} map[string]string
// @Router /groups [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateGroup(c.Request.Context(), &group); err != nil {
		if err == services.ErrInvalidGroup {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// ListGroups godoc
// @Summary List all groups
// @Description Get a paginated list of study groups
// @Tags groups
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {array} models.Group
// @Failure 500 {object} map[string]string
// @Router /groups [get]
func (h *GroupHandler) ListGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	groups, err := h.service.ListGroups(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}
