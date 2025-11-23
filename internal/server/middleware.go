package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupMiddleware(echoServer *echo.Echo) {
	echoServer.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${status} ${method} ${host} ${path} ${latency_human} ${bytes_in} ${bytes_out}\n",
	}))

	echoServer.Use(middleware.Recover())
}
