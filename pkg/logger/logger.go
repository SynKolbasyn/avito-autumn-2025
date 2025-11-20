// Package logger provides initialization and configuration of the global logger
// using the zerolog library.
package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the global logger using zerolog.
// The time format for log entries is set to RFC3339 for readability and standardization.
// Logs are output to the console (standard output) using a human-friendly format.
// The logging level is set dynamically based on the ENV environment variable:
// - In "production" mode, the log level is set to InfoLevel (info and above).
// - In other environments, it is set to DebugLevel (verbose logging for development).
func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	env := os.Getenv("ENV")
	if env == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
