package dto

import "github.com/google/uuid"

type PullRequestStatus string

const (
	PullRequestStatusOpen   PullRequestStatus = "OPEN"
	PullRequestStatusClosed PullRequestStatus = "CLOSED"
)

type PullRequestCreate struct {
	PullRequestID   uuid.UUID `json:"pull_request_id"`
	PullRequestName string    `json:"pull_request_name"`
	AuthorID        uuid.UUID `json:"author_id"`
}

type PullRequestCreated struct {
	PullRequestCreate

	Status            PullRequestStatus `json:"status"`
	AssignedReviewers []uuid.UUID       `json:"assigned_reviewers"`
}
