package services

import (
	"autumn-2025/internal/models/dto"
	"autumn-2025/internal/repositories"
	"context"
	"fmt"
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
	err := t.teamRepository.WithTransaction(ctx, func(txCtx context.Context) error {
		teamID, err := t.teamRepository.CreateTeam(txCtx, team.TeamName)
		if err != nil {
			return fmt.Errorf("cannot create team: %w", err)
		}

		userIDs, err := t.teamRepository.InsertOrUpdateUsers(txCtx, team.Members)
		if err != nil {
			return fmt.Errorf("cannot insert or update users: %w", err)
		}

		err = t.teamRepository.AddTeamMembers(txCtx, teamID, userIDs)
		if err != nil {
			return fmt.Errorf("cannot add team members: %w", err)
		}

		return nil
	})
	if err != nil {
		return dto.Team{}, fmt.Errorf("cannot add team members: %w", err)
	}

	return team, nil
}
