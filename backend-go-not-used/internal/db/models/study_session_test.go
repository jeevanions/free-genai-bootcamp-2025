package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStudySession_TableName(t *testing.T) {
	session := StudySession{}
	assert.Equal(t, "study_sessions", session.TableName())
}

func TestStudySession_Fields(t *testing.T) {
	now := time.Now()
	session := StudySession{
		ID:              1,
		GroupID:         2,
		StudyActivityID: 3,
		TotalWords:      10,
		CorrectWords:    8,
		DurationSeconds: 300,
		StartTime:       now,
		EndTime:         now.Add(5 * time.Minute),
		CreatedAt:       now,
	}

	// Test field types and values
	assert.Equal(t, int64(1), session.ID)
	assert.Equal(t, int64(2), session.GroupID)
	assert.Equal(t, int64(3), session.StudyActivityID)
	assert.Equal(t, 10, session.TotalWords)
	assert.Equal(t, 8, session.CorrectWords)
	assert.Equal(t, 300, session.DurationSeconds)
	assert.Equal(t, now, session.StartTime)
	assert.Equal(t, now.Add(5*time.Minute), session.EndTime)
	assert.Equal(t, now, session.CreatedAt)
}
