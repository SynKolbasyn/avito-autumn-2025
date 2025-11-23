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

type PullRequestRepository interface {
	TransactionManager
	CreatePR(ctx context.Context, prID uuid.UUID, prName string, prAuthor uuid.UUID) bool
	GetActiveTeamMembersIDs(ctx context.Context, prAuthorID uuid.UUID) ([]uuid.UUID, error)
	AssignMembers(ctx context.Context, pullRequestID uuid.UUID, teamMembersIDs []uuid.UUID) error
	Merged(ctx context.Context, pullRequestID uuid.UUID) (dto.PullRequestMerged, bool)
	Merge(ctx context.Context, pullRequestID uuid.UUID) (dto.PullRequestMerged, error)
	GetReviewersIDs(ctx context.Context, pullRequestID uuid.UUID) ([]uuid.UUID, error)
}

type UserRepository interface {
	TransactionManager
	SetUserIsActive(ctx context.Context, userID uuid.UUID, isActive bool) (dto.TeamMember, bool)
	GetUserTeam(ctx context.Context, userID uuid.UUID) (string, error)
}
