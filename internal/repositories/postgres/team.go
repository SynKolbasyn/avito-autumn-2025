package postgres

import (
	"autumn-2025/internal/models/dto"
	"autumn-2025/internal/repositories"
	"context"
	"errors"
	"fmt"

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

func (t *TeamRepository) WithTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Error().Err(err).Msg("failed to rollback transaction")
		}
	}()

	if err := f(context.WithValue(ctx, repositories.TxContextKey, tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
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
	var id uuid.UUID
	err := row.Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.UUID{}, repositories.TeamAlreadyExistsError
	}
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (t *TeamRepository) InsertOrUpdateUsers(ctx context.Context, teamMembers []dto.TeamMember) ([]uuid.UUID, error) {
	var (
		valuesCount = 3
		query       = "INSERT INTO users (id, name, is_active) VALUES "
		params      = make([]interface{}, 0, len(teamMembers)*valuesCount)
		ids         = make([]uuid.UUID, 0, len(teamMembers))
	)
	for i, user := range teamMembers {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d)", valuesCount*i+1, valuesCount*i+2, valuesCount*i+3)
		params = append(params, user.UserID, user.Username, user.IsActive)
		ids = append(ids, user.UserID)
	}
	query += `
    ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		is_active = EXCLUDED.is_active,
		updated_at = CURRENT_TIMESTAMP;
    `
	executor := t.GetExecutor(ctx)
	_, err := executor.Exec(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	return ids, err
}

func (t *TeamRepository) AddTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDS []uuid.UUID) error {
	var (
		valuesCount = 2
		query       = "INSERT INTO user_teams (user_id, team_id) VALUES "
		params      = make([]interface{}, 0, len(memberIDS)*valuesCount)
	)
	for i, memberID := range memberIDS {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d)", valuesCount*i+1, valuesCount*i+2)
		params = append(params, memberID, teamID)
	}
	query += "ON CONFLICT (user_id, team_id) DO NOTHING;"
	executor := t.GetExecutor(ctx)
	_, err := executor.Exec(ctx, query, params...)
	return err
}
