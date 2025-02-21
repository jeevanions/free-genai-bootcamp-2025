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

func TestStudySessionAPI(t *testing.T) {
	db := setupTestDB(t)
	router := setupTestServer(t, db)

	// Create test dependencies
	group := createTestGroup(t, router)
	activity := createTestActivity(t, router)

	t.Run("Create Session", func(t *testing.T) {
		session := models.CreateStudySessionRequest{
			GroupID:         group.ID,
			StudyActivityID: activity.ID, // Add missing required field
			TotalWords:      10,
			CorrectWords:    8,
			DurationSeconds: 300,
		}

		body, err := json.Marshal(session)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/sessions/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.StudySessionResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotZero(t, response.ID)
		assert.Equal(t, session.TotalWords, response.TotalWords)
	})

	t.Run("Get Group Stats", func(t *testing.T) {
		// Test empty group first
		emptyGroup := createTestGroup(t, router)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/sessions/groups/%d/stats/", emptyGroup.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var emptyStats models.GroupStats
		err := json.Unmarshal(w.Body.Bytes(), &emptyStats)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), emptyStats.TotalSessions)
		assert.Equal(t, 0, emptyStats.TotalWords)
		assert.Equal(t, 0, emptyStats.CorrectWords)
		assert.Equal(t, float64(0), emptyStats.AverageAccuracy)

		// Test group with session
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/sessions/groups/%d/stats/", group.ID), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var stats models.GroupStats
		err = json.Unmarshal(w.Body.Bytes(), &stats)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), stats.TotalSessions)
		assert.Equal(t, 10, stats.TotalWords)
		assert.Equal(t, 8, stats.CorrectWords)
		assert.Equal(t, float64(80), stats.AverageAccuracy) // 8/10 * 100
	})

	t.Run("List Group Sessions", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/sessions/groups/%d/", group.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var sessions []models.StudySession
		err := json.Unmarshal(w.Body.Bytes(), &sessions)
		assert.NoError(t, err)
		assert.Len(t, sessions, 1)
	})
}

func createTestGroup(t *testing.T, router *gin.Engine) *models.Group {
	group := models.CreateGroupRequest{
		Name:            "Test Group",
		Description:     "Test Description",
		DifficultyLevel: 1,
		Category:        "Vocabulary",
	}

	body, err := json.Marshal(group)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/groups/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response models.Group
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotZero(t, response.ID)

	return &response
}

