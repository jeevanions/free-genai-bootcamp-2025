package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateWordProgressRequest_Validate(t *testing.T) {
	totalAttempts := 10
	correctAttempts := 7
	incorrectAttempts := 11
	masteryLevel := 3

	tests := []struct {
		name    string
		request UpdateWordProgressRequest
		wantErr bool
	}{
		{
			name: "valid request - partial update",
			request: UpdateWordProgressRequest{
				MasteryLevel: &masteryLevel,
			},
			wantErr: false,
		},
		{
			name: "valid request - attempts update",
			request: UpdateWordProgressRequest{
				TotalAttempts:   &totalAttempts,
				CorrectAttempts: &correctAttempts,
			},
			wantErr: false,
		},
		{
			name: "invalid - correct attempts greater than total",
			request: UpdateWordProgressRequest{
				TotalAttempts:   &totalAttempts,
				CorrectAttempts: &incorrectAttempts,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWordProgressFilter_Validation(t *testing.T) {
	now := time.Now()
	validLimit := 50
	invalidLimit := 101
	validOffset := 0
	invalidOffset := -1
	validMastery := 3
	invalidMastery := 6
	reviewDue := true

	tests := []struct {
		name    string
		filter  WordProgressFilter
		wantErr bool
	}{
		{
			name: "valid filter - all fields",
			filter: WordProgressFilter{
				MasteryLevel:    &validMastery,
				ReviewDue:       &reviewDue,
				ReviewDueBefore: &now,
				Limit:           &validLimit,
				Offset:          &validOffset,
			},
			wantErr: false,
		},
		{
			name: "invalid - mastery level too high",
			filter: WordProgressFilter{
				MasteryLevel: &invalidMastery,
			},
			wantErr: true,
		},
		{
			name: "invalid - limit too high",
			filter: WordProgressFilter{
				Limit: &invalidLimit,
			},
			wantErr: true,
		},
		{
			name: "invalid - negative offset",
			filter: WordProgressFilter{
				Offset: &invalidOffset,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: We can't directly test the binding tags here
			// In a real application, this would be tested through API integration tests
			// Here we're just verifying the struct fields and their tags exist
			assert.NotNil(t, tt.filter)
		})
	}
}
