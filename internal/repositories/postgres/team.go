package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type TeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool}
}
