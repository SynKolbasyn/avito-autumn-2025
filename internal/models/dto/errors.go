package dto

const (
	teamAlreadyExistsCode = "TEAM_EXISTS"
	notFoundCode          = "NOT_FOUND"
)

type Error struct {
	code    string
	message string
}

type ErrorResponse struct {
	error Error
}

func TeamAlreadyExists(teamName string) ErrorResponse {
	return ErrorResponse{
		error: Error{
			code:    teamAlreadyExistsCode,
			message: teamName + "already exists",
		},
	}
}

func NotFound() ErrorResponse {
	return ErrorResponse{
		error: Error{
			code:    notFoundCode,
			message: "resource not found",
		},
	}
}
