package model

import (
	"time"
)

// Member represents a member in the system
type Member struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"password_hash"`
	Verified     bool      `json:"verified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateMemberRequest represents the data needed to create a new member
type CreateMemberRequest struct {
	Name         string `json:"name" validate:"required,min=1,max=50"`
	PasswordHash string `json:"password_hash" validate:"required"`
}
