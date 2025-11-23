// Package main starts the web service.
package main

import (
	"autumn-2025/config"
	"autumn-2025/internal/server"
	"autumn-2025/pkg/database"
	"autumn-2025/pkg/logger"
	"context"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger.InitLogger(cfg.LogLevel)

	pool := database.Pool(context.Background(), cfg)
	defer pool.Close()

	echoServer := echo.New()
	echoServer.HideBanner = true
	echoServer.HidePort = true

	server.SetupMiddleware(echoServer)
	server.RegisterRoutes(echoServer, context.Background(), pool)
	server.Start(echoServer, cfg.ServerAddress())
}
