package models

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWord_TableName(t *testing.T) {
	word := Word{}
	assert.Equal(t, "words", word.TableName())
}

func TestWord_Fields(t *testing.T) {
	now := time.Now()
	verbConj := json.RawMessage(`{"present": {"io": "mangio", "tu": "mangi"}}`)
	
	word := Word{
		ID:              1,
		Italian:         "mangiare",
		English:         "to eat",
		PartsOfSpeech:   "verb",
		Gender:          sql.NullString{String: "", Valid: false},
		Number:          sql.NullString{String: "", Valid: false},
		DifficultyLevel: 1,
		VerbConjugation: verbConj,
		Notes:           sql.NullString{String: "Common verb for eating", Valid: true},
		CreatedAt:       now,
	}

	// Test field types and values
	assert.Equal(t, int64(1), word.ID)
	assert.Equal(t, "mangiare", word.Italian)
	assert.Equal(t, "to eat", word.English)
	assert.Equal(t, "verb", word.PartsOfSpeech)
	assert.False(t, word.Gender.Valid)
	assert.False(t, word.Number.Valid)
	assert.Equal(t, 1, word.DifficultyLevel)
	assert.Equal(t, verbConj, word.VerbConjugation)
	assert.True(t, word.Notes.Valid)
	assert.Equal(t, "Common verb for eating", word.Notes.String)
	assert.Equal(t, now, word.CreatedAt)

	// Test with noun
	noun := Word{
		ID:              2,
		Italian:         "casa",
		English:         "house",
		PartsOfSpeech:   "noun",
		Gender:          sql.NullString{String: "feminine", Valid: true},
		Number:          sql.NullString{String: "singular", Valid: true},
		DifficultyLevel: 1,
		VerbConjugation: nil,
		Notes:           sql.NullString{String: "", Valid: false},
		CreatedAt:       now,
	}

	assert.Equal(t, int64(2), noun.ID)
	assert.Equal(t, "casa", noun.Italian)
	assert.Equal(t, "house", noun.English)
	assert.Equal(t, "noun", noun.PartsOfSpeech)
	assert.True(t, noun.Gender.Valid)
	assert.Equal(t, "feminine", noun.Gender.String)
	assert.True(t, noun.Number.Valid)
	assert.Equal(t, "singular", noun.Number.String)
	assert.Equal(t, 1, noun.DifficultyLevel)
	assert.Nil(t, noun.VerbConjugation)
	assert.False(t, noun.Notes.Valid)
	assert.Equal(t, now, noun.CreatedAt)
}
