package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sb-luis/creative-coding-bookclub/internal/model"
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

	// Sketches table
	sketchesTable := `
	CREATE TABLE IF NOT EXISTS sketches (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		member_id INTEGER NOT NULL,
		slug TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
		keywords TEXT DEFAULT '',
		tags TEXT DEFAULT '[]',
		external_libs TEXT DEFAULT '[]',
		source_code TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (member_id) REFERENCES members (id) ON DELETE CASCADE,
		UNIQUE(member_id, slug)
	);`

	if _, err := d.db.Exec(sketchesTable); err != nil {
		return fmt.Errorf("failed to create sketches table: %w", err)
	}

	// Index for better query performance
	sketchesIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_sketches_member_id ON sketches(member_id);",
		"CREATE INDEX IF NOT EXISTS idx_sketches_slug ON sketches(slug);",
		"CREATE INDEX IF NOT EXISTS idx_sketches_created_at ON sketches(created_at);",
	}

	for _, indexSQL := range sketchesIndexes {
		if _, err := d.db.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create sketches index: %w", err)
		}
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

// GetAllMembers returns all members (ID and Name only)
func (d *Database) GetAllMembers() ([]Member, error) {
	rows, err := d.db.Query(`
		SELECT id, name, created_at, updated_at 
		FROM members ORDER BY name ASC`)
	if err != nil {
		log.Printf("Database error while getting all members: %v", err)
		return nil, fmt.Errorf("failed to get all members: %w", err)
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		member := Member{}
		err := rows.Scan(
			&member.ID, &member.Name, &member.CreatedAt, &member.UpdatedAt)
		if err != nil {
			log.Printf("Database error while scanning member: %v", err)
			continue
		}
		members = append(members, member)
	}

	return members, nil
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

// CreateSketch creates a new sketch for a member
func (d *Database) CreateSketch(memberID int, req *model.CreateSketchRequest) (*model.Sketch, error) {
	// Generate slug from title
	slug := GenerateSlug(req.Title)
	if slug == "" {
		return nil, errors.New("title cannot be empty or invalid")
	}

	// Check if sketch slug already exists for this member
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM sketches WHERE member_id = ? AND slug = ?", memberID, slug).Scan(&count)
	if err != nil {
		log.Printf("Database error while checking if sketch slug exists for member %d, slug '%s': %v", memberID, slug, err)
		return nil, fmt.Errorf("failed to check if sketch slug exists: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("a sketch with the title '%s' already exists for this member (slug conflict: %s)", req.Title, slug)
	}

	// Marshal tags and external_libs to JSON
	tagsJSON, err := json.Marshal(req.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tags: %w", err)
	}
	externalLibsJSON, err := json.Marshal(req.ExternalLibs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal external libs: %w", err)
	}

	now := time.Now()
	result, err := d.db.Exec(`
		INSERT INTO sketches (member_id, slug, title, description, keywords, tags, external_libs, source_code, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		memberID, slug, req.Title, req.Description, req.Keywords, string(tagsJSON), string(externalLibsJSON), req.SourceCode, now, now)
	if err != nil {
		log.Printf("Database error while creating sketch for member %d: %v", memberID, err)
		return nil, fmt.Errorf("failed to create sketch: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Database error while getting sketch ID for member %d: %v", memberID, err)
		return nil, fmt.Errorf("failed to get sketch ID: %w", err)
	}

	return d.GetSketchByID(int(id))
}

// GetSketchByID returns a sketch by ID
func (d *Database) GetSketchByID(id int) (*model.Sketch, error) {
	sketch := &model.Sketch{}
	err := d.db.QueryRow(`
		SELECT id, member_id, slug, title, description, keywords, tags, external_libs, source_code, created_at, updated_at 
		FROM sketches WHERE id = ?`, id).Scan(
		&sketch.ID, &sketch.MemberID, &sketch.Slug, &sketch.Title, &sketch.Description,
		&sketch.Keywords, &sketch.TagsJSON, &sketch.ExternalLibsJSON, &sketch.SourceCode,
		&sketch.CreatedAt, &sketch.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("sketch not found")
	}
	if err != nil {
		log.Printf("Database error while getting sketch by ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get sketch by ID: %w", err)
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal([]byte(sketch.TagsJSON), &sketch.Tags); err != nil {
		log.Printf("Warning: failed to unmarshal tags for sketch %d: %v", id, err)
		sketch.Tags = []string{} // fallback to empty slice
	}
	if err := json.Unmarshal([]byte(sketch.ExternalLibsJSON), &sketch.ExternalLibs); err != nil {
		log.Printf("Warning: failed to unmarshal external libs for sketch %d: %v", id, err)
		sketch.ExternalLibs = []string{} // fallback to empty slice
	}

	return sketch, nil
}

// GetSketchByMemberAndSlug returns a sketch by member ID and slug
func (d *Database) GetSketchByMemberAndSlug(memberID int, slug string) (*model.Sketch, error) {
	sketch := &model.Sketch{}
	err := d.db.QueryRow(`
		SELECT id, member_id, slug, title, description, keywords, tags, external_libs, source_code, created_at, updated_at 
		FROM sketches WHERE member_id = ? AND slug = ?`, memberID, slug).Scan(
		&sketch.ID, &sketch.MemberID, &sketch.Slug, &sketch.Title, &sketch.Description,
		&sketch.Keywords, &sketch.TagsJSON, &sketch.ExternalLibsJSON, &sketch.SourceCode,
		&sketch.CreatedAt, &sketch.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("sketch not found")
	}
	if err != nil {
		log.Printf("Database error while getting sketch by member %d and slug '%s': %v", memberID, slug, err)
		return nil, fmt.Errorf("failed to get sketch by member and slug: %w", err)
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal([]byte(sketch.TagsJSON), &sketch.Tags); err != nil {
		log.Printf("Warning: failed to unmarshal tags for sketch %d: %v", sketch.ID, err)
		sketch.Tags = []string{} // fallback to empty slice
	}
	if err := json.Unmarshal([]byte(sketch.ExternalLibsJSON), &sketch.ExternalLibs); err != nil {
		log.Printf("Warning: failed to unmarshal external libs for sketch %d: %v", sketch.ID, err)
		sketch.ExternalLibs = []string{} // fallback to empty slice
	}

	return sketch, nil
}

// GetSketchesByMember returns all sketches for a member
func (d *Database) GetSketchesByMember(memberID int) ([]*model.Sketch, error) {
	rows, err := d.db.Query(`
		SELECT id, member_id, slug, title, description, keywords, tags, external_libs, source_code, created_at, updated_at 
		FROM sketches WHERE member_id = ? ORDER BY created_at DESC`, memberID)
	if err != nil {
		log.Printf("Database error while getting sketches for member %d: %v", memberID, err)
		return nil, fmt.Errorf("failed to get sketches by member: %w", err)
	}
	defer rows.Close()

	var sketches []*model.Sketch
	for rows.Next() {
		sketch := &model.Sketch{}
		err := rows.Scan(
			&sketch.ID, &sketch.MemberID, &sketch.Slug, &sketch.Title, &sketch.Description,
			&sketch.Keywords, &sketch.TagsJSON, &sketch.ExternalLibsJSON, &sketch.SourceCode,
			&sketch.CreatedAt, &sketch.UpdatedAt)
		if err != nil {
			log.Printf("Database error while scanning sketch for member %d: %v", memberID, err)
			continue
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal([]byte(sketch.TagsJSON), &sketch.Tags); err != nil {
			log.Printf("Warning: failed to unmarshal tags for sketch %d: %v", sketch.ID, err)
			sketch.Tags = []string{} // fallback to empty slice
		}
		if err := json.Unmarshal([]byte(sketch.ExternalLibsJSON), &sketch.ExternalLibs); err != nil {
			log.Printf("Warning: failed to unmarshal external libs for sketch %d: %v", sketch.ID, err)
			sketch.ExternalLibs = []string{} // fallback to empty slice
		}

		sketches = append(sketches, sketch)
	}

	return sketches, nil
}

// GetAllSketches returns all sketches from all members
func (d *Database) GetAllSketches() ([]*model.Sketch, error) {
	rows, err := d.db.Query(`
		SELECT id, member_id, slug, title, description, keywords, tags, external_libs, source_code, created_at, updated_at 
		FROM sketches ORDER BY created_at DESC`)
	if err != nil {
		log.Printf("Database error while getting all sketches: %v", err)
		return nil, fmt.Errorf("failed to get all sketches: %w", err)
	}
	defer rows.Close()

	var sketches []*model.Sketch
	for rows.Next() {
		sketch := &model.Sketch{}
		err := rows.Scan(
			&sketch.ID, &sketch.MemberID, &sketch.Slug, &sketch.Title, &sketch.Description,
			&sketch.Keywords, &sketch.TagsJSON, &sketch.ExternalLibsJSON, &sketch.SourceCode,
			&sketch.CreatedAt, &sketch.UpdatedAt)
		if err != nil {
			log.Printf("Database error while scanning sketch: %v", err)
			continue
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal([]byte(sketch.TagsJSON), &sketch.Tags); err != nil {
			log.Printf("Warning: failed to unmarshal tags for sketch %d: %v", sketch.ID, err)
			sketch.Tags = []string{} // fallback to empty slice
		}
		if err := json.Unmarshal([]byte(sketch.ExternalLibsJSON), &sketch.ExternalLibs); err != nil {
			log.Printf("Warning: failed to unmarshal external libs for sketch %d: %v", sketch.ID, err)
			sketch.ExternalLibs = []string{} // fallback to empty slice
		}

		sketches = append(sketches, sketch)
	}

	return sketches, nil
}

// GetAllSketchesGroupedByMember returns all sketches grouped by member name
func (d *Database) GetAllSketchesGroupedByMember() ([]model.MemberSketchInfo, error) {
	rows, err := d.db.Query(`
		SELECT s.id, s.member_id, s.slug, s.title, s.description, s.keywords, s.tags, s.external_libs, s.source_code, s.created_at, s.updated_at, m.name as member_name
		FROM sketches s
		JOIN members m ON s.member_id = m.id
		ORDER BY m.name, s.created_at DESC`)
	if err != nil {
		log.Printf("Database error while getting all sketches grouped by member: %v", err)
		return nil, fmt.Errorf("failed to get all sketches grouped by member: %w", err)
	}
	defer rows.Close()

	memberSketchMap := make(map[string][]model.SketchInfo)

	for rows.Next() {
		var sketch model.Sketch
		var memberName string
		err := rows.Scan(
			&sketch.ID, &sketch.MemberID, &sketch.Slug, &sketch.Title, &sketch.Description,
			&sketch.Keywords, &sketch.TagsJSON, &sketch.ExternalLibsJSON, &sketch.SourceCode,
			&sketch.CreatedAt, &sketch.UpdatedAt, &memberName)
		if err != nil {
			log.Printf("Database error while scanning sketch: %v", err)
			continue
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal([]byte(sketch.TagsJSON), &sketch.Tags); err != nil {
			log.Printf("Warning: failed to unmarshal tags for sketch %d: %v", sketch.ID, err)
			sketch.Tags = []string{}
		}

		// Convert database sketch to SketchInfo for the lister
		sketchInfo := model.SketchInfo{
			Slug:  sketch.Slug,
			URL:   fmt.Sprintf("/members/%s/%s", memberName, sketch.Slug),
			Alias: memberName,
		}

		// Set pointers for optional fields
		if sketch.Title != "" {
			sketchInfo.Title = &sketch.Title
		}
		if sketch.Description != "" {
			sketchInfo.Description = &sketch.Description
		}
		if sketch.Keywords != "" {
			sketchInfo.Keywords = &sketch.Keywords
		}
		if len(sketch.Tags) > 0 {
			sketchInfo.Tags = sketch.Tags
		}

		memberSketchMap[memberName] = append(memberSketchMap[memberName], sketchInfo)
	}

	// Convert map to slice
	var result []model.MemberSketchInfo
	for memberName, sketches := range memberSketchMap {
		if len(sketches) > 0 {
			result = append(result, model.MemberSketchInfo{
				Name:     memberName,
				Sketches: sketches,
			})
		}
	}

	return result, nil
}

// UpdateSketch updates an existing sketch
func (d *Database) UpdateSketch(id int, req *model.UpdateSketchRequest) (*model.Sketch, error) {
	// Start building the update query
	setParts := []string{}
	args := []interface{}{}

	// Get current sketch to check if title is changing
	currentSketch, err := d.GetSketchByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get current sketch: %w", err)
	}

	if req.Title != nil {
		// If title is changing, generate new slug and check for conflicts
		newSlug := GenerateSlug(*req.Title)
		if newSlug == "" {
			return nil, errors.New("title cannot be empty or invalid")
		}

		// Only check for slug conflicts if the slug is actually changing
		if newSlug != currentSketch.Slug {
			var count int
			err := d.db.QueryRow("SELECT COUNT(*) FROM sketches WHERE member_id = ? AND slug = ? AND id != ?",
				currentSketch.MemberID, newSlug, id).Scan(&count)
			if err != nil {
				log.Printf("Database error while checking slug conflict for sketch %d: %v", id, err)
				return nil, fmt.Errorf("failed to check slug conflict: %w", err)
			}
			if count > 0 {
				return nil, fmt.Errorf("a sketch with the title '%s' already exists for this member (slug conflict: %s)", *req.Title, newSlug)
			}

			setParts = append(setParts, "slug = ?")
			args = append(args, newSlug)
		}

		setParts = append(setParts, "title = ?")
		args = append(args, *req.Title)
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
	}
	if req.Keywords != nil {
		setParts = append(setParts, "keywords = ?")
		args = append(args, *req.Keywords)
	}
	if req.Tags != nil {
		tagsJSON, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags: %w", err)
		}
		setParts = append(setParts, "tags = ?")
		args = append(args, string(tagsJSON))
	}
	if req.ExternalLibs != nil {
		externalLibsJSON, err := json.Marshal(req.ExternalLibs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal external libs: %w", err)
		}
		setParts = append(setParts, "external_libs = ?")
		args = append(args, string(externalLibsJSON))
	}
	if req.SourceCode != nil {
		setParts = append(setParts, "source_code = ?")
		args = append(args, *req.SourceCode)
	}

	if len(setParts) == 0 {
		return nil, errors.New("no fields to update")
	}

	// Always update the updated_at field
	setParts = append(setParts, "updated_at = ?")
	args = append(args, time.Now())

	// Add the ID for the WHERE clause
	args = append(args, id)

	query := fmt.Sprintf("UPDATE sketches SET %s WHERE id = ?", setParts[0])
	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf("%s, %s", query, setParts[i])
	}

	_, err = d.db.Exec(query, args...)
	if err != nil {
		log.Printf("Database error while updating sketch %d: %v", id, err)
		return nil, fmt.Errorf("failed to update sketch: %w", err)
	}

	return d.GetSketchByID(id)
}

// DeleteSketch deletes a sketch by ID
func (d *Database) DeleteSketch(id int) error {
	result, err := d.db.Exec("DELETE FROM sketches WHERE id = ?", id)
	if err != nil {
		log.Printf("Database error while deleting sketch %d: %v", id, err)
		return fmt.Errorf("failed to delete sketch: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("sketch not found")
	}

	return nil
}

// DeleteSketchByMemberAndSlug deletes a sketch by member ID and slug
func (d *Database) DeleteSketchByMemberAndSlug(memberID int, slug string) error {
	result, err := d.db.Exec("DELETE FROM sketches WHERE member_id = ? AND slug = ?", memberID, slug)
	if err != nil {
		log.Printf("Database error while deleting sketch for member %d, slug '%s': %v", memberID, slug, err)
		return fmt.Errorf("failed to delete sketch: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("sketch not found")
	}

	return nil
}
