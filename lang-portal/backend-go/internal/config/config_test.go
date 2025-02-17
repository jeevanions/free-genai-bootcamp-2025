package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Clear environment before each test
	os.Clearenv()

	t.Run("default values when no env vars", func(t *testing.T) {
		cfg, err := Load()
		assert.NoError(t, err)
		assert.Equal(t, "development", cfg.Environment)
		assert.Equal(t, "./data/italian-learning.db", cfg.DatabasePath)
		assert.Equal(t, ":8080", cfg.ServerAddress)
		assert.Equal(t, "info", cfg.LogLevel)
	})

	t.Run("custom values from env vars", func(t *testing.T) {
		os.Setenv("APP_ENV", "staging")
		os.Setenv("DB_PATH", "./test.db")
		os.Setenv("SERVER_ADDR", ":3000")
		os.Setenv("LOG_LEVEL", "debug")
		defer os.Clearenv()

		cfg, err := Load()
		assert.NoError(t, err)
		assert.Equal(t, "staging", cfg.Environment)
		assert.Equal(t, "./test.db", cfg.DatabasePath)
		assert.Equal(t, ":3000", cfg.ServerAddress)
		assert.Equal(t, "debug", cfg.LogLevel)
	})
}

func TestLoad_InvalidEnvironment(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_ENV", "invalid")
	defer os.Clearenv()

	_, err := Load()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid environment")
}

func TestGetEnvOrDefault(t *testing.T) {
	t.Run("existing environment variable", func(t *testing.T) {
		os.Setenv("TEST_KEY", "test_value")
		defer os.Unsetenv("TEST_KEY")

		value := getEnvOrDefault("TEST_KEY", "default")
		assert.Equal(t, "test_value", value)
	})

	t.Run("missing environment variable", func(t *testing.T) {
		value := getEnvOrDefault("MISSING_KEY", "default")
		assert.Equal(t, "default", value)
	})
}
