package handlers

import (
	"autumn-2025/internal/models/dto"
	"autumn-2025/internal/repositories"
	"autumn-2025/internal/services"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type UsersHandler struct {
	usersService *services.UsersService
}

func NewUsersHandler(repository repositories.UserRepository) *UsersHandler {
	return &UsersHandler{usersService: services.NewUsersService(repository)}
}

//nolint:dupl
func (handler *UsersHandler) SetIsActive(ctx echo.Context) error {
	var setUserIsActive dto.SetUserIsActive

	err := ctx.Bind(&setUserIsActive)
	if err != nil {
		log.Error().Err(err).Msg("cannot parse Users.SetIsActive body")

		err = ctx.JSON(http.StatusBadRequest, err)
		if err != nil {
			return fmt.Errorf("failed to serialize StatusBadRequest: %w", err)
		}

		return nil
	}

	userWithTeam, err := handler.usersService.SetUserIsActive(ctx.Request().Context(), setUserIsActive)
	if err != nil {
		log.Error().Err(err).Fields(setUserIsActive).Msg("cannot merge pull request")

		var errorResponse dto.ResponseError

		statusCode := http.StatusBadRequest

		if errors.As(err, &errorResponse) {
			switch errorResponse.ErrorDesc.Code {
			case dto.NotFoundCode:
				statusCode = http.StatusNotFound
			default:
				statusCode = http.StatusInternalServerError
			}
		}

		err = ctx.JSON(statusCode, err)
		if err != nil {
			return fmt.Errorf("failed to serialize service error: %w", err)
		}

		return nil
	}

	err = ctx.JSON(http.StatusOK, map[string]dto.UserWithTeam{"user": userWithTeam})
	if err != nil {
		return fmt.Errorf("failed to serialize StatusOK: %w", err)
	}

	return nil
}

func (handler *UsersHandler) GetReview(ctx echo.Context) error {
	return nil
}
