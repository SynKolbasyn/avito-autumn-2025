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

	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
}
