// Package handlers provides HTTP request handlers for the API.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/services"
)

// ErrorResponse represents an error response
// swagger:model
type ErrorResponse struct {
	// Error message
	// required: true
	Error string `json:"error" example:"Invalid request format"`
}

type LLMHandler struct {
	service services.LLMServiceInterface
}

func NewLLMHandler(service services.LLMServiceInterface) *LLMHandler {
	return &LLMHandler{service: service}
}

// GenerateWords godoc
// @Summary Generate Italian words for a thematic category
// @Description Uses LLM to generate Italian words with translations and grammatical details for a given thematic category
// @Tags words
// @Accept json
// @Produce json
// @Param request body models.GenerateWordsRequest true "Category for word generation"
// @Success 200 {object} models.GenerateWordsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/words/llm/generate-words [post]
func (h *LLMHandler) GenerateWords(c *gin.Context) {
	var req models.GenerateWordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
		return
	}

	response, err := h.service.GenerateWords(req.Category)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate words")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate words"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateThematicGroup godoc
// @Summary Add words to a group
// @Description Adds new words to an existing group
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param request body models.AddWordsToGroupRequest true "Words to add"
// @Success 200 {object} models.AddWordsToGroupResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/groups/{id}/words [post]
func (h *LLMHandler) CreateThematicGroup(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid group ID")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid group ID"})
		return
	}

	var req models.AddWordsToGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
		return
	}

	// Check if group exists
	_, err = h.service.GetGroupByID(groupID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check group existence")
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Group not found"})
		return
	}

	// Add words to the group
	wordsAdded := 0
	for _, word := range req.Words {
		wordID, err := h.service.CreateWord(&word)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create word")
			continue
		}

		if err := h.service.AddWordToGroup(wordID, groupID); err != nil {
			log.Error().Err(err).Msg("Failed to add word to group")
			continue
		}

		wordsAdded++
	}

	// Update group words count
	if err := h.service.UpdateGroupWordsCount(groupID); err != nil {
		log.Error().Err(err).Msg("Failed to update group words count")
	}

	c.JSON(http.StatusOK, models.AddWordsToGroupResponse{
		Success:    wordsAdded > 0,
		WordsAdded: wordsAdded,
	})
}
