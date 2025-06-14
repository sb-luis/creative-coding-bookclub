package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/model"
	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// SketchListerPageData holds all data for the sketch lister page template.
type SketchListerPageData struct {
	utils.PageData
	Sketches []model.SketchInfo
}

// SketchListerPageHandler handles requests to display the list of all sketches.
func SketchListerPageHandler(services *services.Services) func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	return func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
		// Override metadata fields specific to the sketch lister page
		pageData.Title = utils.Translate(pageData.Lang, "pages.sketchLister.meta.title")
		pageData.Description = utils.Translate(pageData.Lang, "pages.sketchLister.meta.description")
		pageData.Keywords = utils.Translate(pageData.Lang, "pages.sketchLister.meta.keywords")

		// Get all sketches from services in chronological order
		if services == nil {
			log.Printf("Services not initialized")
			http.Error(w, "Failed to load sketches", http.StatusInternalServerError)
			return
		}

		sketchesData, err := services.Sketch.GetAllSketchesChronological()
		if err != nil {
			log.Printf("Error getting sketches from services: %v", err)
			http.Error(w, "Failed to load sketches", http.StatusInternalServerError)
			return
		}

		// Prepare the data for the template
		templateData := SketchListerPageData{
			PageData: *pageData,
			Sketches: sketchesData,
		}

		err = tmpl.ExecuteTemplate(w, "page-sketch-lister", templateData)
		if err != nil {
			log.Printf("Error executing page-sketch-lister template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
