package postgres

import (
	"autumn-2025/internal/models/dto"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository struct {
	*Executor
}

func NewUsersRepository(pool *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{NewExecutor(pool)}
}

func (u *UsersRepository) SetUserIsActive(ctx context.Context, userID uuid.UUID, isActive bool) (dto.TeamMember, bool) {
	query := `
	UPDATE users
	SET is_active = $1,
	    updated_at = CURRENT_TIMESTAMP
	WHERE id = $2
	RETURNING id, name, is_active;`
	executor := u.GetExecutor(ctx)
	row := executor.QueryRow(ctx, query, isActive, userID)

	var (
		uID       uuid.UUID
		uName     string
		uIsActive bool
	)

	err := row.Scan(&uID, &uName, &uIsActive)
	if err != nil {
		return dto.TeamMember{}, false
	}

	return dto.TeamMember{
		UserID:   uID,
		Username: uName,
		IsActive: uIsActive,
	}, true
}

func (u *UsersRepository) GetUserTeam(ctx context.Context, userID uuid.UUID) (string, error) {
	query := `
	SELECT name
	FROM teams AS t
	INNER JOIN user_teams AS u ON (u.user_id = $1) AND (u.team_id = t.id);
	`
	executor := u.GetExecutor(ctx)
	row := executor.QueryRow(ctx, query, userID)
	teamName := ""

	err := row.Scan(&teamName)
	if err != nil {
		return teamName, fmt.Errorf("error getting team name: %w", err)
	}

	return teamName, nil
}
