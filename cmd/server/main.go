// Package main starts the web service.
package main

import (
	"autumn-2025/internal/server"
	"autumn-2025/pkg/logger"

	"github.com/labstack/echo/v4"
)

func main() {
	logger.InitLogger()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	server.SetupMiddleware(e)
	server.RegisterRoutes(e)
	server.StartServer(e)
}


