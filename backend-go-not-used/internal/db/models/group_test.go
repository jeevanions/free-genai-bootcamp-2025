package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGroup_TableName(t *testing.T) {
	group := Group{}
	assert.Equal(t, "groups", group.TableName())
}

func TestWordGroup_TableName(t *testing.T) {
	wordGroup := WordGroup{}
	assert.Equal(t, "words_groups", wordGroup.TableName())
}

func TestGroup_Fields(t *testing.T) {
	now := time.Now()
	group := Group{
		ID:              1,
		Name:            "Food Vocabulary",
		Description:     "Common Italian words related to food and dining",
		DifficultyLevel: 1,
		Category:        "thematic",
		CreatedAt:       now,
	}

	// Test field types and values
	assert.Equal(t, int64(1), group.ID)
	assert.Equal(t, "Food Vocabulary", group.Name)
	assert.Equal(t, "Common Italian words related to food and dining", group.Description)
	assert.Equal(t, 1, group.DifficultyLevel)
	assert.Equal(t, "thematic", group.Category)
	assert.Equal(t, now, group.CreatedAt)
}

func TestWordGroup_Fields(t *testing.T) {
	now := time.Now()
	wordGroup := WordGroup{
		ID:        1,
		WordID:    2,
		GroupID:   3,
		CreatedAt: now,
	}

	// Test field types and values
	assert.Equal(t, int64(1), wordGroup.ID)
	assert.Equal(t, int64(2), wordGroup.WordID)
	assert.Equal(t, int64(3), wordGroup.GroupID)
	assert.Equal(t, now, wordGroup.CreatedAt)
}
