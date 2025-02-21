package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/domain/services"
)

type WordReviewHandler struct {
	service services.WordReviewService
}

func NewWordReviewHandler(service services.WordReviewService) *WordReviewHandler {
	return &WordReviewHandler{service: service}
}

// CreateReview godoc
// @Summary Create a new word review
// @Description Record a word review result during a study session
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body models.WordReviewItem true "Word review details"
// @Success 201 {object} models.WordReviewItem
// @Failure 400 {object} map[string]string
// @Router /reviews [post]
func (h *WordReviewHandler) CreateReview(c *gin.Context) {
	var review models.WordReviewItem
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateReview(c.Request.Context(), &review); err != nil {
		if err == services.ErrInvalidReview {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// ListSessionReviews godoc
// @Summary List reviews for a study session
// @Description Get all word reviews for a specific study session
// @Tags reviews
// @Produce json
// @Param sessionId path int true "Study Session ID"
// @Success 200 {array} models.WordReviewItem
// @Failure 404 {object} map[string]string
// @Router /reviews/sessions/{sessionId} [get]
func (h *WordReviewHandler) ListSessionReviews(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("sessionId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	reviews, err := h.service.ListSessionReviews(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// GetWordStats godoc
// @Summary Get word review statistics
// @Description Get review statistics for a specific word
// @Tags reviews
// @Produce json
// @Param wordId path int true "Word ID"
// @Success 200 {object} models.WordStats
// @Failure 404 {object} map[string]string
// @Router /reviews/words/{wordId}/stats [get]
func (h *WordReviewHandler) GetWordStats(c *gin.Context) {
	wordID, err := strconv.ParseInt(c.Param("wordId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	stats, err := h.service.GetWordStats(c.Request.Context(), wordID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *WordReviewHandler) RecordReview(c *gin.Context) {
	var review models.WordReviewItem
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.service.RecordReview(c.Request.Context(), review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
