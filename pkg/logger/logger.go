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
func InitLogger(logLevel zerolog.Level) {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	zerolog.SetGlobalLevel(logLevel)
}
