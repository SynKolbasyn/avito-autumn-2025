package repositories

import (
	"autumn-2025/internal/models/dto"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type repositoryContextKey string

const TxContextKey repositoryContextKey = "pgx_tx"

type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TransactionManager interface {
	WithTransaction(ctx context.Context, function func(ctx context.Context) error) error
	GetExecutor(ctx context.Context) Executor
}

type TeamRepository interface {
	TransactionManager
	CreateTeam(ctx context.Context, teamName string) (uuid.UUID, error)
	InsertOrUpdateUsers(ctx context.Context, teamMembers []dto.TeamMember) ([]uuid.UUID, error)
	AddTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDS []uuid.UUID) error
	GetTeamByName(ctx context.Context, teamName string) (dto.Team, error)
}
