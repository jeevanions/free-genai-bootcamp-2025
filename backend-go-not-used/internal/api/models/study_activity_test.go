package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateStudyActivityRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request CreateStudyActivityRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateStudyActivityRequest{
				Name:            "Basic Vocabulary Quiz",
				Type:            "vocabulary",
				DifficultyLevel: 1,
				RequiresAudio:   false,
				Instructions:    "Match the Italian words with their English translations",
				ThumbnailURL:    "https://example.com/thumbnail.jpg",
				Category:        "basics",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			request: CreateStudyActivityRequest{
				Name:            "Basic Vocabulary Quiz",
				Type:            "invalid_type",
				DifficultyLevel: 1,
				RequiresAudio:   false,
				Instructions:    "Match the Italian words with their English translations",
				ThumbnailURL:    "https://example.com/thumbnail.jpg",
				Category:        "basics",
			},
			wantErr: true,
		},
		{
			name: "invalid difficulty level - too low",
			request: CreateStudyActivityRequest{
				Name:            "Basic Vocabulary Quiz",
				Type:            "vocabulary",
				DifficultyLevel: 0,
				RequiresAudio:   false,
				Instructions:    "Match the Italian words with their English translations",
				ThumbnailURL:    "https://example.com/thumbnail.jpg",
				Category:        "basics",
			},
			wantErr: true,
		},
		{
			name: "invalid difficulty level - too high",
			request: CreateStudyActivityRequest{
				Name:            "Basic Vocabulary Quiz",
				Type:            "vocabulary",
				DifficultyLevel: 6,
				RequiresAudio:   false,
				Instructions:    "Match the Italian words with their English translations",
				ThumbnailURL:    "https://example.com/thumbnail.jpg",
				Category:        "basics",
			},
			wantErr: true,
		},
		{
			name: "invalid category",
			request: CreateStudyActivityRequest{
				Name:            "Basic Vocabulary Quiz",
				Type:            "vocabulary",
				DifficultyLevel: 1,
				RequiresAudio:   false,
				Instructions:    "Match the Italian words with their English translations",
				ThumbnailURL:    "https://example.com/thumbnail.jpg",
				Category:        "invalid_category",
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
