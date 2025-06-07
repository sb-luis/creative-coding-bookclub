package model

import (
	"time"
)

// Session represents a member session
type Session struct {
	ID        string    `json:"id"`
	MemberID  int       `json:"member_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
