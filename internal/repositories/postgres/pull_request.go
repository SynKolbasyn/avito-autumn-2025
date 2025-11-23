package postgres

import (
	"autumn-2025/internal/models/dto"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PullRequestRepository struct {
	*Executor
}

func NewPullRequestRepository(pool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{NewExecutor(pool)}
}

func (r *PullRequestRepository) Exists(ctx context.Context, pullRequestID uuid.UUID) bool {
	query := "SELECT EXISTS(SELECT 1 FROM pull_requests WHERE id = $1);"
	executor := r.GetExecutor(ctx)
	row := executor.QueryRow(ctx, query, pullRequestID)

	var exists bool

	err := row.Scan(&exists)
	if err != nil {
		log.Error().Err(err).Str("pull_request_id", pullRequestID.String()).Msg("error checking existence of pull request")

		return false
	}

	return exists
}

func (r *PullRequestRepository) CreatePR(ctx context.Context, prID uuid.UUID, prName string, prAuthor uuid.UUID) bool {
	query := `
	INSERT INTO pull_requests (id, title, author_id) VALUES ($1, $2, $3)
	ON CONFLICT DO NOTHING
	RETURNING TRUE;
	`
	executor := r.GetExecutor(ctx)
	row := executor.QueryRow(ctx, query, prID, prName, prAuthor)

	ok := false
	if row.Scan(&ok) != nil {
		return false
	}

	return ok
}

func (r *PullRequestRepository) GetTeamMembersIDs(ctx context.Context, prAuthorID uuid.UUID) ([]uuid.UUID, error) {
	query := `
	SELECT ut.user_id
	FROM user_teams AS ut
	WHERE (ut.team_id = (SELECT u.team_id FROM user_teams AS u WHERE u.user_id = $1)) AND ut.user_id != $1;
	`
	executor := r.GetExecutor(ctx)

	rows, err := executor.Query(ctx, query, prAuthorID)
	if err != nil {
		return nil, fmt.Errorf("error getting team members IDs: %w", err)
	}
	defer rows.Close()

	var members []uuid.UUID

	for rows.Next() {
		var memberID uuid.UUID

		err = rows.Scan(&memberID)
		if err != nil {
			return nil, fmt.Errorf("error getting team member ID: %w", err)
		}

		members = append(members, memberID)
	}

	return members, nil
}

func (r *PullRequestRepository) AssignMembers(
	ctx context.Context,
	pullRequestID uuid.UUID,
	teamMembersIDs []uuid.UUID,
) error {
	var (
		argsCount = 2
		args      = make([]any, 0, argsCount*len(teamMembersIDs))
		query     strings.Builder
	)
	query.WriteString("INSERT INTO reviewers (pr_id, user_id) VALUES ")

	for i, memberID := range teamMembersIDs {
		if i > 0 {
			query.WriteString(", ")
		}

		query.WriteString(fmt.Sprintf("($%d, $%d)", i*argsCount+1, i*argsCount+2))

		args = append(args, pullRequestID, memberID)
	}

	executor := r.GetExecutor(ctx)

	_, err := executor.Exec(ctx, query.String(), args...)
	if err != nil {
		log.Error().
			Err(err).
			Str("pull_request_id", pullRequestID.String()).
			Fields(teamMembersIDs).
			Msg("error assigning members to reviewers")

		return fmt.Errorf("error assigning members to reviewers: %w", err)
	}

	return nil
}

func (r *PullRequestRepository) Merged(ctx context.Context, pullRequestID uuid.UUID) (dto.PullRequestMerged, bool) {
	query := "SELECT id, title, author_id, status, merged_at FROM pull_requests WHERE id = $1 AND status = 'MERGED';"
	executor := r.GetExecutor(ctx)
	row := executor.QueryRow(ctx, query, pullRequestID)
	var (
		prID       uuid.UUID
		prName     string
		prAuthorID uuid.UUID
		prStatus   dto.PullRequestStatus
		prMergedAt time.Time
	)
	err := row.Scan(&prID, &prName, &prAuthorID, &prStatus, &prMergedAt)
	if err != nil {
		return dto.PullRequestMerged{}, false
	}
	return dto.PullRequestMerged{
		PullRequestCreated: dto.PullRequestCreated{
			PullRequestCreate: dto.PullRequestCreate{
				PullRequestID:   prID,
				PullRequestName: prName,
				AuthorID:        prAuthorID,
			},
			Status:            prStatus,
			AssignedReviewers: nil,
		},
		MergedAt: prMergedAt,
	}, true
}

func (r *PullRequestRepository) Merge(ctx context.Context, pullRequestID uuid.UUID) (dto.PullRequestMerged, error) {
	var (
		prID       uuid.UUID
		prName     string
		prAuthorID uuid.UUID
		prStatus   dto.PullRequestStatus
		prMergedAt time.Time
	)
	query := `
	UPDATE pull_requests
	SET status = $1,
		merged_at = CURRENT_TIMESTAMP
	WHERE id = $2
	RETURNING id, title, author_id, status, merged_at;
	`
	executor := r.GetExecutor(ctx)
	row := executor.QueryRow(ctx, query, dto.PullRequestStatusMerged, pullRequestID)
	err := row.Scan(&prID, &prName, &prAuthorID, &prStatus, &prMergedAt)
	if err != nil {
		return dto.PullRequestMerged{}, fmt.Errorf("error merging pull request: %w", err)
	}
	return dto.PullRequestMerged{
		PullRequestCreated: dto.PullRequestCreated{
			PullRequestCreate: dto.PullRequestCreate{
				PullRequestID:   prID,
				PullRequestName: prName,
				AuthorID:        prAuthorID,
			},
			Status:            prStatus,
			AssignedReviewers: nil,
		},
		MergedAt: prMergedAt,
	}, nil
}

func (r *PullRequestRepository) GetReviewersIDs(ctx context.Context, pullRequestID uuid.UUID) ([]uuid.UUID, error) {
	query := "SELECT user_id FROM reviewers WHERE pr_id = $1;"
	executor := r.GetExecutor(ctx)
	rows, err := executor.Query(ctx, query, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("error getting reviewers IDs: %w", err)
	}
	defer rows.Close()
	var reviewers []uuid.UUID
	for rows.Next() {
		var reviewerID uuid.UUID
		err = rows.Scan(&reviewerID)
		if err != nil {
			return nil, fmt.Errorf("error getting reviewer ID: %w", err)
		}
		reviewers = append(reviewers, reviewerID)
	}
	if len(reviewers) == 0 {
		reviewers = []uuid.UUID{}
	}
	return reviewers, nil
}
