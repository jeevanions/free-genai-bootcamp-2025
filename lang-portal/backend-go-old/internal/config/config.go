package config

import (
	"fmt"
	"os"
)

type Config struct {
	Environment   string
	DatabasePath  string
	ServerAddress string
	LogLevel      string
}

func Load() (*Config, error) {
	cfg := &Config{
		Environment:   getEnvOrDefault("APP_ENV", "development"),
		DatabasePath:  getEnvOrDefault("DB_PATH", "./data/italian-learning.db"),
		ServerAddress: getEnvOrDefault("SERVER_ADDR", ":8080"),
		LogLevel:      getEnvOrDefault("LOG_LEVEL", "info"),
	}

	// Validate environment
	switch cfg.Environment {
	case "development", "staging", "production":
		// valid
	default:
		return nil, fmt.Errorf("invalid environment: %s", cfg.Environment)
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
