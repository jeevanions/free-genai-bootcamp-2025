package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(writer ...io.Writer) {
	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set global log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if os.Getenv("ENVIRONMENT") == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Configure logger output
	var output io.Writer = os.Stdout
	if len(writer) > 0 && writer[0] != nil {
		output = writer[0]
	}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}
