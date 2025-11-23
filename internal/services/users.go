package services

import (
	"autumn-2025/internal/models/dto"
	"autumn-2025/internal/repositories"
	"context"
)

type UsersService struct {
	userRepository repositories.UserRepository
}

func NewUsersService(repository repositories.UserRepository) *UsersService {
	return &UsersService{repository}
}

func (u *UsersService) SetUserIsActive(ctx context.Context, user dto.SetUserIsActive) (dto.UserWithTeam, error) {
	var userWithTeam dto.UserWithTeam
	err := u.userRepository.WithTransaction(ctx, func(txCtx context.Context) error {
		teamMember, updated := u.userRepository.SetUserIsActive(txCtx, user.UserID, user.IsActive)
		if !updated {
			return dto.NotFound()
		}
		teamName, err := u.userRepository.GetUserTeam(txCtx, user.UserID)
		if err != nil {
			return err
		}
		userWithTeam.TeamMember = teamMember
		userWithTeam.TeamName = teamName
		return nil
	})
	if err != nil {
		return dto.UserWithTeam{}, err
	}
	return userWithTeam, nil
}
