package services

import (
	"context"
	"testing"
	"time"

	"github.com/jeevanions/italian-learning/internal/api/models"
	dbmodels "github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/db/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockWordRepository struct {
	repository.WordRepository
}

func (m *mockWordRepository) Create(ctx context.Context, word *dbmodels.Word) error {
	word.ID = 1 // Mock ID assignment
	word.CreatedAt = time.Now()
	return nil
}

func TestListWords(t *testing.T) {
	// ... (existing code)

	// Add new test cases
	tc := []struct {
		name         string
		expectGender *string
		expectNumber *string
		expectNotes  *string
	}{
		{
			name:         "Test case 1",
			expectGender: nil,
			expectNumber: nil,
			expectNotes:  nil,
		},
		{
			name:         "Test case 2",
			expectGender: nil,
			expectNumber: nil,
			expectNotes:  nil,
		},
		{
			name:         "Test case 3",
			expectGender: nil,
			expectNumber: nil,
			expectNotes:  nil,
		},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			// ... (existing code)

			// ... (existing assertions)

			// ... (new assertions)
		})
	}
}

func TestGetWord(t *testing.T) {
	// ... (existing code)

	// Add new test cases
	tc := []struct {
		name         string
		expectGender *string
		expectNumber *string
		expectNotes  *string
	}{
		{
			name:         "Test case 1",
			expectGender: nil,
			expectNumber: nil,
			expectNotes:  nil,
		},
		{
			name:         "Test case 2",
			expectGender: nil,
			expectNumber: nil,
			expectNotes:  nil,
		},
		{
			name:         "Test case 3",
			expectGender: nil,
			expectNumber: nil,
			expectNotes:  nil,
		},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			// ... (existing code)

			// ... (existing assertions)

			// ... (new assertions)
		})
	}
}

func TestCreateWord(t *testing.T) {
	testCases := []struct {
		name     string
		input    *models.CreateWordRequest
		wantErr  bool
		validate func(*testing.T, *models.Word)
	}{
		{
			name: "valid word with all fields",
			input: &models.CreateWordRequest{
				Italian:         "essere",
				English:         "to be",
				PartsOfSpeech:   "verb",
				DifficultyLevel: 1,
				VerbConjugation: strPtr(`{"present":{"io":"sono","tu":"sei"}}`),
				Gender:          strPtr("masculine"),
				Number:          strPtr("singular"),
				Notes:           strPtr("Most common verb"),
			},
			validate: func(t *testing.T, w *models.Word) {
				assert.NotNil(t, w.VerbConjugation)
				assert.NotNil(t, w.Gender)
				assert.NotNil(t, w.Number)
				assert.NotNil(t, w.Notes)
			},
		},
		{
			name: "valid word with only required fields",
			input: &models.CreateWordRequest{
				Italian:         "ciao",
				English:         "hello",
				PartsOfSpeech:   "interjection",
				DifficultyLevel: 1,
			},
			validate: func(t *testing.T, w *models.Word) {
				assert.Nil(t, w.VerbConjugation)
				assert.Nil(t, w.Gender)
				assert.Nil(t, w.Number)
				assert.Nil(t, w.Notes)
			},
		},
		{
			name: "invalid difficulty level",
			input: &models.CreateWordRequest{
				Italian:         "ciao",
				English:         "hello",
				PartsOfSpeech:   "interjection",
				DifficultyLevel: 6, // Invalid: should be 1-5
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ... (existing code)

			// ... (existing assertions)

			// ... (new assertions)
		})
	}
}

func TestUpdateWord(t *testing.T) {
	// ... (existing code)

	// Add new test cases
	tc := []struct {
		name         string
		expectGender *string
		expectNumber *string
		expectNotes  *string
	}{
		{
			name:         "Test case 1",
			expectGender: nil,
			expectNumber: nil,
			expectNotes:  nil,
		},
		{
			name:         "Test case 2",
			expectGender: nil,
			expectNumber: nil,
			expectNotes:  nil,
		},
		{
			name:         "Test case 3",
			expectGender: nil,
			expectNumber: nil,
			expectNotes:  nil,
		},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			// ... (existing code)

			// ... (existing assertions)

			// ... (new assertions)
		})
	}
}

func TestDeleteWord(t *testing.T) {
	// ... (existing code)

	// Add new test cases
	tc := []struct {
		name         string
		expectGender *string
		expectNumber *string
		expectNotes  *string
	}{
		{
			name:         "Test case 1",
			expectGender: nil,
			expectNumber: nil,
			expectNotes:  nil,
		},
		{
			name:         "Test case 2",
			expectGender: nil,
			expectNumber: nil,
			expectNotes:  nil,
		},
		{
			name:         "Test case 3",
			expectGender: nil,
			expectNumber: nil,
			expectNotes:  nil,
		},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			// ... (existing code)

			// ... (existing assertions)

			// ... (new assertions)
		})
	}
}

func TestWordService_CreateWord(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockWordRepository{}
	service := NewWordService(mockRepo)

	testCases := []struct {
		name  string
		input *models.CreateWordRequest
	}{
		{
			name: "Basic word without optional fields",
			input: &models.CreateWordRequest{
				Italian:         "ciao",
				English:         "hello",
				PartsOfSpeech:   "interjection",
				DifficultyLevel: 1,
			},
		},
		{
			name: "Word with all optional fields",
			input: &models.CreateWordRequest{
				Italian:         "gatto",
				English:         "cat",
				PartsOfSpeech:   "noun",
				Gender:          strPtr("masculine"),
				Number:          strPtr("singular"),
				DifficultyLevel: 1,
				Notes:           strPtr("Common pet animal"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.CreateWord(ctx, tc.input)
			require.NoError(t, err)
			assert.NotNil(t, result)

			// Add these assertions
			if tc.input.Gender != nil {
				assert.NotNil(t, result.Gender)
				assert.Equal(t, *tc.input.Gender, *result.Gender)
			} else {
				assert.Nil(t, result.Gender)
			}

			if tc.input.Number != nil {
				assert.NotNil(t, result.Number)
				assert.Equal(t, *tc.input.Number, *result.Number)
			} else {
				assert.Nil(t, result.Number)
			}

			if tc.input.Notes != nil {
				assert.NotNil(t, result.Notes)
				assert.Equal(t, *tc.input.Notes, *result.Notes)
			} else {
				assert.Nil(t, result.Notes)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
