package config

import (
	"os"
)

type Config struct {
	ServerAddress string
	DBPath        string
	Environment   string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerAddress: getEnvOrDefault("SERVER_ADDRESS", ":8080"),
		DBPath:        getEnvOrDefault("DB_PATH", "./words.db"),
		Environment:   getEnvOrDefault("ENVIRONMENT", "development"),
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
