package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateWord(t *testing.T) {
	tests := []struct {
		name    string
		word    *Word
		wantErr bool
	}{
		{
			name: "valid word",
			word: &Word{
				Italian:         "ciao",
				English:         "hello",
				PartsOfSpeech:   "interjection",
				DifficultyLevel: 1,
			},
			wantErr: false,
		},
		{
			name: "missing required field",
			word: &Word{
				Italian: "",
				English: "hello",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWord(tt.word)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
