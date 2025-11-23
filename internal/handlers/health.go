package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
