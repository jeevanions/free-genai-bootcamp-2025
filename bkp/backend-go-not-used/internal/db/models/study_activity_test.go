package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStudyActivity_TableName(t *testing.T) {
	activity := StudyActivity{}
	assert.Equal(t, "study_activities", activity.TableName())
}

func TestStudyActivity_Fields(t *testing.T) {
	now := time.Now()
	activity := StudyActivity{
		ID:              1,
		Name:            "Basic Vocabulary Quiz",
		Type:            "vocabulary",
		RequiresAudio:   false,
		DifficultyLevel: 1,
		Instructions:    "Match the Italian words with their English translations",
		ThumbnailURL:    "https://example.com/thumbnail.jpg",
		Category:        "basics",
		CreatedAt:       now,
	}

	// Test field types and values
	assert.Equal(t, int64(1), activity.ID)
	assert.Equal(t, "Basic Vocabulary Quiz", activity.Name)
	assert.Equal(t, "vocabulary", activity.Type)
	assert.Equal(t, false, activity.RequiresAudio)
	assert.Equal(t, 1, activity.DifficultyLevel)
	assert.Equal(t, "Match the Italian words with their English translations", activity.Instructions)
	assert.Equal(t, "https://example.com/thumbnail.jpg", activity.ThumbnailURL)
	assert.Equal(t, "basics", activity.Category)
	assert.Equal(t, now, activity.CreatedAt)
}
