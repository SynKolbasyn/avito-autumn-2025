package handlers

import (
	"autumn-2025/internal/repositories"
	"autumn-2025/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(repository repositories.TeamRepository) *TeamHandler {
	return &TeamHandler{
		teamService: services.NewTeamService(repository),
	}
}

func (handler *TeamHandler) Add(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "not implemented")
}

func (handler *TeamHandler) Get(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "not implemented")
}
