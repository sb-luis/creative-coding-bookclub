package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// SketchPageData holds all data for the sketch page template.
type SketchPageData struct {
	utils.PageData
	MemberName   string
	SketchSlug   string
	SketchJsPath string
	ExternalLibs []string
	SourceCode   string
}

// SketchEditorPageHandler handles requests to display a sketch page.
func SketchEditorPageHandler(services *services.Services) func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	return func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
		memberName := utils.PathVariable(r, "memberName")
		sketchSlug := utils.PathVariable(r, "sketchSlug")

		// Get member and sketch from services
		if services == nil {
			log.Printf("Services not initialized")
			NotFoundHandler(w, r, tmpl, pageData)
			return
		}

		member, err := services.Member.GetMemberByName(memberName)
		if err != nil {
			log.Printf("Member not found: %s", memberName)
			NotFoundHandler(w, r, tmpl, pageData)
			return
		}

		sketch, err := services.Sketch.GetSketchByMemberAndSlug(member.ID, sketchSlug)
		if err != nil {
			log.Printf("Sketch not found: %s by member %s", sketchSlug, memberName)
			NotFoundHandler(w, r, tmpl, pageData)
			return
		}

		// Use metadata from database
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

		templateData := SketchPageData{
			PageData:     *pageData,
			MemberName:   memberName,
			SketchSlug:   sketchSlug,
			SketchJsPath: sketchJsPath,
			ExternalLibs: sketch.ExternalLibs,
			SourceCode:   sketch.SourceCode,
		}

		log.Printf("Rendering sketch page for member: %s, sketch: %s (JS served from database)",
			memberName, sketchSlug)

		err = tmpl.ExecuteTemplate(w, "page-sketch-editor", templateData)
		if err != nil {
			log.Printf("Error executing page-sketch-editor template: %v", err)
			http.Error(w, "Internal Server Error executing template", http.StatusInternalServerError)
		}
	}
}
