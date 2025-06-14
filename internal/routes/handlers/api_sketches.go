package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/sb-luis/creative-coding-bookclub/internal/model"
	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// Validation constants
const (
	MaxTitleLength       = 100
	MaxDescriptionLength = 500
	MaxKeywordsLength    = 200
	MaxTagsCount         = 10
	MaxExternalLibsCount = 5
)

// Precompiled regex patterns for validation
var (
	titleRegex       = regexp.MustCompile(`^[a-zA-Z0-9\s\-_.,:;!?()]+$`)
	descriptionRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\-_.,:;!?()]*$`)
	keywordsRegex    = regexp.MustCompile(`^[a-zA-Z0-9\s,\-_]*$`)
	tagRegex         = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
)

// Validation functions
func validateTitle(title string) error {
	if len(title) == 0 {
		return fmt.Errorf("title is required")
	}
	if len(title) > MaxTitleLength {
		return fmt.Errorf("title must be %d characters or less", MaxTitleLength)
	}
	// Allow alphanumeric, spaces, and common punctuation
	if !titleRegex.MatchString(title) {
		return fmt.Errorf("title contains invalid characters")
	}
	return nil
}

func validateDescription(description string) error {
	if len(description) > MaxDescriptionLength {
		return fmt.Errorf("description must be %d characters or less", MaxDescriptionLength)
	}
	// Allow alphanumeric, spaces, and common punctuation
	if description != "" && !descriptionRegex.MatchString(description) {
		return fmt.Errorf("description contains invalid characters")
	}
	return nil
}

func validateKeywords(keywords string) error {
	if len(keywords) > MaxKeywordsLength {
		return fmt.Errorf("keywords must be %d characters or less", MaxKeywordsLength)
	}
	// Allow alphanumeric, spaces, commas, and hyphens
	if keywords != "" && !keywordsRegex.MatchString(keywords) {
		return fmt.Errorf("keywords contains invalid characters")
	}
	return nil
}

func validateTags(tags []string) error {
	if len(tags) > MaxTagsCount {
		return fmt.Errorf("maximum %d tags allowed", MaxTagsCount)
	}
	// Alphanumeric only, no spaces or whitespace
	for _, tag := range tags {
		if strings.TrimSpace(tag) == "" {
			return fmt.Errorf("tags cannot be empty")
		}
		if !tagRegex.MatchString(tag) {
			return fmt.Errorf("tag '%s' contains invalid characters (only alphanumeric, hyphens, and underscores allowed)", tag)
		}
	}
	return nil
}

func validateExternalLibs(libs []string) error {
	if len(libs) > MaxExternalLibsCount {
		return fmt.Errorf("maximum %d external libraries allowed", MaxExternalLibsCount)
	}
	for _, lib := range libs {
		if strings.TrimSpace(lib) == "" {
			return fmt.Errorf("external library URLs cannot be empty")
		}
		// Validate URL format
		if _, err := url.ParseRequestURI(lib); err != nil {
			return fmt.Errorf("invalid URL format: %s", lib)
		}
		// Ensure it's HTTPS
		if !strings.HasPrefix(lib, "https://") {
			return fmt.Errorf("external library URLs must use HTTPS: %s", lib)
		}
	}
	return nil
}

// Request structs for sketch endpoints

// SketchCreateRequest represents the payload for creating a new sketch
type SketchCreateRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Keywords     string   `json:"keywords"`
	Tags         []string `json:"tags"`
	ExternalLibs []string `json:"external_libs"`
	SourceCode   string   `json:"source_code"`
}

// SketchUpdateRequest represents the payload for updating source code only (PUT)
type SketchUpdateRequest struct {
	SourceCode string `json:"source_code"`
}

// SketchMetadataUpdateRequest represents the payload for updating metadata only (PATCH)
type SketchMetadataUpdateRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description,omitempty"`
	Keywords     string   `json:"keywords,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	ExternalLibs []string `json:"external_libs,omitempty"`
}

// Response structs

// SketchResponse represents a sketch in API responses (without source code)
type SketchResponse struct {
	ID           int      `json:"id"`
	Slug         string   `json:"slug"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Keywords     string   `json:"keywords"`
	Tags         []string `json:"tags"`
	ExternalLibs []string `json:"external_libs"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

// PUBLIC ENDPOINTS (NO AUTH REQUIRED)

// GetMemberSketchesHandler handles GET requests to return all sketches for a specific member
func GetMemberSketchesHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type for JSON response
		w.Header().Set("Content-Type", "application/json")

		// Get member name from URL path
		memberName := utils.PathVariable(r, "memberName")
		if memberName == "" {
			log.Printf("Invalid request: missing member name")
			http.Error(w, "Bad Request: missing member name", http.StatusBadRequest)
			return
		}

		if services == nil {
			log.Printf("Services not initialized")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Get member from service
		member, err := services.Member.GetMemberByName(memberName)
		if err != nil {
			log.Printf("Member not found: %s", memberName)
			http.Error(w, "Member not found", http.StatusNotFound)
			return
		}

		// Get all sketches for this member
		sketches, err := services.Sketch.GetSketchesByMember(member.ID)
		if err != nil {
			log.Printf("Error getting sketches for member %s: %v", memberName, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert to response format (exclude source code for listing)
		var sketchResponses []SketchResponse
		for _, sketch := range sketches {
			sketchResponses = append(sketchResponses, SketchResponse{
				ID:           sketch.ID,
				Slug:         sketch.Slug,
				Title:        sketch.Title,
				Description:  sketch.Description,
				Keywords:     sketch.Keywords,
				Tags:         sketch.Tags,
				ExternalLibs: sketch.ExternalLibs,
				CreatedAt:    sketch.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
				UpdatedAt:    sketch.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			})
		}

		// Return JSON response
		if err := json.NewEncoder(w).Encode(sketchResponses); err != nil {
			log.Printf("Error encoding JSON response for member sketches: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Printf("Served %d sketches for member %s via API", len(sketchResponses), memberName)
	}
}

// SketchCodeHandler handles requests to serve JavaScript files from the database.
// This handler serves the JS source code stored in the database for a specific sketch.
func SketchCodeHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memberName := utils.PathVariable(r, "memberName")
		sketchSlug := utils.PathVariable(r, "sketchSlug")

		if memberName == "" || sketchSlug == "" {
			log.Printf("Invalid request: missing member or sketch slug")
			http.Error(w, "Bad Request: missing member or sketch slug", http.StatusBadRequest)
			return
		}

		if services == nil {
			log.Printf("Services not initialized")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Get member from service
		member, err := services.Member.GetMemberByName(memberName)
		if err != nil {
			log.Printf("Member not found: %s", memberName)
			http.NotFound(w, r)
			return
		}

		// Get sketch from service
		sketch, err := services.Sketch.GetSketchByMemberAndSlug(member.ID, sketchSlug)
		if err != nil {
			log.Printf("Sketch not found: %s by member %s", sketchSlug, memberName)
			http.NotFound(w, r)
			return
		}

		// Set appropriate content type for JavaScript
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")

		// Set cache headers to allow reasonable caching but still allow updates
		w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes cache

		// Write the JavaScript source code
		_, err = w.Write([]byte(sketch.SourceCode))
		if err != nil {
			log.Printf("Error writing JavaScript response for sketch %s/%s: %v", memberName, sketchSlug, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Printf("Served JavaScript for sketch: %s/%s (from database)", memberName, sketchSlug)
	}
}

// PROTECTED ENDPOINTS (AUTH REQUIRED)

// CreateSketchHandler handles POST requests to create a new sketch
func CreateSketchHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type for JSON response
		w.Header().Set("Content-Type", "application/json")

		// Get authenticated member ID from context
		memberID, ok := r.Context().Value("authenticated_member_id").(int)
		if !ok {
			http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
			return
		}

		// Get member name and sketch slug from URL
		memberName := utils.PathVariable(r, "memberName")
		sketchSlug := utils.PathVariable(r, "sketchSlug")

		if memberName == "" || sketchSlug == "" {
			log.Printf("Invalid request: missing member name or sketch slug")
			http.Error(w, `{"error":"Member name and sketch slug are required"}`, http.StatusBadRequest)
			return
		}

		// Verify the authenticated user is creating a sketch for themselves
		member, err := services.Member.GetMemberByID(memberID)
		if err != nil || member.Name != memberName {
			http.Error(w, `{"error":"You can only create sketches for your own account"}`, http.StatusForbidden)
			return
		}

		// Check if sketch with this slug already exists for this member
		exists, err := services.Sketch.SketchExistsByMemberAndSlug(memberID, sketchSlug)
		if err != nil {
			log.Printf("Error checking if sketch exists: %v", err)
			http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}

		if exists {
			http.Error(w, `{"error":"A sketch with this slug already exists"}`, http.StatusConflict)
			return
		}

		// Parse request body
		var req SketchCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding create sketch request: %v", err)
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		// Validate required fields
		if req.Title == "" {
			http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
			return
		}
		if req.SourceCode == "" {
			http.Error(w, `{"error":"Source code is required"}`, http.StatusBadRequest)
			return
		}

		// Validate metadata fields
		if err := validateTitle(req.Title); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		if err := validateDescription(req.Description); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		if err := validateKeywords(req.Keywords); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		if err := validateTags(req.Tags); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		if err := validateExternalLibs(req.ExternalLibs); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		// Create sketch request
		createReq := &model.CreateSketchRequest{
			Title:        req.Title,
			Description:  req.Description,
			Keywords:     req.Keywords,
			Tags:         req.Tags,
			ExternalLibs: req.ExternalLibs,
			SourceCode:   req.SourceCode,
		}

		// Override slug with the one from URL
		sketch, err := services.Sketch.CreateSketchWithSlug(memberID, createReq, sketchSlug)
		if err != nil {
			log.Printf("Error creating sketch for member %d: %v", memberID, err)
			http.Error(w, `{"error":"Failed to create sketch"}`, http.StatusInternalServerError)
			return
		}

		// Return created sketch
		if err := json.NewEncoder(w).Encode(sketch); err != nil {
			log.Printf("Error encoding sketch response: %v", err)
			http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}

		log.Printf("Created sketch %s for member %s (ID: %d)", sketchSlug, memberName, memberID)
	}
}

// UpdateSketchHandler handles PUT requests to update sketch source code only
func UpdateSketchHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type for JSON response
		w.Header().Set("Content-Type", "application/json")

		// Get authenticated member ID from context (set by authMiddleware)
		memberID, ok := r.Context().Value("authenticated_member_id").(int)
		if !ok {
			http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
			return
		}

		// Get member name and sketch slug from URL
		memberName := utils.PathVariable(r, "memberName")
		sketchSlug := utils.PathVariable(r, "sketchSlug")

		if memberName == "" || sketchSlug == "" {
			log.Printf("Invalid request: missing member name or sketch slug")
			http.Error(w, `{"error":"Member name and sketch slug are required"}`, http.StatusBadRequest)
			return
		}

		// Verify the authenticated user is updating their own sketch
		member, err := services.Member.GetMemberByID(memberID)
		if err != nil || member.Name != memberName {
			http.Error(w, `{"error":"You can only update your own sketches"}`, http.StatusForbidden)
			return
		}

		// Get the sketch to update
		sketch, err := services.Sketch.GetSketchByMemberAndSlug(memberID, sketchSlug)
		if err != nil {
			log.Printf("Sketch not found: %s for member %s (ID: %d)", sketchSlug, memberName, memberID)
			http.Error(w, `{"error":"Sketch not found"}`, http.StatusNotFound)
			return
		}

		// Parse request body
		var req SketchUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding update sketch request: %v", err)
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		// Validate required field
		if req.SourceCode == "" {
			http.Error(w, `{"error":"Source code is required"}`, http.StatusBadRequest)
			return
		}

		// Create update request (only source code)
		updateReq := &model.UpdateSketchRequest{
			SourceCode: &req.SourceCode,
		}

		// Update sketch
		updatedSketch, err := services.Sketch.UpdateSketch(sketch.ID, updateReq)
		if err != nil {
			log.Printf("Error updating sketch %s for member %s: %v", sketchSlug, memberName, err)
			http.Error(w, `{"error":"Failed to update sketch"}`, http.StatusInternalServerError)
			return
		}

		// Return updated sketch
		if err := json.NewEncoder(w).Encode(updatedSketch); err != nil {
			log.Printf("Error encoding updated sketch response: %v", err)
			http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}

		log.Printf("Updated sketch source code %s for member %s", sketchSlug, memberName)
	}
}

// UpdateSketchMetadataHandler handles PATCH requests to update sketch metadata only
func UpdateSketchMetadataHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type for JSON response
		w.Header().Set("Content-Type", "application/json")

		// Get authenticated member ID from context (set by authMiddleware)
		memberID, ok := r.Context().Value("authenticated_member_id").(int)
		if !ok {
			http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
			return
		}

		// Get member name and sketch slug from URL
		memberName := utils.PathVariable(r, "memberName")
		sketchSlug := utils.PathVariable(r, "sketchSlug")

		if memberName == "" || sketchSlug == "" {
			log.Printf("Invalid request: missing member name or sketch slug")
			http.Error(w, `{"error":"Member name and sketch slug are required"}`, http.StatusBadRequest)
			return
		}

		// Verify the authenticated user is updating their own sketch
		member, err := services.Member.GetMemberByID(memberID)
		if err != nil || member.Name != memberName {
			http.Error(w, `{"error":"You can only update your own sketches"}`, http.StatusForbidden)
			return
		}

		// Get the sketch to update
		sketch, err := services.Sketch.GetSketchByMemberAndSlug(memberID, sketchSlug)
		if err != nil {
			log.Printf("Sketch not found: %s for member %s (ID: %d)", sketchSlug, memberName, memberID)
			http.Error(w, `{"error":"Sketch not found"}`, http.StatusNotFound)
			return
		}

		// Parse request body
		var req SketchMetadataUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding update sketch metadata request: %v", err)
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		// Validate required field
		if req.Title == "" {
			http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
			return
		}

		// Validate metadata fields
		if err := validateTitle(req.Title); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		if err := validateDescription(req.Description); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		if err := validateKeywords(req.Keywords); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		if err := validateTags(req.Tags); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}
		if err := validateExternalLibs(req.ExternalLibs); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		// Create update request (only metadata fields)
		updateReq := &model.UpdateSketchRequest{
			Title:        &req.Title,
			Description:  &req.Description,
			Keywords:     &req.Keywords,
			Tags:         req.Tags,
			ExternalLibs: req.ExternalLibs,
		}

		// Update sketch metadata (this will also update the slug and updated_at automatically)
		updatedSketch, err := services.Sketch.UpdateSketch(sketch.ID, updateReq)
		if err != nil {
			log.Printf("Error updating sketch metadata %s for member %s: %v", sketchSlug, memberName, err)
			http.Error(w, `{"error":"Failed to update sketch metadata"}`, http.StatusInternalServerError)
			return
		}

		// Return updated sketch
		if err := json.NewEncoder(w).Encode(updatedSketch); err != nil {
			log.Printf("Error encoding updated sketch response: %v", err)
			http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}

		log.Printf("Updated sketch metadata %s for member %s", sketchSlug, memberName)
	}
}

// DeleteSketchHandler handles DELETE requests to delete a sketch
func DeleteSketchHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type for JSON response
		w.Header().Set("Content-Type", "application/json")

		// Get authenticated member ID from context (set by authMiddleware)
		memberID, ok := r.Context().Value("authenticated_member_id").(int)
		if !ok {
			http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
			return
		}

		// Get member name and sketch slug from URL
		memberName := utils.PathVariable(r, "memberName")
		sketchSlug := utils.PathVariable(r, "sketchSlug")

		if memberName == "" || sketchSlug == "" {
			log.Printf("Invalid request: missing member name or sketch slug")
			http.Error(w, `{"error":"Member name and sketch slug are required"}`, http.StatusBadRequest)
			return
		}

		// Verify the authenticated user is deleting their own sketch
		member, err := services.Member.GetMemberByID(memberID)
		if err != nil || member.Name != memberName {
			http.Error(w, `{"error":"You can only delete your own sketches"}`, http.StatusForbidden)
			return
		}

		// Get the sketch to delete
		sketch, err := services.Sketch.GetSketchByMemberAndSlug(memberID, sketchSlug)
		if err != nil {
			log.Printf("Sketch not found: %s for member %s (ID: %d)", sketchSlug, memberName, memberID)
			http.Error(w, `{"error":"Sketch not found"}`, http.StatusNotFound)
			return
		}

		// Delete sketch
		err = services.Sketch.DeleteSketch(sketch.ID)
		if err != nil {
			log.Printf("Error deleting sketch %s for member %s: %v", sketchSlug, memberName, err)
			http.Error(w, `{"error":"Failed to delete sketch"}`, http.StatusInternalServerError)
			return
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Sketch deleted successfully"})

		log.Printf("Deleted sketch %s for member %s", sketchSlug, memberName)
	}
}
