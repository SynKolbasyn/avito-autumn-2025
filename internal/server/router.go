// Package server contains functions for work with echo web server framework.
package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers HTTP routes on the provided Echo instance.
func RegisterRoutes(echoServer *echo.Echo) {
	echoServer.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
}
