package handlers

import (
	"autumn-2025/internal/models/dto"
	"autumn-2025/internal/repositories"
	"autumn-2025/internal/services"
	"fmt"
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

func (handler *TeamHandler) Add(ctx echo.Context) error {
	var team dto.Team

	err := ctx.Bind(&team)
	if err != nil {
		log.Info().Err(err).Msg("StatusBadRequest")

		err = ctx.JSON(http.StatusBadRequest, err)
		if err != nil {
			return fmt.Errorf("failed to serialize validation error: %w", err)
		}

		return nil
	}

	team, err = handler.teamService.CreateTeam(ctx.Request().Context(), team)
	if err != nil {
		log.Info().Err(err).Fields(team).Msg("teamService.CreateTeam failed")

		err = ctx.JSON(http.StatusBadRequest, dto.TeamAlreadyExists(team.TeamName))
		if err != nil {
			return fmt.Errorf("failed to serialize TeamAlreadyExists: %w", err)
		}

		return nil
	}

	err = ctx.JSON(http.StatusCreated, map[string]dto.Team{"team": team})
	if err != nil {
		return fmt.Errorf("failed to serialize team: %w", err)
	}

	return nil
}

func (handler *TeamHandler) Get(ctx echo.Context) error {
	teamName := ctx.QueryParam("team_name")

	team, err := handler.teamService.GetTeam(ctx.Request().Context(), teamName)
	if err != nil {
		log.Error().Err(err).Str("team_name", teamName).Msg("teamService.GetTeam failed")

		err = ctx.JSON(http.StatusNotFound, dto.NotFound())
		if err != nil {
			return fmt.Errorf("failed to serialize dto.NotFound: %w", err)
		}

		return nil
	}

	err = ctx.JSON(http.StatusOK, team)
	if err != nil {
		return fmt.Errorf("failed to serialize team: %w", err)
	}

	return nil
}
