package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/api/models"
	"github.com/jeevanions/italian-learning/internal/domain/services"
)

// WordHandler handles word-related HTTP requests
type WordHandler struct {
	wordService services.WordService
}

// NewWordHandler creates a new WordHandler instance
func NewWordHandler(wordService services.WordService) *WordHandler {
	return &WordHandler{
		wordService: wordService,
	}
}

// ListWords godoc
// @Summary List all words
// @Description Get a list of all Italian vocabulary words
// @Tags words
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {array} models.Word
// @Failure 500 {object} models.ErrorResponse
// @Router /words [get]
func (h *WordHandler) ListWords(c *gin.Context) {
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := c.Query("page_size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	words, err := h.wordService.ListWords(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	if words == nil {
		words = []models.Word{}
	}

	c.JSON(http.StatusOK, gin.H{"words": words})
}

// CreateWord godoc
// @Summary Create a new word
// @Description Create a new Italian vocabulary word
// @Tags words
// @Accept json
// @Produce json
// @Param word body models.CreateWordRequest true "Word object"
// @Success 201 {object} models.Word
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /words [post]
func (h *WordHandler) CreateWord(c *gin.Context) {
	var req models.CreateWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	word, err := h.wordService.CreateWord(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, word)
}

// GetWord godoc
// @Summary Get a word by ID
// @Description Get an Italian vocabulary word by its ID
// @Tags words
// @Accept json
// @Produce json
// @Param id path int true "Word ID"
// @Success 200 {object} models.Word
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /words/{id} [get]
func (h *WordHandler) GetWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid id"})
		return
	}

	word, err := h.wordService.GetWord(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, word)
}
