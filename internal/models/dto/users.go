package dto

import "github.com/google/uuid"

type TeamMember struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	IsActive bool      `json:"is_active"`
}

type SetUserIsActive struct {
	UserID   uuid.UUID `json:"user_id"`
	IsActive bool      `json:"is_active"`
}

type UserWithTeam struct {
	TeamMember
	TeamName string `json:"team_name"`
}
