package postgres

import (
	"autumn-2025/internal/repositories"
	"context"

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

func (t *TeamRepository) CreateTeam(ctx context.Context, teamName string) error {
	query := `
	INSERT INTO teams (name) VALUES ($1)
	ON CONFLICT (name) DO NOTHING;
	`
	executor := t.GetExecutor(ctx)
	tag, err := executor.Exec(ctx, query, teamName)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return repositories.TeamAlreadyExistsError
	}
	return nil
}
