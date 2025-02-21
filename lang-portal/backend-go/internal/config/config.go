package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Port    int
	DBPath  string
	EnvMode string
}

// Load returns a Config struct populated with values from environment variables
func Load() *Config {
	port, _ := strconv.Atoi(getEnvOrDefault("PORT", "8080"))
	
	return &Config{
		Port:    port,
		DBPath:  getEnvOrDefault("DB_PATH", "words.db"),
		EnvMode: getEnvOrDefault("ENV_MODE", "development"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
