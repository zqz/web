package domain

import (
	"time"
)

// User represents a user in the system (domain model)
type User struct {
	ID                  int32
	Name                string
	Email               string
	Provider            string
	ProviderID          string
	Role                string
	DisplayTag          string // 1-3 char tag for display; empty if not set
	Colour              string // hex e.g. #RRGGBB; empty if not set
	Banned              bool
	MaxFileSizeOverride *int64 // bytes; nil = use site default. Admin-set only.
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// IsAdmin returns true if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsMember returns true if the user is a regular member
func (u *User) IsMember() bool {
	return u.Role == "member"
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Name       string
	Email      string
	Provider   string
	ProviderID string
	Role       string
}
