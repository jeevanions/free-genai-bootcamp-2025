package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWordRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request CreateWordRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateWordRequest{
				Italian:         "ciao",
				English:         "hello",
				PartsOfSpeech:   "interjection",
				DifficultyLevel: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid parts of speech",
			request: CreateWordRequest{
				Italian:         "ciao",
				English:         "hello",
				PartsOfSpeech:   "invalid",
				DifficultyLevel: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid difficulty level",
			request: CreateWordRequest{
				Italian:         "ciao",
				English:         "hello",
				PartsOfSpeech:   "noun",
				DifficultyLevel: 6,
			},
			wantErr: true,
		},
		{
			name: "invalid gender",
			request: CreateWordRequest{
				Italian:         "ciao",
				English:         "hello",
				PartsOfSpeech:   "noun",
				Gender:          strPtr("invalid"),
				DifficultyLevel: 1,
			},
			wantErr: true,
		},
		{
			name: "invalid number",
			request: CreateWordRequest{
				Italian:         "ciao",
				English:         "hello",
				PartsOfSpeech:   "noun",
				Number:          strPtr("invalid"),
				DifficultyLevel: 1,
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

func TestUpdateWordRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateWordRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: UpdateWordRequest{
				Italian:         strPtr("ciao"),
				English:         strPtr("hello"),
				PartsOfSpeech:   strPtr("interjection"),
				DifficultyLevel: intPtr(1),
			},
			wantErr: false,
		},
		{
			name: "invalid parts of speech",
			request: UpdateWordRequest{
				PartsOfSpeech: strPtr("invalid"),
			},
			wantErr: true,
		},
		{
			name: "invalid difficulty level",
			request: UpdateWordRequest{
				DifficultyLevel: intPtr(6),
			},
			wantErr: true,
		},
		{
			name: "invalid gender",
			request: UpdateWordRequest{
				Gender: strPtr("invalid"),
			},
			wantErr: true,
		},
		{
			name: "invalid number",
			request: UpdateWordRequest{
				Number: strPtr("invalid"),
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

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
