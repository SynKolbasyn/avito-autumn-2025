package services

import (
	"autumn-2025/internal/models/dto"
	"autumn-2025/internal/repositories"
	"context"
	"math/rand/v2"

	"github.com/google/uuid"
)

type PullRequestService struct {
	pullRequestRepository repositories.PullRequestRepository
}

func NewPullRequestService(repository repositories.PullRequestRepository) *PullRequestService {
	return &PullRequestService{
		pullRequestRepository: repository,
	}
}

func (pr *PullRequestService) CreatePullRequest(
	ctx context.Context,
	pullRequest dto.PullRequestCreate,
) (dto.PullRequestCreated, error) {
	var assignedReviewers []uuid.UUID

	err := pr.pullRequestRepository.WithTransaction(ctx, func(txCtx context.Context) error {
		prExists := pr.pullRequestRepository.CreatePR(
			txCtx,
			pullRequest.PullRequestID,
			pullRequest.PullRequestName,
			pullRequest.AuthorID,
		)

		if !prExists {
			return dto.PullRequestExists(pullRequest.PullRequestID)
		}

		teamMembersIDs, err := pr.pullRequestRepository.GetTeamMembersIDs(txCtx, pullRequest.AuthorID)
		if err != nil || len(teamMembersIDs) == 0 {
			return dto.NotFound()
		}

		if len(teamMembersIDs) > 1 {
			rand.Shuffle(len(teamMembersIDs), func(i, j int) {
				teamMembersIDs[i], teamMembersIDs[j] = teamMembersIDs[j], teamMembersIDs[i]
			})
			teamMembersIDs = teamMembersIDs[:2]
		}

		err = pr.pullRequestRepository.AssignMembers(txCtx, pullRequest.PullRequestID, teamMembersIDs)
		if err != nil {
			return dto.InternalError()
		}

		assignedReviewers = teamMembersIDs

		return nil
	})
	if err != nil {
		//nolint:wrapcheck
		return dto.PullRequestCreated{}, err
	}

	return dto.PullRequestCreated{
		PullRequestCreate: pullRequest,
		Status:            dto.PullRequestStatusOpen,
		AssignedReviewers: assignedReviewers,
	}, nil
}

func (pr *PullRequestService) MergePullRequest(
	ctx context.Context,
	prMerge dto.PullRequestMerge,
) (dto.PullRequestMerged, error) {
	var mergedPullRequest dto.PullRequestMerged

	err := pr.pullRequestRepository.WithTransaction(ctx, func(txCtx context.Context) error {
		existingPR, exists := pr.pullRequestRepository.Merged(txCtx, prMerge.PullRequestID)
		if !exists {
			mergedPr, err := pr.pullRequestRepository.Merge(txCtx, prMerge.PullRequestID)
			if err != nil {
				return dto.NotFound()
			}

			existingPR = mergedPr
		}

		prReviewers, err := pr.pullRequestRepository.GetReviewersIDs(txCtx, prMerge.PullRequestID)
		if err != nil {
			return dto.NotFound()
		}

		existingPR.AssignedReviewers = prReviewers
		mergedPullRequest = existingPR

		return nil
	})
	if err != nil {
		//nolint:wrapcheck
		return dto.PullRequestMerged{}, err
	}

	return mergedPullRequest, nil
}

func (pr *PullRequestService) ReassignPullRequest() {

}
