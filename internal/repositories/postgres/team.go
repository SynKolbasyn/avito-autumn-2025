package postgres

import (
	"autumn-2025/internal/models/dto"
	"autumn-2025/internal/repositories"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type TeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool}
}

func (t *TeamRepository) WithTransaction(ctx context.Context, function func(ctx context.Context) error) error {
	transaction, err := t.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("cannot start transaction: %w", err)
	}

	defer func() {
		rbErr := transaction.Rollback(ctx)
		if rbErr != nil {
			log.Error().Err(rbErr).Msg("failed to rollback transaction")
		}
	}()

	err = function(context.WithValue(ctx, repositories.TxContextKey, transaction))
	if err != nil {
		return fmt.Errorf("transaction error: %w", err)
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return fmt.Errorf("cannot commit transaction: %w", err)
	}

	return nil
}

func (t *TeamRepository) GetExecutor(ctx context.Context) repositories.Executor {
	if tx, ok := ctx.Value(repositories.TxContextKey).(pgx.Tx); ok {
		return tx
	}

	return t.pool
}

func (t *TeamRepository) CreateTeam(ctx context.Context, teamName string) (uuid.UUID, error) {
	query := `
	INSERT INTO teams (name) VALUES ($1)
	ON CONFLICT (name) DO NOTHING
	RETURNING id;
	`
	executor := t.GetExecutor(ctx)
	row := executor.QueryRow(ctx, query, teamName)

	var teamID uuid.UUID

	err := row.Scan(&teamID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert team: %w", err)
	}

	return teamID, nil
}

func (t *TeamRepository) InsertOrUpdateUsers(ctx context.Context, teamMembers []dto.TeamMember) ([]uuid.UUID, error) {
	if len(teamMembers) == 0 {
		return []uuid.UUID{}, nil
	}

	var (
		valuesCount = 3
		query       strings.Builder
		params      = make([]any, 0, len(teamMembers)*valuesCount)
		ids         = make([]uuid.UUID, 0, len(teamMembers))
	)
	query.WriteString("INSERT INTO users (id, name, is_active) VALUES ")

	for i, user := range teamMembers {
		if i > 0 {
			query.WriteString(", ")
		}

		query.WriteString(fmt.Sprintf("($%d, $%d, $%d)", valuesCount*i+1, valuesCount*i+2, valuesCount*i+3))

		params = append(params, user.UserID, user.Username, user.IsActive)
		ids = append(ids, user.UserID)
	}

	query.WriteString(`
    ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		is_active = EXCLUDED.is_active,
		updated_at = CURRENT_TIMESTAMP;
    `)

	executor := t.GetExecutor(ctx)

	_, err := executor.Exec(ctx, query.String(), params...)
	if err != nil {
		return nil, fmt.Errorf("failed to update users: %w", err)
	}

	return ids, nil
}

func (t *TeamRepository) AddTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDS []uuid.UUID) error {
	if len(memberIDS) == 0 {
		return nil
	}

	var (
		valuesCount = 2
		query       strings.Builder
		params      = make([]any, 0, len(memberIDS)*valuesCount)
	)
	query.WriteString("INSERT INTO user_teams (user_id, team_id) VALUES ")

	for i, memberID := range memberIDS {
		if i > 0 {
			query.WriteString(", ")
		}

		query.WriteString(fmt.Sprintf("($%d, $%d)", valuesCount*i+1, valuesCount*i+2))

		params = append(params, memberID, teamID)
	}

	query.WriteString("ON CONFLICT (user_id, team_id) DO NOTHING;")

	executor := t.GetExecutor(ctx)

	_, err := executor.Exec(ctx, query.String(), params...)
	if err != nil {
		return fmt.Errorf("failed to add team members to users: %w", err)
	}

	return nil
}

func (t *TeamRepository) GetTeamByName(ctx context.Context, teamName string) (dto.Team, error) {
	query := `
	SELECT u.id, u.name, u.is_active
	FROM teams AS t
	LEFT JOIN user_teams AS ut ON ut.team_id = t.id
	LEFT JOIN users AS u ON u.id = ut.user_id
	WHERE t.name = $1
	`
	executor := t.GetExecutor(ctx)

	rows, err := executor.Query(ctx, query, teamName)
	if err != nil {
		return dto.Team{}, fmt.Errorf("failed to get query team by name: %w", err)
	}
	defer rows.Close()

	var teamMembers []dto.TeamMember

	for rows.Next() {
		var (
			userID   uuid.UUID
			name     string
			isActive bool
		)

		err = rows.Scan(&userID, &name, &isActive)
		if err != nil {
			continue
		}

		teamMembers = append(teamMembers, dto.TeamMember{
			UserID:   userID,
			Username: name,
			IsActive: isActive,
		})
	}

	if rows.CommandTag().RowsAffected() == 0 {
		return dto.Team{}, pgx.ErrNoRows
	}

	if len(teamMembers) == 0 {
		teamMembers = []dto.TeamMember{}
	}

	return dto.Team{
		TeamName: teamName,
		Members:  teamMembers,
	}, nil
}
