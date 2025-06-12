package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// SketchViewPageData holds all data for the clean sketch view template.
type SketchViewPageData struct {
	utils.PageData
	MemberName   string
	SketchSlug   string
	SketchJsPath string
	ExternalLibs []string
	Title        string
}

// SketchViewerPageHandler handles requests to display a clean sketch view (no editor).
// This renders the sketch.html template for pure sketch viewing without editing capabilities.
func SketchViewerPageHandler(services *services.Services) func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	return func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
		memberName := utils.PathVariable(r, "memberName")
		sketchSlug := utils.PathVariable(r, "sketchSlug")

		// Get member and sketch from services
		if services == nil {
			log.Printf("Services not initialized")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		member, err := services.Member.GetMemberByName(memberName)
		if err != nil {
			log.Printf("Member not found: %s, error: %v", memberName, err)
			http.Error(w, "Member not found", http.StatusNotFound)
			return
		}

		sketch, err := services.Sketch.GetSketchByMemberAndSlug(member.ID, sketchSlug)
		if err != nil {
			log.Printf("Sketch not found: %s/%s, error: %v", memberName, sketchSlug, err)
			http.Error(w, "Sketch not found", http.StatusNotFound)
			return
		}

		// Use metadata from database for SEO
		if sketch.Title != "" {
			pageData.Title = sketch.Title
		}
		if sketch.Description != "" {
			pageData.Description = sketch.Description
		}
		if sketch.Keywords != "" {
			pageData.Keywords = sketch.Keywords
		}

		// Point to the database-served JavaScript endpoint
		sketchJsPath := "/api/sketches/" + memberName + "/" + sketchSlug

		templateData := SketchViewPageData{
			PageData:     *pageData,
			MemberName:   memberName,
			SketchSlug:   sketchSlug,
			SketchJsPath: sketchJsPath,
			ExternalLibs: sketch.ExternalLibs,
			Title:        sketch.Title,
		}

		log.Printf("Rendering clean sketch view for member: %s, sketch: %s (JS served from database)",
			memberName, sketchSlug)

		err = tmpl.ExecuteTemplate(w, "page-sketch-viewer", templateData)
		if err != nil {
			log.Printf("Error executing page-sketch-viewer template: %v", err)
			http.Error(w, "Internal Server Error executing template", http.StatusInternalServerError)
		}
	}
}
