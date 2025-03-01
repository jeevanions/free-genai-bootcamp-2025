package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeevanions/italian-learning/internal/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite" // Replace mattn/go-sqlite3
)

func TestStudyActivityAPI(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestServer(t, db)

	t.Run("Create Activity", func(t *testing.T) {
		activity := models.CreateStudyActivityRequest{
			Name:            "Basic Vocabulary Quiz",
			Type:            "vocabulary",
			DifficultyLevel: 1,
			RequiresAudio:   false,
			Instructions:    "Match the Italian words with their English translations"}

		body, err := json.Marshal(activity)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/activities/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.StudyActivity
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotZero(t, response.ID)
		assert.Equal(t, activity.Type, response.Type)
	})

	t.Run("List Activities", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/activities/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.StudyActivity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response)
	})
}
