package handlers

import (
	"autumn-2025/internal/models/dto"
	"autumn-2025/internal/repositories"
	"autumn-2025/internal/services"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
	var team dto.Team
	if err := c.Bind(&team); err != nil {
		log.Info().Err(err).Msg("StatusBadRequest")
		return c.JSON(http.StatusBadRequest, err)
	}

	team, err := handler.teamService.CreateTeam(c.Request().Context(), team)
	if err != nil {
		log.Info().Err(err).Msg("teamService.CreateTeam failed")
		if errors.Is(err, repositories.TeamAlreadyExistsError) {
			return c.JSON(http.StatusConflict, dto.TeamAlreadyExists(team.TeamName))
		}
		return c.JSON(http.StatusBadRequest, "Something went wrong")
	}

	return c.JSON(http.StatusCreated, map[string]dto.Team{"team": team})
}

func (handler *TeamHandler) Get(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "not implemented")
}
