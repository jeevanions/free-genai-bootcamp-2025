package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWordReviewRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request CreateWordReviewRequest
		wantErr bool
	}{
		{
			name: "valid request - translation",
			request: CreateWordReviewRequest{
				WordID:         1,
				StudySessionID: 1,
				Correct:        true,
				ReviewType:     "translation",
			},
			wantErr: false,
		},
		{
			name: "valid request - pronunciation",
			request: CreateWordReviewRequest{
				WordID:         1,
				StudySessionID: 1,
				Correct:        false,
				ReviewType:     "pronunciation",
			},
			wantErr: false,
		},
		{
			name: "invalid review type",
			request: CreateWordReviewRequest{
				WordID:         1,
				StudySessionID: 1,
				Correct:        true,
				ReviewType:     "invalid_type",
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
