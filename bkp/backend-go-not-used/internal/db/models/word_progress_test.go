package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWordProgress_TableName(t *testing.T) {
	progress := WordProgress{}
	assert.Equal(t, "word_progress", progress.TableName())
}

func TestWordProgress_Fields(t *testing.T) {
	now := time.Now()
	nextReview := now.Add(24 * time.Hour)
	
	progress := WordProgress{
		ID:              1,
		UserID:          2,
		WordID:          3,
		MasteryLevel:    3,
		NextReviewAt:    nextReview,
		TotalAttempts:   10,
		CorrectAttempts: 7,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Test field types and values
	assert.Equal(t, int64(1), progress.ID)
	assert.Equal(t, int64(2), progress.UserID)
	assert.Equal(t, int64(3), progress.WordID)
	assert.Equal(t, 3, progress.MasteryLevel)
	assert.Equal(t, nextReview, progress.NextReviewAt)
	assert.Equal(t, 10, progress.TotalAttempts)
	assert.Equal(t, 7, progress.CorrectAttempts)
	assert.Equal(t, now, progress.CreatedAt)
	assert.Equal(t, now, progress.UpdatedAt)
}

func TestWordProgress_CalculateAccuracy(t *testing.T) {
	tests := []struct {
		name     string
		progress WordProgress
		want     float64
	}{
		{
			name: "zero attempts",
			progress: WordProgress{
				TotalAttempts:   0,
				CorrectAttempts: 0,
			},
			want: 0,
		},
		{
			name: "all correct",
			progress: WordProgress{
				TotalAttempts:   10,
				CorrectAttempts: 10,
			},
			want: 100,
		},
		{
			name: "70% correct",
			progress: WordProgress{
				TotalAttempts:   10,
				CorrectAttempts: 7,
			},
			want: 70,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.progress.CalculateAccuracy()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWordProgress_UpdateMasteryLevel(t *testing.T) {
	tests := []struct {
		name           string
		progress      *WordProgress
		wantMastery   int
	}{
		{
			name: "not enough attempts",
			progress: &WordProgress{
				TotalAttempts:   3,
				CorrectAttempts: 3,
				MasteryLevel:    0,
			},
			wantMastery: 0,
		},
		{
			name: "mastery level 5",
			progress: &WordProgress{
				TotalAttempts:   10,
				CorrectAttempts: 9, // 90%
				MasteryLevel:    4,
			},
			wantMastery: 5,
		},
		{
			name: "mastery level 3",
			progress: &WordProgress{
				TotalAttempts:   10,
				CorrectAttempts: 7, // 70%
				MasteryLevel:    2,
			},
			wantMastery: 3,
		},
		{
			name: "mastery level 0",
			progress: &WordProgress{
				TotalAttempts:   10,
				CorrectAttempts: 4, // 40%
				MasteryLevel:    1,
			},
			wantMastery: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.progress.UpdateMasteryLevel()
			assert.Equal(t, tt.wantMastery, tt.progress.MasteryLevel)
		})
	}
}

func TestWordProgress_UpdateNextReviewTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name         string
		progress    *WordProgress
		wantMinDiff time.Duration
		wantMaxDiff time.Duration
	}{
		{
			name: "mastery level 5",
			progress: &WordProgress{
				MasteryLevel: 5,
			},
			wantMinDiff: 27 * 24 * time.Hour,  // ~1 month (allowing for some variance)
			wantMaxDiff: 31 * 24 * time.Hour,
		},
		{
			name: "mastery level 3",
			progress: &WordProgress{
				MasteryLevel: 3,
			},
			wantMinDiff: 6 * 24 * time.Hour,   // 1 week
			wantMaxDiff: 8 * 24 * time.Hour,
		},
		{
			name: "mastery level 0",
			progress: &WordProgress{
				MasteryLevel: 0,
			},
			wantMinDiff: 0,
			wantMaxDiff: time.Hour, // Small buffer for test execution
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.progress.UpdateNextReviewTime()
			diff := tt.progress.NextReviewAt.Sub(now)
			assert.GreaterOrEqual(t, diff, tt.wantMinDiff)
			assert.LessOrEqual(t, diff, tt.wantMaxDiff)
		})
	}
}
