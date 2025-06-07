package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	db *sql.DB
)

// InitDatabase initializes the database connection
func InitDatabase() error {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	log.Printf("Connecting to PostgreSQL database")

	dbConn, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err = dbConn.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db = dbConn

	// Create tables
	if err = createTables(dbConn); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Printf("Database initialized successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return db
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if db != nil {
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
		log.Printf("Database connection closed")
	}
	return nil
}

// createTables creates the necessary database tables
func createTables(db *sql.DB) error {
	// Members table
	membersTable := `
	CREATE TABLE IF NOT EXISTS members (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		verified BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(membersTable); err != nil {
		return fmt.Errorf("failed to create members table: %w", err)
	}

	// Sessions table
	sessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		member_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NOT NULL,
		FOREIGN KEY (member_id) REFERENCES members (id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(sessionsTable); err != nil {
		return fmt.Errorf("failed to create sessions table: %w", err)
	}

	// Sketches table
	sketchesTable := `
	CREATE TABLE IF NOT EXISTS sketches (
		id SERIAL PRIMARY KEY,
		member_id INTEGER NOT NULL,
		slug TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT DEFAULT '',
		keywords TEXT DEFAULT '',
		tags TEXT DEFAULT '[]',
		external_libs TEXT DEFAULT '[]',
		source_code TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (member_id) REFERENCES members (id) ON DELETE CASCADE,
		UNIQUE(member_id, slug)
	);`

	if _, err := db.Exec(sketchesTable); err != nil {
		return fmt.Errorf("failed to create sketches table: %w", err)
	}

	// Index for better query performance
	sketchesIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_sketches_member_id ON sketches(member_id);",
		"CREATE INDEX IF NOT EXISTS idx_sketches_slug ON sketches(slug);",
		"CREATE INDEX IF NOT EXISTS idx_sketches_created_at ON sketches(created_at);",
	}

	for _, indexSQL := range sketchesIndexes {
		if _, err := db.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create sketches index: %w", err)
		}
	}

	return nil
}
