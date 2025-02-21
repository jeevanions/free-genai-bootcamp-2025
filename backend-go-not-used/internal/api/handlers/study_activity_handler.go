package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/domain/services"
)

type StudyActivityHandler struct {
	service services.StudyActivityService
}

func NewStudyActivityHandler(service services.StudyActivityService) *StudyActivityHandler {
	return &StudyActivityHandler{service: service}
}

func (h *StudyActivityHandler) CreateActivity(c *gin.Context) {
	var activity models.StudyActivity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateActivity(c.Request.Context(), &activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity"})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

func (h *StudyActivityHandler) GetActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	activity, err := h.service.GetActivity(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

func (h *StudyActivityHandler) ListActivities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	activities, err := h.service.ListActivities(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}

	c.JSON(http.StatusOK, activities)
}

// GetCategories returns a list of available study activity categories
func (h *StudyActivityHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// GetRecommended returns a list of recommended study activities based on user progress
func (h *StudyActivityHandler) GetRecommended(c *gin.Context) {
	count, _ := strconv.Atoi(c.DefaultQuery("count", "5"))

	activities, err := h.service.GetRecommended(c.Request.Context(), count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommended activities"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recommended": activities})
}
