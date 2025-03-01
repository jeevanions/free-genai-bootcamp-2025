package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *Config
		wantErr  bool
	}{
		{
			name: "all environment variables set",
			envVars: map[string]string{
				"DB_PATH": "./test.db",
				"PORT":    "8080",
				"ENV":     "development",
			},
			expected: &Config{
				DatabasePath:  "./test.db",
				ServerAddress: ":8080",
				Environment:   "development",
			},
			wantErr: false,
		},
		// Add more test cases...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				// Clean up environment variables
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			cfg, err := Load()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, cfg)
		})
	}
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
