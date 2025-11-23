package dto

const (
	teamAlreadyExistsCode = "TEAM_EXISTS"
	notFoundCode          = "NOT_FOUND"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

func TeamAlreadyExists(teamName string) ErrorResponse {
	return ErrorResponse{
		Error: Error{
			Code:    teamAlreadyExistsCode,
			Message: teamName + "already exists",
		},
	}
}

func NotFound() ErrorResponse {
	return ErrorResponse{
		Error: Error{
			Code:    notFoundCode,
			Message: "resource not found",
		},
	}
}
