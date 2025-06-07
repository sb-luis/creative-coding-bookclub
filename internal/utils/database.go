package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
)

var (
	db *sql.DB
)

// InitDatabase initializes the database connection
func InitDatabase() error {
	tursoURL := os.Getenv("TURSO_DATABASE_URL")
	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")

	var dbConn *sql.DB
	var err error

	if tursoURL != "" && tursoToken != "" {
		// Production: Use Turso with embedded replica for better performance
		log.Printf("Connecting to Turso database with embedded replica")

		// Create data directory for local replica
		if err := os.MkdirAll("data", 0755); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}

		dbPath := filepath.Join("data", "replica.db")

		connector, err := libsql.NewEmbeddedReplicaConnector(
			dbPath,
			tursoURL,
			libsql.WithAuthToken(tursoToken),
		)
		if err != nil {
			return fmt.Errorf("failed to create Turso connector: %w", err)
		}

		dbConn = sql.OpenDB(connector)

		// Sync with remote database
		frames, syncErr := connector.Sync()
		if syncErr != nil {
			log.Printf("Warning: failed to sync with remote database: %v", syncErr)
		} else {
			log.Printf("Synced %d frames from remote database", frames)
		}
	} else {
		// Development: Use local SQLite only
		log.Printf("Using local SQLite database")

		if err := os.MkdirAll("data", 0755); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}

		dbPath := filepath.Join("data", "bookclub.db")

		// For local-only development, use the standard sql.Open with file: scheme
		dbConn, err = sql.Open("libsql", "file:"+dbPath)
		if err != nil {
			return fmt.Errorf("failed to open local database: %w", err)
		}

		log.Printf("Local database setup complete (no remote sync)")
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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		verified BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(membersTable); err != nil {
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

	if _, err := db.Exec(sessionsTable); err != nil {
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
