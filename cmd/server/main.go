// Package main starts the web service.
package main

import (
	"autumn-2025/internal/server"
	"autumn-2025/pkg/logger"

	"github.com/labstack/echo/v4"
)

func main() {
	logger.InitLogger()

	echoServer := echo.New()
	echoServer.HideBanner = true
	echoServer.HidePort = true

	server.SetupMiddleware(echoServer)
	server.RegisterRoutes(echoServer)
	server.StartServer(echoServer)
}


