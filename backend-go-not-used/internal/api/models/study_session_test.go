package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateStudySessionRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request CreateStudySessionRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateStudySessionRequest{
				GroupID:         1,
				StudyActivityID: 1,
				TotalWords:      10,
				CorrectWords:    8,
				DurationSeconds: 300,
			},
			wantErr: false,
		},
		{
			name: "valid request - zero correct words",
			request: CreateStudySessionRequest{
				GroupID:         1,
				StudyActivityID: 1,
				TotalWords:      10,
				CorrectWords:    0,
				DurationSeconds: 300,
			},
			wantErr: false,
		},
		{
			name: "invalid total words - zero",
			request: CreateStudySessionRequest{
				GroupID:         1,
				StudyActivityID: 1,
				TotalWords:      0,
				CorrectWords:    8,
				DurationSeconds: 300,
			},
			wantErr: true,
		},
		{
			name: "invalid total words - negative",
			request: CreateStudySessionRequest{
				GroupID:         1,
				StudyActivityID: 1,
				TotalWords:      -5,
				CorrectWords:    8,
				DurationSeconds: 300,
			},
			wantErr: true,
		},
		{
			name: "invalid correct words - negative",
			request: CreateStudySessionRequest{
				GroupID:         1,
				StudyActivityID: 1,
				TotalWords:      10,
				CorrectWords:    -1,
				DurationSeconds: 300,
			},
			wantErr: true,
		},
		{
			name: "invalid correct words - greater than total",
			request: CreateStudySessionRequest{
				GroupID:         1,
				StudyActivityID: 1,
				TotalWords:      10,
				CorrectWords:    11,
				DurationSeconds: 300,
			},
			wantErr: true,
		},
		{
			name: "invalid duration - zero",
			request: CreateStudySessionRequest{
				GroupID:         1,
				StudyActivityID: 1,
				TotalWords:      10,
				CorrectWords:    8,
				DurationSeconds: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid duration - negative",
			request: CreateStudySessionRequest{
				GroupID:         1,
				StudyActivityID: 1,
				TotalWords:      10,
				CorrectWords:    8,
				DurationSeconds: -300,
			},
			wantErr: true,
		},
		{
			name: "invalid group ID - zero",
			request: CreateStudySessionRequest{
				GroupID:         0,
				StudyActivityID: 1,
				TotalWords:      10,
				CorrectWords:    8,
				DurationSeconds: 300,
			},
			wantErr: true,
		},
		{
			name: "invalid activity ID - zero",
			request: CreateStudySessionRequest{
				GroupID:         1,
				StudyActivityID: 0,
				TotalWords:      10,
				CorrectWords:    8,
				DurationSeconds: 300,
			},
			wantErr: true,
		},
				StudyActivityID: 1,
				TotalWords:      10,
				CorrectWords:    8,
				DurationSeconds: 0,
			},
			wantErr: true,
		},
		{
			name: "correct words greater than total words",
			request: CreateStudySessionRequest{
				GroupID:         1,
				StudyActivityID: 1,
				TotalWords:      10,
				CorrectWords:    11,
				DurationSeconds: 300,
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
