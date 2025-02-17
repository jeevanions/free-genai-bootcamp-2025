package seed

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jeevanions/italian-learning/internal/db/repository"
	"github.com/jeevanions/italian-learning/internal/db/testutil"
	"github.com/jeevanions/italian-learning/internal/domain/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeedData(t *testing.T) {
	db, err := testutil.NewTestDB()
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()

	// Setup test DB and service
	wordRepo := repository.NewSQLiteWordRepository(db)
	wordService := services.NewWordService(wordRepo)
	groupRepo := repository.NewSQLiteGroupRepository(db)
	groupService := services.NewGroupService(groupRepo)

	// Seed data
	err = SeedBasicData(wordService, groupService)
	require.NoError(t, err)

	t.Run("verb_conjugation_seeding", func(t *testing.T) {
		word, err := wordService.GetWord(ctx, 1)
		require.NoError(t, err)

		// First verify the word was found
		require.NotNil(t, word)

		// Print the word details for debugging
		t.Logf("Word details: ID=%d, Italian=%s, VerbConjugation=%v",
			word.ID, word.Italian, word.VerbConjugation != nil)

		// Verify the verb conjugation exists
		require.NotNil(t, word.VerbConjugation, "VerbConjugation should not be nil")

		// Only try to unmarshal if we have data
		if word.VerbConjugation != nil {
			var conjugation map[string]interface{}
			err = json.Unmarshal([]byte(*word.VerbConjugation), &conjugation)
			assert.NoError(t, err)
			assert.NotEmpty(t, conjugation, "Conjugation data should not be empty")
		}
	})
}
