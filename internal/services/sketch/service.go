package sketch

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/sb-luis/creative-coding-bookclub/internal/model"
)

// generateSlug creates a URL-friendly slug from a title
func generateSlug(title string) string {
	if title == "" {
		return ""
	}

	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces and non-alphanumeric characters with hyphens
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// Remove consecutive hyphens
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	return slug
}

// Service handles sketch-related business logic
type Service struct {
	db *sql.DB
}

// NewService creates a new sketch service
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// CreateSketch creates a new sketch for a member
func (s *Service) CreateSketch(memberID int, req *model.CreateSketchRequest) (*model.Sketch, error) {
	if memberID <= 0 {
		return nil, errors.New("invalid member ID")
	}
	if req == nil {
		return nil, errors.New("create sketch request cannot be nil")
	}
	if req.Title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if req.SourceCode == "" {
		return nil, errors.New("source code cannot be empty")
	}

	// Generate slug from title
	slug := generateSlug(req.Title)
	if slug == "" {
		return nil, errors.New("title cannot be empty or invalid")
	}

	// Check if sketch slug already exists for this member
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM sketches WHERE member_id = $1 AND slug = $2", memberID, slug).Scan(&count)
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
	createdAt := now
	updatedAt := now

	// Use custom dates if provided
	if req.CreatedAt != nil {
		createdAt = *req.CreatedAt
	}
	if req.UpdatedAt != nil {
		updatedAt = *req.UpdatedAt
	}

	var id int
	err = s.db.QueryRow(`
		INSERT INTO sketches (member_id, slug, title, description, keywords, tags, external_libs, source_code, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		memberID, slug, req.Title, req.Description, req.Keywords, string(tagsJSON), string(externalLibsJSON), req.SourceCode, createdAt, updatedAt).Scan(&id)
	if err != nil {
		log.Printf("Database error while creating sketch for member %d: %v", memberID, err)
		return nil, fmt.Errorf("failed to create sketch: %w", err)
	}

	return s.GetSketchByID(id)
}

// GetSketchByID returns a sketch by ID
func (s *Service) GetSketchByID(id int) (*model.Sketch, error) {
	if id <= 0 {
		return nil, errors.New("invalid sketch ID")
	}

	sketch := &model.Sketch{}
	err := s.db.QueryRow(`
		SELECT id, member_id, slug, title, description, keywords, tags, external_libs, source_code, created_at, updated_at 
		FROM sketches WHERE id = $1`, id).Scan(
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
func (s *Service) GetSketchByMemberAndSlug(memberID int, slug string) (*model.Sketch, error) {
	if memberID <= 0 {
		return nil, errors.New("invalid member ID")
	}
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}

	// Handle the special "new" slug - return a default new sketch template
	if slug == "new" {
		defaultCode := `
function setup() {
    createCanvas(400, 400);
}

function draw() {
    background(220);
    fill(255, 0, 150);
    ellipse(mouseX, mouseY, 50, 50);
}`

		return &model.Sketch{
			ID:           0,
			MemberID:     memberID,
			Slug:         "new",
			Title:        "New Sketch",
			Description:  "A new creative coding sketch",
			Keywords:     "creative coding, p5js, sketch",
			Tags:         []string{"creative-coding", "p5js"},
			ExternalLibs: []string{"https://cdn.jsdelivr.net/npm/p5@1.11.7/lib/p5.min.js"},
			SourceCode:   defaultCode,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}, nil
	}

	sketch := &model.Sketch{}
	err := s.db.QueryRow(`
		SELECT id, member_id, slug, title, description, keywords, tags, external_libs, source_code, created_at, updated_at 
		FROM sketches WHERE member_id = $1 AND slug = $2`, memberID, slug).Scan(
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
func (s *Service) GetSketchesByMember(memberID int) ([]*model.Sketch, error) {
	if memberID <= 0 {
		return nil, errors.New("invalid member ID")
	}

	rows, err := s.db.Query(`
		SELECT id, member_id, slug, title, description, keywords, tags, external_libs, source_code, created_at, updated_at 
		FROM sketches WHERE member_id = $1 ORDER BY updated_at DESC`, memberID)
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
func (s *Service) GetAllSketches() ([]*model.Sketch, error) {
	rows, err := s.db.Query(`
		SELECT id, member_id, slug, title, description, keywords, tags, external_libs, source_code, created_at, updated_at 
		FROM sketches ORDER BY updated_at DESC`)
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
func (s *Service) GetAllSketchesGroupedByMember() ([]model.MemberSketchInfo, error) {
	rows, err := s.db.Query(`
		SELECT s.id, s.member_id, s.slug, s.title, s.description, s.keywords, s.tags, s.external_libs, s.source_code, s.created_at, s.updated_at, m.name as member_name
		FROM sketches s
		JOIN members m ON s.member_id = m.id
		ORDER BY s.updated_at DESC`)
	if err != nil {
		log.Printf("Database error while getting all sketches grouped by member: %v", err)
		return nil, fmt.Errorf("failed to get all sketches grouped by member: %w", err)
	}
	defer rows.Close()

	// Keep track of members in order and build result as we go
	var result []model.MemberSketchInfo
	memberIndices := make(map[string]int) // maps member name to index in result slice

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
		if err := json.Unmarshal([]byte(sketch.ExternalLibsJSON), &sketch.ExternalLibs); err != nil {
			log.Printf("Warning: failed to unmarshal external libs for sketch %d: %v", sketch.ID, err)
			sketch.ExternalLibs = []string{}
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

		// Check if we already have this member in our result
		if index, exists := memberIndices[memberName]; exists {
			// Add sketch to existing member
			result[index].Sketches = append(result[index].Sketches, sketchInfo)
		} else {
			// Create new member entry
			memberIndices[memberName] = len(result)
			result = append(result, model.MemberSketchInfo{
				Name:     memberName,
				Sketches: []model.SketchInfo{sketchInfo},
			})
		}
	}

	return result, nil
}

// GetAllSketchesChronological returns all sketches in chronological order (not grouped by member)
func (s *Service) GetAllSketchesChronological() ([]model.SketchInfo, error) {
	rows, err := s.db.Query(`
		SELECT s.id, s.member_id, s.slug, s.title, s.description, s.keywords, s.tags, s.external_libs, s.source_code, s.created_at, s.updated_at, m.name as member_name
		FROM sketches s
		JOIN members m ON s.member_id = m.id
		ORDER BY s.updated_at DESC`)
	if err != nil {
		log.Printf("Database error while getting all sketches chronologically: %v", err)
		return nil, fmt.Errorf("failed to get all sketches chronologically: %w", err)
	}
	defer rows.Close()

	var result []model.SketchInfo
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
		if err := json.Unmarshal([]byte(sketch.ExternalLibsJSON), &sketch.ExternalLibs); err != nil {
			log.Printf("Warning: failed to unmarshal external libs for sketch %d: %v", sketch.ID, err)
			sketch.ExternalLibs = []string{}
		}

		// Convert to SketchInfo for the lister
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

		result = append(result, sketchInfo)
	}

	return result, nil
}

// UpdateSketch updates an existing sketch
func (s *Service) UpdateSketch(id int, req *model.UpdateSketchRequest) (*model.Sketch, error) {
	if id <= 0 {
		return nil, errors.New("invalid sketch ID")
	}
	if req == nil {
		return nil, errors.New("update sketch request cannot be nil")
	}

	// Start building the update query
	setParts := []string{}
	args := []interface{}{}
	paramCount := 0

	// Get current sketch to check if title is changing
	currentSketch, err := s.GetSketchByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get current sketch: %w", err)
	}

	if req.Title != nil {
		// If title is changing, generate new slug and check for conflicts
		newSlug := generateSlug(*req.Title)
		if newSlug == "" {
			return nil, errors.New("title cannot be empty or invalid")
		}

		// Only check for slug conflicts if the slug is actually changing
		if newSlug != currentSketch.Slug {
			var count int
			err := s.db.QueryRow("SELECT COUNT(*) FROM sketches WHERE member_id = $1 AND slug = $2 AND id != $3",
				currentSketch.MemberID, newSlug, id).Scan(&count)
			if err != nil {
				log.Printf("Database error while checking slug conflict for sketch %d: %v", id, err)
				return nil, fmt.Errorf("failed to check slug conflict: %w", err)
			}
			if count > 0 {
				return nil, fmt.Errorf("a sketch with the title '%s' already exists for this member (slug conflict: %s)", *req.Title, newSlug)
			}

			paramCount++
			setParts = append(setParts, fmt.Sprintf("slug = $%d", paramCount))
			args = append(args, newSlug)
		}

		paramCount++
		setParts = append(setParts, fmt.Sprintf("title = $%d", paramCount))
		args = append(args, *req.Title)
	}
	if req.Description != nil {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("description = $%d", paramCount))
		args = append(args, *req.Description)
	}
	if req.Keywords != nil {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("keywords = $%d", paramCount))
		args = append(args, *req.Keywords)
	}
	if req.Tags != nil {
		tagsJSON, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags: %w", err)
		}
		paramCount++
		setParts = append(setParts, fmt.Sprintf("tags = $%d", paramCount))
		args = append(args, string(tagsJSON))
	}
	if req.ExternalLibs != nil {
		externalLibsJSON, err := json.Marshal(req.ExternalLibs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal external libs: %w", err)
		}
		paramCount++
		setParts = append(setParts, fmt.Sprintf("external_libs = $%d", paramCount))
		args = append(args, string(externalLibsJSON))
	}
	if req.SourceCode != nil {
		paramCount++
		setParts = append(setParts, fmt.Sprintf("source_code = $%d", paramCount))
		args = append(args, *req.SourceCode)
	}

	if len(setParts) == 0 {
		return nil, errors.New("no fields to update")
	}

	// Always update the updated_at field
	paramCount++
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", paramCount))
	args = append(args, time.Now())

	// Add the ID for the WHERE clause
	paramCount++
	args = append(args, id)

	query := fmt.Sprintf("UPDATE sketches SET %s WHERE id = $%d", strings.Join(setParts, ", "), paramCount)

	_, err = s.db.Exec(query, args...)
	if err != nil {
		log.Printf("Database error while updating sketch %d: %v", id, err)
		return nil, fmt.Errorf("failed to update sketch: %w", err)
	}

	return s.GetSketchByID(id)
}

// DeleteSketch deletes a sketch by ID
func (s *Service) DeleteSketch(id int) error {
	if id <= 0 {
		return errors.New("invalid sketch ID")
	}

	result, err := s.db.Exec("DELETE FROM sketches WHERE id = $1", id)
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
func (s *Service) DeleteSketchByMemberAndSlug(memberID int, slug string) error {
	if memberID <= 0 {
		return errors.New("invalid member ID")
	}
	if slug == "" {
		return errors.New("slug cannot be empty")
	}

	result, err := s.db.Exec("DELETE FROM sketches WHERE member_id = $1 AND slug = $2", memberID, slug)
	if err != nil {
		log.Printf("Database error while deleting sketch for member %d, slug '%s': %v", memberID, slug, err)
		return fmt.Errorf("failed to delete sketch by member and slug: %w", err)
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

// SketchExistsByMemberAndSlug checks if a sketch with the given slug exists for a member
func (s *Service) SketchExistsByMemberAndSlug(memberID int, slug string) (bool, error) {
	if memberID <= 0 {
		return false, errors.New("invalid member ID")
	}
	if slug == "" {
		return false, errors.New("slug cannot be empty")
	}

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM sketches WHERE member_id = $1 AND slug = $2", memberID, slug).Scan(&count)
	if err != nil {
		log.Printf("Database error while checking if sketch exists for member %d and slug '%s': %v", memberID, slug, err)
		return false, fmt.Errorf("failed to check if sketch exists: %w", err)
	}

	return count > 0, nil
}

// CreateSketchWithSlug creates a new sketch with a specific slug
func (s *Service) CreateSketchWithSlug(memberID int, req *model.CreateSketchRequest, slug string) (*model.Sketch, error) {
	if memberID <= 0 {
		return nil, errors.New("invalid member ID")
	}
	if req == nil {
		return nil, errors.New("create sketch request cannot be nil")
	}
	if req.Title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if req.SourceCode == "" {
		return nil, errors.New("source code cannot be empty")
	}
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}

	// Check if sketch slug already exists for this member
	exists, err := s.SketchExistsByMemberAndSlug(memberID, slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("a sketch with slug '%s' already exists for this member", slug)
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
	createdAt := now
	updatedAt := now

	// Use custom dates if provided
	if req.CreatedAt != nil {
		createdAt = *req.CreatedAt
	}
	if req.UpdatedAt != nil {
		updatedAt = *req.UpdatedAt
	}

	var id int
	err = s.db.QueryRow(`
		INSERT INTO sketches (member_id, slug, title, description, keywords, tags, external_libs, source_code, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		memberID, slug, req.Title, req.Description, req.Keywords, string(tagsJSON), string(externalLibsJSON), req.SourceCode, createdAt, updatedAt).Scan(&id)
	if err != nil {
		log.Printf("Database error while creating sketch for member %d with slug '%s': %v", memberID, slug, err)
		return nil, fmt.Errorf("failed to create sketch: %w", err)
	}

	return s.GetSketchByID(id)
}
