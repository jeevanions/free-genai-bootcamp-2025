package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGroupRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request CreateGroupRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateGroupRequest{
				Name:            "Food Vocabulary",
				Description:     "Common Italian words related to food and dining",
				DifficultyLevel: 1,
				Category:        "thematic",
			},
			wantErr: false,
		},
		{
			name: "invalid difficulty level - too low",
			request: CreateGroupRequest{
				Name:            "Food Vocabulary",
				Description:     "Common Italian words related to food and dining",
				DifficultyLevel: 0,
				Category:        "thematic",
			},
			wantErr: true,
		},
		{
			name: "invalid difficulty level - too high",
			request: CreateGroupRequest{
				Name:            "Food Vocabulary",
				Description:     "Common Italian words related to food and dining",
				DifficultyLevel: 6,
				Category:        "thematic",
			},
			wantErr: true,
		},
		{
			name: "invalid category",
			request: CreateGroupRequest{
				Name:            "Food Vocabulary",
				Description:     "Common Italian words related to food and dining",
				DifficultyLevel: 1,
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
