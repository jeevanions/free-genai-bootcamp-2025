package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateUserSettingsRequest_Validate(t *testing.T) {
	goal := 15
	validTime := "Morning"
	invalidTime := "Night"
	enabled := true

	tests := []struct {
		name    string
		request UpdateUserSettingsRequest
		wantErr bool
	}{
		{
			name: "valid request - all fields",
			request: UpdateUserSettingsRequest{
				DailyWordGoal:      &goal,
				PreferredStudyTime: &validTime,
				NotificationEnabled: &enabled,
			},
			wantErr: false,
		},
		{
			name: "valid request - daily word goal only",
			request: UpdateUserSettingsRequest{
				DailyWordGoal: &goal,
			},
			wantErr: false,
		},
		{
			name: "valid request - preferred study time only",
			request: UpdateUserSettingsRequest{
				PreferredStudyTime: &validTime,
			},
			wantErr: false,
		},
		{
			name: "valid request - notification enabled only",
			request: UpdateUserSettingsRequest{
				NotificationEnabled: &enabled,
			},
			wantErr: false,
		},
		{
			name: "invalid preferred study time",
			request: UpdateUserSettingsRequest{
				PreferredStudyTime: &invalidTime,
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
