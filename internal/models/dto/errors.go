package dto

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	TeamAlreadyExistsCode = "TEAM_EXISTS"
	NotFoundCode          = "NOT_FOUND"
	PRExistsCode          = "PR_EXISTS"
	InternalErrorCode     = "INTERNAL_ERROR"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	ErrorDesc Error `json:"error"`
}

func (e ErrorResponse) Error() string {
	var str strings.Builder
	err := json.NewEncoder(&str).Encode(e)
	if err != nil {
		return fmt.Errorf("cannot encode error response: %w", err).Error()
	}
	return str.String()
}

func TeamAlreadyExists(teamName string) ErrorResponse {
	return ErrorResponse{
		ErrorDesc: Error{
			Code:    TeamAlreadyExistsCode,
			Message: teamName + "already exists",
		},
	}
}

func NotFound() ErrorResponse {
	return ErrorResponse{
		ErrorDesc: Error{
			Code:    NotFoundCode,
			Message: "resource not found",
		},
	}
}

func InternalError() ErrorResponse {
	return ErrorResponse{
		ErrorDesc: Error{
			Code:    InternalErrorCode,
			Message: "internal error",
		},
	}
}

func PullRequestExists(pullRequestID uuid.UUID) ErrorResponse {
	return ErrorResponse{
		ErrorDesc: Error{
			Code:    PRExistsCode,
			Message: fmt.Sprintf("PR %s already exists", pullRequestID),
		},
	}
}
