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

type PullRequestHandler struct {
	pullRequestService *services.PullRequestService
}

func NewPullRequestHandler(repository repositories.PullRequestRepository) *PullRequestHandler {
	return &PullRequestHandler{
		pullRequestService: services.NewPullRequestService(repository),
	}
}

func (handler *PullRequestHandler) Create(ctx echo.Context) error {
	var pullRequestCreate dto.PullRequestCreate

	err := ctx.Bind(&pullRequestCreate)
	if err != nil {
		log.Error().Err(err).Msg("cannot parse PullRequest.Create body")

		err = ctx.JSON(http.StatusBadRequest, err)
		if err != nil {
			return fmt.Errorf("failed to serialize StatusBadRequest: %w", err)
		}

		return nil
	}

	pullRequestCreated, err := handler.pullRequestService.CreatePullRequest(ctx.Request().Context(), pullRequestCreate)
	if err != nil {
		log.Error().Err(err).Fields(pullRequestCreate).Msg("cannot create pull request")

		var errorResponse dto.ResponseError

		statusCode := http.StatusBadRequest

		if errors.As(err, &errorResponse) {
			switch errorResponse.ErrorDesc.Code {
			case dto.PRExistsCode:
				statusCode = http.StatusConflict
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

	err = ctx.JSON(http.StatusCreated, map[string]dto.PullRequestCreated{"pr": pullRequestCreated})
	if err != nil {
		return fmt.Errorf("failed to serialize StatusCreated: %w", err)
	}

	return nil
}

//nolint:dupl
func (handler *PullRequestHandler) Merge(ctx echo.Context) error {
	var pullRequestMerge dto.PullRequestMerge

	err := ctx.Bind(&pullRequestMerge)
	if err != nil {
		log.Error().Err(err).Msg("cannot parse PullRequest.Merge body")

		err = ctx.JSON(http.StatusBadRequest, err)
		if err != nil {
			return fmt.Errorf("failed to serialize StatusBadRequest: %w", err)
		}

		return nil
	}

	pullRequestMerged, err := handler.pullRequestService.MergePullRequest(ctx.Request().Context(), pullRequestMerge)
	if err != nil {
		log.Error().Err(err).Fields(pullRequestMerge).Msg("cannot merge pull request")

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

	err = ctx.JSON(http.StatusOK, map[string]dto.PullRequestMerged{"pr": pullRequestMerged})
	if err != nil {
		return fmt.Errorf("failed to serialize StatusOK: %w", err)
	}

	return nil
}

func (handler *PullRequestHandler) Reassign(ctx echo.Context) error {
	err := ctx.JSON(http.StatusNotImplemented, dto.InternalError())
	if err != nil {
		return fmt.Errorf("failed to serialize service error: %w", err)
	}

	return nil
}
