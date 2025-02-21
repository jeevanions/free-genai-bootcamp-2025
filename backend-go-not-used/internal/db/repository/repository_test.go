package repository

import (
	"context"
	"testing"

	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestWordRepositoryOperations(t *testing.T) {
	db := SetupTestDB(t)
	repo := NewSQLiteWordRepository(db)

	t.Run("Create and Get Word", func(t *testing.T) {
		word := &models.Word{
			Italian:         "ciao",
			English:         "hello",
			PartsOfSpeech:   "interjection",
			DifficultyLevel: 1,
		}

		// Test Create
		err := repo.Create(context.Background(), word)
		assert.NoError(t, err)
		assert.Greater(t, word.ID, int64(0))

		// Test Get
		retrieved, err := repo.GetByID(context.Background(), word.ID)
		assert.NoError(t, err)
		assert.Equal(t, word.Italian, retrieved.Italian)
		assert.Equal(t, word.English, retrieved.English)
	})

	// Add more test cases...
}
