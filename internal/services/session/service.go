package session

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/sb-luis/creative-coding-bookclub/internal/model"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// Service handles session management
type Service struct {
	db *sql.DB
}

// NewService creates a new session service
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// CreateSession creates a new session for a member
func (s *Service) CreateSession(memberID int) (*model.Session, error) {
	if memberID <= 0 {
		return nil, errors.New("invalid member ID")
	}

	// Generate session ID
	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		log.Printf("Error generating session ID for member %d: %v", memberID, err)
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	expiresAt := time.Now().Add(24 * time.Hour) // 24 hours
	createdAt := time.Now()

	// Insert new session
	_, err = s.db.Exec(`
		INSERT INTO sessions (id, member_id, created_at, expires_at) 
		VALUES (?, ?, ?, ?)`,
		sessionID, memberID, createdAt, expiresAt)
	if err != nil {
		log.Printf("Database error while creating session for member %d: %v", memberID, err)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &model.Session{
		ID:        sessionID,
		MemberID:  memberID,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}

// GetSession returns a session by ID and validates it's not expired
func (s *Service) GetSession(sessionID string) (*model.Session, error) {
	if sessionID == "" {
		return nil, errors.New("session ID cannot be empty")
	}

	session := &model.Session{}
	err := s.db.QueryRow(`
		SELECT id, member_id, created_at, expires_at 
		FROM sessions WHERE id = ?`, sessionID).Scan(
		&session.ID, &session.MemberID, &session.CreatedAt, &session.ExpiresAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("session not found")
	}
	if err != nil {
		log.Printf("Database error while getting session '%s': %v", sessionID, err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		// Delete expired session
		s.DeleteSession(sessionID)
		return nil, errors.New("session expired")
	}

	return session, nil
}

// DeleteSession deletes a session
func (s *Service) DeleteSession(sessionID string) error {
	if sessionID == "" {
		return errors.New("session ID cannot be empty")
	}

	_, err := s.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		log.Printf("Database error while deleting session '%s': %v", sessionID, err)
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *Service) CleanupExpiredSessions() error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	if err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}
	return nil
}

// IsSessionValid checks if a session exists and is not expired
func (s *Service) IsSessionValid(sessionID string) bool {
	if sessionID == "" {
		return false
	}

	session, err := s.GetSession(sessionID)
	if err != nil {
		return false
	}

	return time.Now().Before(session.ExpiresAt)
}

// GetMemberIDFromSession retrieves the member ID associated with a session
func (s *Service) GetMemberIDFromSession(sessionID string) (int, error) {
	if sessionID == "" {
		return 0, errors.New("session ID cannot be empty")
	}

	session, err := s.GetSession(sessionID)
	if err != nil {
		return 0, err
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		return 0, errors.New("session expired")
	}

	return session.MemberID, nil
}
