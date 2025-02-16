package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *Config
		wantErr        bool
	}{
		{
			name:    "default values when no env vars",
			envVars: map[string]string{},
			expectedConfig: &Config{
				ServerAddress: ":8080",
				DBPath:        "./words.db",
				Environment:   "development",
			},
			wantErr: false,
		},
		{
			name: "custom values from env vars",
			envVars: map[string]string{
				"SERVER_ADDRESS": ":3000",
				"DB_PATH":        "./custom.db",
				"ENVIRONMENT":    "production",
			},
			expectedConfig: &Config{
				ServerAddress: ":3000",
				DBPath:        "./custom.db",
				Environment:   "production",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment before each test
			os.Clearenv()

			// Set environment variables for test
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Load configuration
			cfg, err := Load()

			// Assert results
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedConfig.ServerAddress, cfg.ServerAddress)
			assert.Equal(t, tt.expectedConfig.DBPath, cfg.DBPath)
			assert.Equal(t, tt.expectedConfig.Environment, cfg.Environment)
		})
	}
}

func TestLoad_InvalidEnvironment(t *testing.T) {
	// Clear environment
	os.Clearenv()

	// Set invalid server address
	os.Setenv("SERVER_ADDRESS", "invalid:port")

	cfg, err := Load()
	assert.NoError(t, err) // Should not error as we don't validate format yet
	assert.Equal(t, "invalid:port", cfg.ServerAddress)
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "existing environment variable",
			key:          "TEST_KEY",
			envValue:     "test_value",
			defaultValue: "default",
			expected:     "test_value",
		},
		{
			name:         "missing environment variable",
			key:          "MISSING_KEY",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := getEnvOrDefault(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
