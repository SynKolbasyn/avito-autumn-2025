package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type PullRequestRepository struct {
	*Executor
}

func NewPullRequestRepository(pool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{NewExecutor(pool)}
}
