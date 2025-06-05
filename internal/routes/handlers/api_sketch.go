package handlers

import (
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// SketchCodeHandler handles requests to serve JavaScript files from the database.
// This handler serves the JS source code stored in the database for a specific sketch.
func SketchCodeHandler(w http.ResponseWriter, r *http.Request) {
	memberName := utils.PathVariable(r, "member")
	sketchSlug := utils.PathVariable(r, "sketch")

	if memberName == "" || sketchSlug == "" {
		log.Printf("Invalid request: missing member or sketch slug")
		http.Error(w, "Bad Request: missing member or sketch slug", http.StatusBadRequest)
		return
	}

	// Get database connection
	db := utils.GetDB()

	// Get member from database
	member, err := db.GetMemberByName(memberName)
	if err != nil {
		log.Printf("Member not found: %s", memberName)
		http.NotFound(w, r)
		return
	}

	// Get sketch from database
	sketch, err := db.GetSketchByMemberAndSlug(member.ID, sketchSlug)
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
