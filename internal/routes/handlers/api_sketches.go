package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

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
