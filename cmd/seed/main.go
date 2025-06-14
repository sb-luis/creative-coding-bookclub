package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sb-luis/creative-coding-bookclub/internal/model"
	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// SketchMetadata represents the structure of sketch JSON files
type SketchMetadata struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Keywords     string   `json:"keywords"`
	Tags         []string `json:"tags"`
	ExternalLibs []string `json:"external_libs"`
	CreatedAt    *string  `json:"created_at,omitempty"`
	UpdatedAt    *string  `json:"updated_at,omitempty"`
}

// MemberConfig represents the structure of members in members.json
type MemberConfig struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func main() {
	// Configure logger to write to stdout
	log.SetOutput(os.Stdout)

	// Load .env file during development only
	// In production, use environment variables that are already set
	appEnv := os.Getenv("APP_ENV")
	if appEnv != "production" {
		if err := utils.LoadEnvFile(); err != nil {
			log.Printf("Note: Could not load .env file (this is normal in production): %v", err)
		}
	}

	// Initialize database
	err := utils.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer utils.CloseDatabase()

	// Initialize services
	globalServices := services.NewServices(utils.GetDB())

	// Load members configuration from JSON file
	membersConfig, err := loadMembersConfig("data/seed/members.json")
	if err != nil {
		log.Fatalf("Failed to load members configuration: %v", err)
	}

	sketchesDir := "data/seed/sketches"

	// Process each member in the order specified in members.json
	for _, memberConfig := range membersConfig {
		memberName := memberConfig.Name
		log.Printf("Processing member: %s", memberName)

		// Check if member directory exists
		memberDir := filepath.Join(sketchesDir, memberName)
		if _, err := os.Stat(memberDir); os.IsNotExist(err) {
			log.Printf("Skipping member %s: directory %s does not exist", memberName, memberDir)
			continue
		}

		// Create member with the password from configuration
		err := createMember(globalServices, memberName, memberConfig.Password)
		if err != nil {
			log.Printf("Failed to create member %s: %v", memberName, err)
			continue
		}

		// Process sketches for this member
		err = processMemberSketches(globalServices, memberName, memberDir)
		if err != nil {
			log.Printf("Failed to process sketches for member %s: %v", memberName, err)
		}
	}

	// Verify all members manually via SQL (for seeding purposes)
	err = verifyAllMembers(utils.GetDB())
	if err != nil {
		log.Printf("Failed to verify all members: %v", err)
	}

	log.Println("Seeding completed!")
}

func loadMembersConfig(configPath string) ([]MemberConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read members config file: %w", err)
	}

	var members []MemberConfig
	err = json.Unmarshal(data, &members)
	if err != nil {
		return nil, fmt.Errorf("failed to parse members config JSON: %w", err)
	}

	log.Printf("Loaded %d members from configuration file", len(members))
	return members, nil
}

func createMember(services *services.Services, name, password string) error {
	if services == nil {
		return errors.New("services not initialized")
	}

	// Hash the password
	passwordHash := utils.HashPassword(password)

	// Create the member using the regular CreateMember method
	member, err := services.Member.CreateMember(name, passwordHash)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("Member %s already exists, skipping creation", name)
			return nil
		}
		return fmt.Errorf("failed to create member: %w", err)
	}

	log.Printf("Created member: %s (ID: %d)", member.Name, member.ID)
	return nil
}

// verifyAllMembers manually sets all members as verified via direct SQL
func verifyAllMembers(db *sql.DB) error {
	if db == nil {
		return errors.New("database connection is nil")
	}

	// Update all members to be verified
	result, err := db.Exec("UPDATE members SET verified = true WHERE verified = false")
	if err != nil {
		return fmt.Errorf("failed to verify members: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Could not get rows affected count: %v", err)
	} else {
		log.Printf("Verified %d members via direct SQL", rowsAffected)
	}

	return nil
}

func processMemberSketches(services *services.Services, memberName, memberDir string) error {
	if services == nil {
		return errors.New("services not initialized")
	}

	// Get member from database
	member, err := services.Member.GetMemberByName(memberName)
	if err != nil {
		return fmt.Errorf("failed to get member %s: %w", memberName, err)
	}

	// Find all .js files in the member directory
	err = filepath.WalkDir(memberDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Only process .js files
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".js") {
			sketchName := strings.TrimSuffix(d.Name(), ".js")
			log.Printf("  Processing sketch: %s", sketchName)

			err := createSketch(services, member.ID, memberDir, sketchName)
			if err != nil {
				log.Printf("    Failed to create sketch %s: %v", sketchName, err)
			} else {
				log.Printf("    Created sketch: %s", sketchName)
			}
		}

		return nil
	})

	return err
}

func createSketch(services *services.Services, memberID int, memberDir, sketchName string) error {
	if services == nil {
		return errors.New("services not initialized")
	}

	// Read JavaScript source code
	jsPath := filepath.Join(memberDir, sketchName+".js")
	sourceCode, err := os.ReadFile(jsPath)
	if err != nil {
		return fmt.Errorf("failed to read JS file %s: %w", jsPath, err)
	}

	// Read JSON metadata (optional)
	jsonPath := filepath.Join(memberDir, sketchName+".json")
	var metadata SketchMetadata

	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Printf("    No JSON metadata for sketch %s, using defaults", sketchName)
		// Use defaults if no JSON file
		metadata = SketchMetadata{
			Title:        sketchName,
			Description:  "",
			Keywords:     "",
			Tags:         []string{},
			ExternalLibs: []string{},
		}
	} else {
		err = json.Unmarshal(jsonData, &metadata)
		if err != nil {
			return fmt.Errorf("failed to parse JSON metadata for %s: %w", sketchName, err)
		}
	}

	// Extract external libraries from external_libs field
	var externalLibs []string
	if metadata.ExternalLibs != nil {
		externalLibs = metadata.ExternalLibs
	}

	// Parse dates if provided in metadata
	var createdAt, updatedAt *time.Time
	if metadata.CreatedAt != nil {
		if parsedCreatedAt, err := time.Parse(time.RFC3339, *metadata.CreatedAt); err == nil {
			createdAt = &parsedCreatedAt
		} else {
			log.Printf("    Warning: Failed to parse created_at for %s: %v", sketchName, err)
		}
	}
	if metadata.UpdatedAt != nil {
		if parsedUpdatedAt, err := time.Parse(time.RFC3339, *metadata.UpdatedAt); err == nil {
			updatedAt = &parsedUpdatedAt
		} else {
			log.Printf("    Warning: Failed to parse updated_at for %s: %v", sketchName, err)
		}
	}

	// Create the sketch request
	req := &model.CreateSketchRequest{
		Title:        metadata.Title,
		Description:  metadata.Description,
		Keywords:     metadata.Keywords,
		Tags:         metadata.Tags,
		ExternalLibs: externalLibs,
		SourceCode:   string(sourceCode),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	// Create the sketch in database
	sketch, err := services.Sketch.CreateSketch(memberID, req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Printf("    Sketch %s already exists for member, skipping", sketchName)
			return nil
		}
		return fmt.Errorf("failed to create sketch in database: %w", err)
	}

	log.Printf("    Created sketch ID: %d", sketch.ID)
	return nil
}
