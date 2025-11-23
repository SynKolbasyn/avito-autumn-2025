package services

import (
	"autumn-2025/internal/models/dto"
	"autumn-2025/internal/repositories"
	"context"
)

type TeamService struct {
	teamRepository repositories.TeamRepository
}

func NewTeamService(repository repositories.TeamRepository) *TeamService {
	return &TeamService{
		teamRepository: repository,
	}
}

func (t *TeamService) CreateTeam(ctx context.Context, team dto.Team) (dto.Team, error) {
	err := t.teamRepository.WithTransaction(ctx, func(c context.Context) error {
		id, err := t.teamRepository.CreateTeam(c, team.TeamName)
		if err != nil {
			return err
		}

		ids, err := t.teamRepository.InsertOrUpdateUsers(c, team.Members)
		if err != nil {
			return err
		}

		return t.teamRepository.AddTeamMembers(c, id, ids)
	})

	if err != nil {
		return dto.Team{}, err
	}

	return team, nil
}
