package services

import (
	"database/sql"

	"github.com/sb-luis/creative-coding-bookclub/internal/services/member"
	"github.com/sb-luis/creative-coding-bookclub/internal/services/session"
	"github.com/sb-luis/creative-coding-bookclub/internal/services/sketch"
)

// Services contains all application services
type Services struct {
	Member  *member.Service
	Session *session.Service
	Sketch  *sketch.Service
}

// NewServices creates a new services container with all services initialized
func NewServices(db *sql.DB) *Services {
	memberService := member.NewService(db)
	return &Services{
		Member:  memberService,
		Session: session.NewService(db),
		Sketch:  sketch.NewService(db),
	}
}
