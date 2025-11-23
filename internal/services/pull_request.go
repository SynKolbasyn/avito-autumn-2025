package services

import "autumn-2025/internal/repositories"

type PullRequestService struct {
	pullRequestRepository repositories.PullRequestRepository
}

func NewPullRequestService(repository repositories.PullRequestRepository) *PullRequestService {
	return &PullRequestService{
		pullRequestRepository: repository,
	}
}

func (pr *PullRequestService) CreatePullRequest() {

}

func (pr *PullRequestService) MergePullRequest() {

}

func (pr *PullRequestService) ReassignPullRequest() {

}
