package handlers

import (
	"autumn-2025/internal/repositories"
	"autumn-2025/internal/services"

	"github.com/labstack/echo/v4"
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

	return nil
}

func (handler *PullRequestHandler) Merge(ctx echo.Context) error {
	return nil
}

func (handler *PullRequestHandler) Reassign(ctx echo.Context) error {
	return nil
}
