package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Setup() {
	// Set default level to info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Configure logger output
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// Create logger
	log = zerolog.New(output).With().Timestamp().Logger()

	// Set log level based on environment
	env := os.Getenv("APP_ENV")
	if env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func SetupWithOutput(w io.Writer) {
	output := zerolog.ConsoleWriter{
		Out:        w,
		TimeFormat: time.RFC3339,
	}
	log = zerolog.New(output).With().Timestamp().Logger()
}

func GetLogger() *zerolog.Logger {
	return &log
}
