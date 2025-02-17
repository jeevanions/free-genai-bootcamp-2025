package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/api/models"
	"github.com/stretchr/testify/require"
)

func createTestActivity(t *testing.T, router *gin.Engine) *models.StudyActivity {
	activity := models.CreateStudyActivityRequest{
		Name:            "Test Activity",
		Type:            "vocabulary",
		DifficultyLevel: 1,
		RequiresAudio:   false,
		Instructions:    "Test instructions",
	}

	body, err := json.Marshal(activity)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/activities/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response models.StudyActivity
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotZero(t, response.ID)

	return &response
}
