package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserSettings_TableName(t *testing.T) {
	settings := UserSettings{}
	assert.Equal(t, "user_settings", settings.TableName())
}

func TestUserSettings_Fields(t *testing.T) {
	now := time.Now()
	settings := UserSettings{
		ID:                  1,
		UserID:             2,
		DailyWordGoal:      15,
		PreferredStudyTime: "Morning",
		NotificationEnabled: true,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	// Test field types and values
	assert.Equal(t, int64(1), settings.ID)
	assert.Equal(t, int64(2), settings.UserID)
	assert.Equal(t, 15, settings.DailyWordGoal)
	assert.Equal(t, "Morning", settings.PreferredStudyTime)
	assert.True(t, settings.NotificationEnabled)
	assert.Equal(t, now, settings.CreatedAt)
	assert.Equal(t, now, settings.UpdatedAt)
}

func TestDefaultSettings(t *testing.T) {
	userID := int64(1)
	settings := DefaultSettings(userID)

	assert.Equal(t, userID, settings.UserID)
	assert.Equal(t, 10, settings.DailyWordGoal)
	assert.Equal(t, "Morning", settings.PreferredStudyTime)
	assert.True(t, settings.NotificationEnabled)
	assert.False(t, settings.CreatedAt.IsZero())
	assert.False(t, settings.UpdatedAt.IsZero())
}

func TestValidStudyTimes(t *testing.T) {
	expected := []string{"Morning", "Afternoon", "Evening"}
	assert.Equal(t, expected, ValidStudyTimes)
}
