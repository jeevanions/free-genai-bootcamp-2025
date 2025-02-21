package logger

import (
	"bytes"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	// Save original environment and restore after test
	originalEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", originalEnv)

	t.Run("development environment sets debug level", func(t *testing.T) {
		os.Setenv("APP_ENV", "development")
		Setup()
		assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
	})

	t.Run("production environment sets info level", func(t *testing.T) {
		os.Setenv("APP_ENV", "production")
		Setup()
		assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
	})
}

func TestSetupWithOutput(t *testing.T) {
	var buf bytes.Buffer
	SetupWithOutput(&buf)

	log.Debug().Msg("debug message")
	log.Info().Msg("info message")

	output := buf.String()
	assert.Contains(t, output, "info message")
}

func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	assert.NotNil(t, logger)
}
