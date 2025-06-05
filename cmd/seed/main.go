package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sb-luis/creative-coding-bookclub/internal/model"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// SketchMetadata represents the structure of sketch JSON files
type SketchMetadata struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Keywords     string   `json:"keywords"`
	Tags         []string `json:"tags"`
	ExternalLibs []string `json:"external_libs"`
}

func main() {
	// Initialize database
	err := utils.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer utils.CloseDatabase()

	sketchesDir := "data/sketches"

	// Read all member directories
	entries, err := os.ReadDir(sketchesDir)
	if err != nil {
		log.Fatalf("Failed to read sketches directory: %v", err)
	}

	// Process each member directory
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		memberName := entry.Name()
		log.Printf("Processing member: %s", memberName)

		// Create member with 'qwerty' password
		err := createMember(memberName, "qwerty")
		if err != nil {
			log.Printf("Failed to create member %s: %v", memberName, err)
			continue
		}

		// Process sketches for this member
		memberDir := filepath.Join(sketchesDir, memberName)
		err = processMemberSketches(memberName, memberDir)
		if err != nil {
			log.Printf("Failed to process sketches for member %s: %v", memberName, err)
		}
	}

	log.Println("Seeding completed!")
}

func createMember(name, password string) error {
	db := utils.GetDB()
	passwordHash := utils.HashPassword(password)

	// Try to create the member
	member, err := db.CreateMember(name, passwordHash)
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

func processMemberSketches(memberName, memberDir string) error {
	db := utils.GetDB()

	// Get member from database
	member, err := db.GetMemberByName(memberName)
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

			err := createSketch(member.ID, memberDir, sketchName)
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

func createSketch(memberID int, memberDir, sketchName string) error {
	db := utils.GetDB()

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

	// Create the sketch request
	req := &model.CreateSketchRequest{
		Title:        metadata.Title,
		Description:  metadata.Description,
		Keywords:     metadata.Keywords,
		Tags:         metadata.Tags,
		ExternalLibs: externalLibs,
		SourceCode:   string(sourceCode),
	}

	// Create the sketch in database
	sketch, err := db.CreateSketch(memberID, req)
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
