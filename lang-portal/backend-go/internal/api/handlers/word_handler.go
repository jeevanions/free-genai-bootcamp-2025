package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/services"
)

type WordHandler struct {
	service services.WordServiceInterface
}

func NewWordHandler(service services.WordServiceInterface) *WordHandler {
	return &WordHandler{service: service}
}

// GetWords godoc
// @Summary Get all words
// @Description Returns a paginated list of words
// @Tags words
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page" default(100)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} models.WordListResponse
// @Router /api/words [get]
func (h *WordHandler) GetWords(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	words, err := h.service.GetWords(limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get words")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, words)
}

// GetWordByID godoc
// @Summary Get word by ID
// @Description Returns details about a specific word
// @Tags words
// @Accept json
// @Produce json
// @Param id path int true "Word ID"
// @Success 200 {object} models.WordResponse
// @Router /api/words/{id} [get]
func (h *WordHandler) GetWordByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid word ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	word, err := h.service.GetWordByID(id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get word")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if word == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
		return
	}

	c.JSON(http.StatusOK, word)
}
