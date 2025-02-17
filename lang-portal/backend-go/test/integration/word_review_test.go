package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWordReviewAPI(t *testing.T) {
	router, _, err := setupTestServer()
	require.NoError(t, err)

	// Create test dependencies
	word := createTestWord(t, router)
	group := createTestGroup(t, router)
	activity := createTestActivity(t, router)
	session := createTestSession(t, router, group.ID, activity.ID)

	t.Run("Create Review", func(t *testing.T) {
		review := models.CreateWordReviewRequest{ // Use correct request model
			WordID:         word.ID,
			StudySessionID: session.ID,
			Correct:        true,
		}

		body, _ := json.Marshal(review)
		req := httptest.NewRequest("POST", "/api/v1/reviews/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.WordReviewItem
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotZero(t, response.ID)
		assert.Equal(t, review.WordID, response.WordID)
		assert.Equal(t, review.Correct, response.Correct)
	})

	t.Run("List Session Reviews", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/reviews/sessions/%d/", session.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var reviews []*models.WordReviewItem
		err := json.Unmarshal(w.Body.Bytes(), &reviews)
		assert.NoError(t, err)
		assert.Len(t, reviews, 1)
		assert.Equal(t, word.ID, reviews[0].WordID)
		assert.NotNil(t, reviews[0].WordDetails)
		assert.Equal(t, word.Italian, reviews[0].WordDetails.Italian)
	})

	t.Run("Get Word Stats", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/reviews/words/%d/stats/", word.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var stats models.WordStats
		err := json.Unmarshal(w.Body.Bytes(), &stats)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), stats.TotalReviews)
		assert.Equal(t, int64(1), stats.CorrectReviews)
		assert.Equal(t, float64(100), stats.Accuracy)
	})
}

func createTestWord(t *testing.T, router *gin.Engine) *models.Word {
	word := models.Word{
		Italian:         "ciao",
		English:         "hello",
		PartsOfSpeech:   "interjection",
		DifficultyLevel: 1,
	}

	body, _ := json.Marshal(word)
	req := httptest.NewRequest("POST", "/api/v1/words/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response models.Word
	json.Unmarshal(w.Body.Bytes(), &response)
	return &response
}

func createTestSession(t *testing.T, router *gin.Engine, groupID, activityID int64) *models.StudySession {
	session := models.CreateStudySessionRequest{
		GroupID:         groupID,
		StudyActivityID: activityID,
		TotalWords:      10,
		CorrectWords:    8,
		DurationSeconds: 300,
	}

	body, _ := json.Marshal(session)
	req := httptest.NewRequest("POST", "/api/v1/sessions/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response models.StudySession
	json.Unmarshal(w.Body.Bytes(), &response)
	return &response
}

