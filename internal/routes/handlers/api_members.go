package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// MemberResponse represents a member in API responses
type MemberResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetMembersHandler handles GET requests to return a list of all members (IDs and Names)
func GetMembersHandler(w http.ResponseWriter, r *http.Request) {
	// Set content type for JSON response
	w.Header().Set("Content-Type", "application/json")

	// Get database connection
	db := utils.GetDB()
	if db == nil {
		log.Printf("Database not initialized")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get all members from database
	members, err := db.GetAllMembers()
	if err != nil {
		log.Printf("Error getting all members: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert to response format (exclude password hash and other sensitive data)
	var memberResponses []MemberResponse
	for _, member := range members {
		memberResponses = append(memberResponses, MemberResponse{
			ID:   member.ID,
			Name: member.Name,
		})
	}

	// Return JSON response
	if err := json.NewEncoder(w).Encode(memberResponses); err != nil {
		log.Printf("Error encoding JSON response for members: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Served %d members via API", len(memberResponses))
}

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

// GetMemberSketchesHandler handles GET requests to return all sketches for a specific member
func GetMemberSketchesHandler(w http.ResponseWriter, r *http.Request) {
	// Set content type for JSON response
	w.Header().Set("Content-Type", "application/json")

	// Get member name from URL path
	memberName := utils.PathVariable(r, "member")
	if memberName == "" {
		log.Printf("Invalid request: missing member name")
		http.Error(w, "Bad Request: missing member name", http.StatusBadRequest)
		return
	}

	// Get database connection
	db := utils.GetDB()
	if db == nil {
		log.Printf("Database not initialized")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get member from database
	member, err := db.GetMemberByName(memberName)
	if err != nil {
		log.Printf("Member not found: %s", memberName)
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	// Get all sketches for this member
	sketches, err := db.GetSketchesByMember(member.ID)
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
