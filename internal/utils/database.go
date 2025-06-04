package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

// Member represents a member in the system
type Member struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Session represents a member session
type Session struct {
	ID        string    `json:"id"`
	MemberID  int       `json:"member_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Database represents our SQLite database
type Database struct {
	db *sql.DB
}

var dbInstance *Database

// InitDatabase initializes the SQLite database
func InitDatabase() error {
	var dbURL string
	var err error

	// Check if we're in production (Turso) or development (local SQLite)
	if tursoURL := os.Getenv("TURSO_DATABASE_URL"); tursoURL != "" {
		// Production: Use Turso
		tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
		if tursoToken == "" {
			return errors.New("TURSO_AUTH_TOKEN environment variable is required for production")
		}
		dbURL = fmt.Sprintf("%s?authToken=%s", tursoURL, tursoToken)
		log.Printf("Connecting to Turso database")
	} else {
		// Development: Use local SQLite
		if err := os.MkdirAll("data", 0755); err != nil {
			return err
		}
		dbURL = "file:data/bookclub.db"
		log.Printf("Using local SQLite database at data/bookclub.db")
	}

	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	dbInstance = &Database{db: db}

	// Create tables
	if err = dbInstance.createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Printf("Database initialized successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *Database {
	return dbInstance
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if dbInstance != nil && dbInstance.db != nil {
		return dbInstance.db.Close()
	}
	return nil
}

// createTables creates the necessary database tables
func (d *Database) createTables() error {
	// Members table
	membersTable := `
	CREATE TABLE IF NOT EXISTS members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := d.db.Exec(membersTable); err != nil {
		return fmt.Errorf("failed to create members table: %w", err)
	}

	// Sessions table
	sessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		member_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (member_id) REFERENCES members (id) ON DELETE CASCADE
	);`

	if _, err := d.db.Exec(sessionsTable); err != nil {
		return fmt.Errorf("failed to create sessions table: %w", err)
	}

	return nil
}

// CreateMember creates a new member
func (d *Database) CreateMember(name, passwordHash string) (*Member, error) {
	// Check if name already exists
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM members WHERE name = ?", name).Scan(&count)
	if err != nil {
		log.Printf("Database error while checking if member exists for name '%s': %v", name, err)
		return nil, fmt.Errorf("failed to check if member exists: %w", err)
	}
	if count > 0 {
		return nil, errors.New("name already exists")
	}

	// Insert new member
	result, err := d.db.Exec(`
		INSERT INTO members (name, password_hash, created_at, updated_at) 
		VALUES (?, ?, ?, ?)`,
		name, passwordHash, time.Now(), time.Now())
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
	return d.GetMemberByID(int(id))
}

// GetMemberByName returns a member by name
func (d *Database) GetMemberByName(name string) (*Member, error) {
	member := &Member{}
	err := d.db.QueryRow(`
		SELECT id, name, password_hash, created_at, updated_at 
		FROM members WHERE name = ?`, name).Scan(
		&member.ID, &member.Name, &member.PasswordHash,
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
func (d *Database) GetMemberByID(id int) (*Member, error) {
	member := &Member{}
	err := d.db.QueryRow(`
		SELECT id, name, password_hash, created_at, updated_at 
		FROM members WHERE id = ?`, id).Scan(
		&member.ID, &member.Name, &member.PasswordHash,
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

// CreateSession creates a new session for a member
func (d *Database) CreateSession(memberID int) (*Session, error) {
	// Generate session ID
	sessionID, err := generateSessionID()
	if err != nil {
		log.Printf("Error generating session ID for member %d: %v", memberID, err)
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	expiresAt := time.Now().Add(24 * time.Hour) // 24 hours
	createdAt := time.Now()

	// Insert new session
	_, err = d.db.Exec(`
		INSERT INTO sessions (id, member_id, created_at, expires_at) 
		VALUES (?, ?, ?, ?)`,
		sessionID, memberID, createdAt, expiresAt)
	if err != nil {
		log.Printf("Database error while creating session for member %d: %v", memberID, err)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &Session{
		ID:        sessionID,
		MemberID:  memberID,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}

// GetSession returns a session by ID
func (d *Database) GetSession(sessionID string) (*Session, error) {
	session := &Session{}
	err := d.db.QueryRow(`
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
		d.DeleteSession(sessionID)
		return nil, errors.New("session expired")
	}

	return session, nil
}

// DeleteSession deletes a session
func (d *Database) DeleteSession(sessionID string) error {
	_, err := d.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		log.Printf("Database error while deleting session '%s': %v", sessionID, err)
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// CleanupExpiredSessions removes expired sessions from the database
func (d *Database) CleanupExpiredSessions() error {
	_, err := d.db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	if err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}
	return nil
}

// generateSessionID generates a random session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
