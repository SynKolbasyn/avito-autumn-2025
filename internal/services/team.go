package services

import "autumn-2025/internal/repositories"

type TeamService struct {
	teamRepository repositories.TeamRepository
}

func NewTeamService(repository repositories.TeamRepository) *TeamService {
	return &TeamService{
		teamRepository: repository,
	}
}
