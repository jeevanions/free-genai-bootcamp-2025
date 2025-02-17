package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/db/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWordRepository(t *testing.T) {
	// Setup test database
	db, err := testutil.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewSQLiteWordRepository(db)
	ctx := context.Background()

	t.Run("null_fields_handling", func(t *testing.T) {
		word := &models.Word{
			Italian:         "ciao",
			English:         "hello",
			PartsOfSpeech:   "interjection",
			DifficultyLevel: 1,
			Gender:          sql.NullString{Valid: false},
			Number:          sql.NullString{Valid: false},
			Notes:           sql.NullString{Valid: false},
		}

		err := repo.Create(ctx, word)
		require.NoError(t, err)

		// Verify retrieval
		retrieved, err := repo.GetByID(ctx, word.ID)
		require.NoError(t, err)
		assert.False(t, retrieved.Gender.Valid)
		assert.False(t, retrieved.Number.Valid)
		assert.False(t, retrieved.Notes.Valid)
	})
}
