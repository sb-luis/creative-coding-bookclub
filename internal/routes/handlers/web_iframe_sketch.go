package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// SketchIframeContentData holds data for the iframe content template.
type SketchIframeContentData struct {
	utils.PageData
	MemberName   string
	SketchSlug   string
	SketchJsPath string
	ExternalLibs []string
	Title        string
}

// SketchIframeContentHandler handles requests to display sketch content inside a sandboxed iframe.
// This renders the sketch-iframe-content.html template for the iframe content.
func SketchIframeContentHandler(services *services.Services) func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
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

		// Point to the database-served JavaScript endpoint
		sketchJsPath := "/api/sketches/" + memberName + "/" + sketchSlug

		templateData := SketchIframeContentData{
			PageData:     *pageData,
			MemberName:   memberName,
			SketchSlug:   sketchSlug,
			SketchJsPath: sketchJsPath,
			ExternalLibs: sketch.ExternalLibs,
			Title:        sketch.Title,
		}

		log.Printf("Rendering iframe content for sketch: %s/%s", memberName, sketchSlug)

		// Add security headers for iframe content
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Content-Security-Policy", "default-src 'self' https:; script-src 'self' 'unsafe-eval' 'unsafe-inline' https:; style-src 'self' 'unsafe-inline';")

		err = tmpl.ExecuteTemplate(w, "page-iframe-sketch", templateData)
		if err != nil {
			log.Printf("Error executing page-iframe-sketch template: %v", err)
			http.Error(w, "Internal Server Error executing template", http.StatusInternalServerError)
		}
	}
}
