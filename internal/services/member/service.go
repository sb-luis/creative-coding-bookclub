package member

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/sb-luis/creative-coding-bookclub/internal/model"
)

// Service handles member-related business logic
type Service struct {
	db *sql.DB
}

// NewService creates a new member service
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// CreateMember creates a new member with validation (password should be pre-hashed)
func (s *Service) CreateMember(name, passwordHash string) (*model.Member, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if passwordHash == "" {
		return nil, errors.New("password hash cannot be empty")
	}

	// Check if name already exists
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM members WHERE name = ?", name).Scan(&count)
	if err != nil {
		log.Printf("Database error while checking if member exists for name '%s': %v", name, err)
		return nil, fmt.Errorf("failed to check if member exists: %w", err)
	}
	if count > 0 {
		return nil, errors.New("name already exists")
	}

	// Insert new member
	result, err := s.db.Exec(`
		INSERT INTO members (name, password_hash, verified, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)`,
		name, passwordHash, false, time.Now(), time.Now())
	if err != nil {
		log.Printf("Database error while creating member '%s': %v", name, err)
		return nil, fmt.Errorf("failed to create member: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Database error while getting member ID for '%s': %v", name, err)
		return nil, fmt.Errorf("failed to get member ID: %w", err)
	}

	// Return the created member
	return s.GetMemberByID(int(id))
}

// GetMemberByName returns a member by name
func (s *Service) GetMemberByName(name string) (*model.Member, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	member := &model.Member{}
	err := s.db.QueryRow(`
		SELECT id, name, password_hash, verified, created_at, updated_at 
		FROM members WHERE name = ?`, name).Scan(
		&member.ID, &member.Name, &member.PasswordHash, &member.Verified,
		&member.CreatedAt, &member.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("member not found")
	}
	if err != nil {
		log.Printf("Database error while getting member by name '%s': %v", name, err)
		return nil, fmt.Errorf("failed to get member by name: %w", err)
	}

	return member, nil
}

// GetMemberByID returns a member by ID
func (s *Service) GetMemberByID(id int) (*model.Member, error) {
	if id <= 0 {
		return nil, errors.New("invalid member ID")
	}

	member := &model.Member{}
	err := s.db.QueryRow(`
		SELECT id, name, password_hash, verified, created_at, updated_at 
		FROM members WHERE id = ?`, id).Scan(
		&member.ID, &member.Name, &member.PasswordHash, &member.Verified,
		&member.CreatedAt, &member.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("member not found")
	}
	if err != nil {
		log.Printf("Database error while getting member by ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get member by ID: %w", err)
	}

	return member, nil
}

// GetAllMembers returns all members
func (s *Service) GetAllMembers() ([]model.Member, error) {
	rows, err := s.db.Query(`
		SELECT id, name, verified, created_at, updated_at 
		FROM members ORDER BY name ASC`)
	if err != nil {
		log.Printf("Database error while getting all members: %v", err)
		return nil, fmt.Errorf("failed to get all members: %w", err)
	}
	defer rows.Close()

	var members []model.Member
	for rows.Next() {
		member := model.Member{}
		err := rows.Scan(
			&member.ID, &member.Name, &member.Verified, &member.CreatedAt, &member.UpdatedAt)
		if err != nil {
			log.Printf("Database error while scanning member: %v", err)
			continue
		}
		members = append(members, member)
	}

	return members, nil
}

// UpdatePassword updates a member's password (only for verified members)
func (s *Service) UpdatePassword(memberID int, currentPasswordHash, newPasswordHash string) error {
	if memberID <= 0 {
		return errors.New("invalid member ID")
	}
	if currentPasswordHash == "" {
		return errors.New("current password hash cannot be empty")
	}
	if newPasswordHash == "" {
		return errors.New("new password hash cannot be empty")
	}

	// Get the current member to verify they exist and are verified
	member, err := s.GetMemberByID(memberID)
	if err != nil {
		return fmt.Errorf("failed to get member: %w", err)
	}

	// Check if member is verified
	if !member.Verified {
		return errors.New("only verified members can update their password")
	}

	// Verify current password
	if member.PasswordHash != currentPasswordHash {
		return errors.New("current password is incorrect")
	}

	// Update the password
	_, err = s.db.Exec(`
		UPDATE members 
		SET password_hash = ?, updated_at = ? 
		WHERE id = ?`,
		newPasswordHash, time.Now(), memberID)

	if err != nil {
		log.Printf("Database error while updating password for member ID %d: %v", memberID, err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
