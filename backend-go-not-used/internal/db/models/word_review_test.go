package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWordReviewItem_TableName(t *testing.T) {
	review := WordReviewItem{}
	assert.Equal(t, "word_review_items", review.TableName())
}

func TestWordReviewItem_Fields(t *testing.T) {
	now := time.Now()
	review := WordReviewItem{
		ID:             1,
		WordID:         2,
		StudySessionID: 3,
		Correct:        true,
		ReviewType:     "translation",
		CreatedAt:      now,
	}

	// Test field types and values
	assert.Equal(t, int64(1), review.ID)
	assert.Equal(t, int64(2), review.WordID)
	assert.Equal(t, int64(3), review.StudySessionID)
	assert.Equal(t, true, review.Correct)
	assert.Equal(t, "translation", review.ReviewType)
	assert.Equal(t, now, review.CreatedAt)
}

func TestWordReviewStats_CalculateAccuracy(t *testing.T) {
	tests := []struct {
		name           string
		stats         *WordReviewStats
		wantAccuracy float64
	}{
		{
			name: "zero attempts",
			stats: &WordReviewStats{
				WordID:          1,
				TotalAttempts:   0,
				CorrectAttempts: 0,
			},
			wantAccuracy: 0,
		},
		{
			name: "all correct",
			stats: &WordReviewStats{
				WordID:          1,
				TotalAttempts:   10,
				CorrectAttempts: 10,
			},
			wantAccuracy: 100,
		},
		{
			name: "half correct",
			stats: &WordReviewStats{
				WordID:          1,
				TotalAttempts:   10,
				CorrectAttempts: 5,
			},
			wantAccuracy: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stats.CalculateAccuracy()
			assert.Equal(t, tt.wantAccuracy, tt.stats.AccuracyRate)
		})
	}
}
