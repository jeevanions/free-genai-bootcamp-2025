package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		wantLevel   zerolog.Level
	}{
		{
			name:        "development environment sets debug level",
			environment: "development",
			wantLevel:   zerolog.DebugLevel,
		},
		{
			name:        "production environment sets info level",
			environment: "production",
			wantLevel:   zerolog.InfoLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original environment
			originalEnv := os.Getenv("ENVIRONMENT")
			defer os.Setenv("ENVIRONMENT", originalEnv)

			// Clear environment
			os.Clearenv()
			os.Setenv("ENVIRONMENT", tt.environment)

			// Reset global logger state
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			log.Logger = zerolog.Nop()

			// Run setup
			Setup()

			// Assert log level
			assert.Equal(t, tt.wantLevel, zerolog.GlobalLevel())
		})
	}
}

func TestSetup_LogLevels(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		wantLevel   zerolog.Level
		wantLog     string
	}{
		{
			name:        "debug level in development",
			environment: "development",
			wantLevel:   zerolog.DebugLevel,
			wantLog:     "debug message",
		},
		{
			name:        "info level in production",
			environment: "production",
			wantLevel:   zerolog.InfoLevel,
			wantLog:     "debug message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Backup and restore environment
			originalEnv := os.Getenv("ENVIRONMENT")
			defer os.Setenv("ENVIRONMENT", originalEnv)

			// Create a buffer to capture log output
			var buf bytes.Buffer

			// Clear environment and set test environment
			os.Clearenv()
			os.Setenv("ENVIRONMENT", tt.environment)

			// Reset logger state
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			log.Logger = zerolog.Nop()

			// Configure logger with our buffer
			Setup(&buf)

			// Test logging
			log.Debug().Msg("debug message")

			// Get log output
			logOutput := buf.String()

			// Verify results
			assert.Equal(t, tt.wantLevel, zerolog.GlobalLevel())
			if tt.wantLevel <= zerolog.DebugLevel {
				assert.Contains(t, strings.ToLower(logOutput), tt.wantLog)
			} else {
				assert.NotContains(t, strings.ToLower(logOutput), tt.wantLog)
			}
		})
	}
}
