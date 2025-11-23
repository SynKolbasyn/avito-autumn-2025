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

type ResponseError struct {
	ErrorDesc Error `json:"error"`
}

func (e ResponseError) Error() string {
	var str strings.Builder

	err := json.NewEncoder(&str).Encode(e)
	if err != nil {
		return fmt.Errorf("cannot encode error response: %w", err).Error()
	}

	return str.String()
}

func TeamAlreadyExists(teamName string) ResponseError {
	return ResponseError{
		ErrorDesc: Error{
			Code:    TeamAlreadyExistsCode,
			Message: teamName + "already exists",
		},
	}
}

func NotFound() ResponseError {
	return ResponseError{
		ErrorDesc: Error{
			Code:    NotFoundCode,
			Message: "resource not found",
		},
	}
}

func InternalError() ResponseError {
	return ResponseError{
		ErrorDesc: Error{
			Code:    InternalErrorCode,
			Message: "internal error",
		},
	}
}

func PullRequestExists(pullRequestID uuid.UUID) ResponseError {
	return ResponseError{
		ErrorDesc: Error{
			Code:    PRExistsCode,
			Message: fmt.Sprintf("PR %s already exists", pullRequestID),
		},
	}
}
