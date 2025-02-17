package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apimodels "github.com/jeevanions/italian-learning/internal/api/models"
	dbmodels "github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/db/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWordAPI(t *testing.T) {
	router, _, err := setupTestServer()
	require.NoError(t, err)

	t.Run("Create Word", func(t *testing.T) {
		word := apimodels.CreateWordRequest{
			Italian:         "ciao",
			English:         "hello",
			PartsOfSpeech:   "interjection",
			DifficultyLevel: 1,
		}

		body, err := json.Marshal(word)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/words/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response apimodels.Word
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotZero(t, response.ID)
		assert.Equal(t, word.Italian, response.Italian)
	})

	t.Run("List Words", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/words/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Words []apimodels.Word `json:"words"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Words)
	})
}

func TestWordDatabaseConstraints(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	t.Run("Difficulty Level Constraint", func(t *testing.T) {
		// Test invalid difficulty level
		_, err := db.ExecContext(ctx, `
			INSERT INTO words (
				italian, english, parts_of_speech, difficulty_level
			) VALUES (?, ?, ?, ?)`,
			"test", "test", "noun", 6, // Invalid difficulty level (>5)
		)
		assert.Error(t, err, "Should reject difficulty level > 5")

		_, err = db.ExecContext(ctx, `
			INSERT INTO words (
				italian, english, parts_of_speech, difficulty_level
			) VALUES (?, ?, ?, ?)`,
			"test", "test", "noun", 0, // Invalid difficulty level (<1)
		)
		assert.Error(t, err, "Should reject difficulty level < 1")

		// Test valid difficulty level
		_, err = db.ExecContext(ctx, `
			INSERT INTO words (
				italian, english, parts_of_speech, difficulty_level
			) VALUES (?, ?, ?, ?)`,
			"test_valid_level", "test", "noun", 3, // Valid difficulty level
		)
		assert.NoError(t, err, "Should accept valid difficulty level")
	})

	t.Run("Gender Constraint", func(t *testing.T) {
		// Test invalid gender
		_, err := db.ExecContext(ctx, `
			INSERT INTO words (
				italian, english, parts_of_speech, difficulty_level, gender
			) VALUES (?, ?, ?, ?, ?)`,
			"test", "test", "noun", 1, "invalid",
		)
		assert.Error(t, err, "Should reject invalid gender")

		// Test valid genders
		validGenders := []string{"masculine", "feminine", "neuter"}
		for _, gender := range validGenders {
			_, err := db.ExecContext(ctx, `
				INSERT INTO words (
					italian, english, parts_of_speech, difficulty_level, gender
				) VALUES (?, ?, ?, ?, ?)`,
				"test_gender_"+gender, "test", "noun", 1, gender,
			)
			assert.NoError(t, err, "Should accept valid gender: "+gender)
		}
	})

	t.Run("Required Fields", func(t *testing.T) {
		// Test missing required field
		_, err := db.ExecContext(ctx, `
			INSERT INTO words (
				italian, parts_of_speech, difficulty_level
			) VALUES (?, ?, ?)`,
			"test", "noun", 1, // Missing english field
		)
		assert.Error(t, err, "Should reject missing required field")

		// Test all required fields present
		_, err = db.ExecContext(ctx, `
			INSERT INTO words (
				italian, english, parts_of_speech, difficulty_level
			) VALUES (?, ?, ?, ?)`,
			"test_required_fields", "test", "noun", 1,
		)
		assert.NoError(t, err, "Should accept all required fields")
	})

	t.Run("Unique Italian Word", func(t *testing.T) {
		// Insert first word
		_, err := db.ExecContext(ctx, `
			INSERT INTO words (
				italian, english, parts_of_speech, difficulty_level
			) VALUES (?, ?, ?, ?)`,
			"test_unique_word", "test", "noun", 1,
		)
		require.NoError(t, err)

		// Try to insert duplicate Italian word
		_, err = db.ExecContext(ctx, `
			INSERT INTO words (
				italian, english, parts_of_speech, difficulty_level
			) VALUES (?, ?, ?, ?)`,
			"test_unique_word", "test2", "noun", 1,
		)
		assert.Error(t, err, "Should reject duplicate Italian word")
	})
}

func TestWordRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewSQLiteWordRepository(db)
	ctx := context.Background()

	t.Run("Create and Retrieve Word", func(t *testing.T) {
		word := &dbmodels.Word{
			Italian:         "ciao",
			English:         "hello",
			PartsOfSpeech:   "interjection",
			DifficultyLevel: 1,
		}

		// Create word
		err := repo.Create(ctx, word)
		require.NoError(t, err)
		assert.NotZero(t, word.ID, "Should assign an ID")

		// Retrieve word
		retrieved, err := repo.GetByID(ctx, word.ID)
		require.NoError(t, err)
		assert.Equal(t, word.Italian, retrieved.Italian)
		assert.Equal(t, word.English, retrieved.English)
	})

	t.Run("Update Word", func(t *testing.T) {
		word := &dbmodels.Word{
			Italian:         "test_update",
			English:         "test",
			PartsOfSpeech:   "noun",
			DifficultyLevel: 1,
		}

		// Create word
		err := repo.Create(ctx, word)
		require.NoError(t, err)

		// Update word
		word.English = "updated"
		err = repo.Update(ctx, word)
		require.NoError(t, err)

		// Verify update
		updated, err := repo.GetByID(ctx, word.ID)
		require.NoError(t, err)
		assert.Equal(t, "updated", updated.English)
	})

	t.Run("Delete Word", func(t *testing.T) {
		word := &dbmodels.Word{
			Italian:         "test_delete",
			English:         "test",
			PartsOfSpeech:   "noun",
			DifficultyLevel: 1,
		}

		// Create word
		err := repo.Create(ctx, word)
		require.NoError(t, err)

		// Delete word
		err = repo.Delete(ctx, word.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.GetByID(ctx, word.ID)
		assert.Error(t, err, "Should not find deleted word")
	})
}

func strPtr(s string) *string {
	return &s
}

func convertAPIModelToDBModel(apiWord *apimodels.Word) *dbmodels.Word {
	dbWord := &dbmodels.Word{
		ID:              apiWord.ID,
		Italian:         apiWord.Italian,
		English:         apiWord.English,
		PartsOfSpeech:   apiWord.PartsOfSpeech,
		DifficultyLevel: apiWord.DifficultyLevel,
		CreatedAt:       apiWord.CreatedAt,
	}

	// Safely handle nullable fields
	if apiWord.Gender != nil {
		dbWord.Gender = sql.NullString{String: *apiWord.Gender, Valid: true}
	}
	if apiWord.Number != nil {
		dbWord.Number = sql.NullString{String: *apiWord.Number, Valid: true}
	}
	if apiWord.Notes != nil {
		dbWord.Notes = sql.NullString{String: *apiWord.Notes, Valid: true}
	}
	if apiWord.VerbConjugation != nil {
		dbWord.VerbConjugation = json.RawMessage(*apiWord.VerbConjugation)
	}

	return dbWord
}
