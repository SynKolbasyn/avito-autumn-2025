package postgres

import (
	"autumn-2025/internal/repositories"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Executor struct {
	pool *pgxpool.Pool
}

func NewExecutor(pool *pgxpool.Pool) *Executor {
	return &Executor{pool: pool}
}

func (e *Executor) WithTransaction(ctx context.Context, function func(ctx context.Context) error) error {
	transaction, err := e.pool.Begin(ctx)
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

func (e *Executor) GetExecutor(ctx context.Context) repositories.Executor {
	if tx, ok := ctx.Value(repositories.TxContextKey).(pgx.Tx); ok {
		return tx
	}

	return e.pool
}
