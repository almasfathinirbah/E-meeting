package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the zerolog logging configuration for the application
// It sets up:
//   - Time format using RFC3339
//   - Global log level to Info
//   - Pretty console logging for development environment
//   - Caller information in log entries
func InitLogger() {
	// Set the time format for log entries to RFC3339
	zerolog.TimeFieldFormat = time.RFC3339

	// Set the default global log level to Info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Enable pretty console logging with colors in development environment
	if os.Getenv("ENV") == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	// Add caller information (file and line number) to log entries
	log.Logger = log.With().Caller().Logger()
}
