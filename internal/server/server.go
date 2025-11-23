// Package server contains functions for work with echo web server framework.
package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const readTimeout = 10 * time.Second
const writeTimeout = 10 * time.Second
const idleTimeout = 120 * time.Second
const shutdownTimeout = 10 * time.Second

// Start starts the Echo HTTP server on the port defined in the SERVER_PORT
// environment variable. It configures the underlying http.Server with read, write,
// and idle timeouts to improve server robustness.
//
// The server runs in a separate goroutine so that the function can listen for OS
// termination signals (SIGINT, SIGTERM). Upon receiving such a signal, Start
// initiates a graceful shutdown with a 10-second timeout to allow ongoing requests
// to complete before closing.
func Start(server *echo.Echo, address string) {
	httpServer := &http.Server{
		Addr:         address,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	go func() {
		log.Info().Msgf("Starting server on %s", address)

		err := server.StartServer(httpServer)

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Server shutdown error")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Server shutdown error")
	}

	log.Info().Msg("Server stopped")
}
