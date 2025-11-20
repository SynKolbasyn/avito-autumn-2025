// Package server contains functions for work with echo web server framework.
package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupMiddleware configures common middleware for the Echo web server instance.
func SetupMiddleware(echoServer *echo.Echo) {
	echoServer.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${status} ${method} ${host} ${path} ${latency_human} ${bytes_in} ${bytes_out}\n",
	}))

	echoServer.Use(middleware.Recover())
}
