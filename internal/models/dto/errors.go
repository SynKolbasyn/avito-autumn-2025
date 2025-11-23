package dto

const (
	teamAlreadyExistsCode = "TEAM_EXISTS"
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
