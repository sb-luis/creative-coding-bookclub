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
	err := s.db.QueryRow("SELECT COUNT(*) FROM sketches WHERE member_id = ? AND slug = ?", memberID, slug).Scan(&count)
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
	result, err := s.db.Exec(`
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

	return s.GetSketchByID(int(id))
}

// GetSketchByID returns a sketch by ID
func (s *Service) GetSketchByID(id int) (*model.Sketch, error) {
	if id <= 0 {
		return nil, errors.New("invalid sketch ID")
	}

	sketch := &model.Sketch{}
	err := s.db.QueryRow(`
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
func (s *Service) GetSketchByMemberAndSlug(memberID int, slug string) (*model.Sketch, error) {
	if memberID <= 0 {
		return nil, errors.New("invalid member ID")
	}
	if slug == "" {
		return nil, errors.New("slug cannot be empty")
	}

	sketch := &model.Sketch{}
	err := s.db.QueryRow(`
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
func (s *Service) GetSketchesByMember(memberID int) ([]*model.Sketch, error) {
	if memberID <= 0 {
		return nil, errors.New("invalid member ID")
	}

	rows, err := s.db.Query(`
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
func (s *Service) GetAllSketches() ([]*model.Sketch, error) {
	rows, err := s.db.Query(`
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
func (s *Service) GetAllSketchesGroupedByMember() ([]model.MemberSketchInfo, error) {
	rows, err := s.db.Query(`
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
			err := s.db.QueryRow("SELECT COUNT(*) FROM sketches WHERE member_id = ? AND slug = ? AND id != ?",
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

	result, err := s.db.Exec("DELETE FROM sketches WHERE id = ?", id)
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

	result, err := s.db.Exec("DELETE FROM sketches WHERE member_id = ? AND slug = ?", memberID, slug)
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
