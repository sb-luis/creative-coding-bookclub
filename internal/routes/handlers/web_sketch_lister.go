package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/model"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// SketchListerPageData holds all data for the sketch lister page template.
type SketchListerPageData struct {
	utils.PageData
	Members []model.MemberSketchInfo
}

// SketchListerPageHandler handles requests to display the list of all sketches.
func SketchListerPageHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	// Override metadata fields specific to the sketch lister page
	pageData.Title = utils.Translate(pageData.Lang, "pages.sketchLister.meta.title")
	pageData.Description = utils.Translate(pageData.Lang, "pages.sketchLister.meta.description")
	pageData.Keywords = utils.Translate(pageData.Lang, "pages.sketchLister.meta.keywords")

	// Get all sketches from database grouped by member
	db := utils.GetDB()
	membersData, err := db.GetAllSketchesGroupedByMember()
	if err != nil {
		log.Printf("Error getting sketches from database: %v", err)
		http.Error(w, "Failed to load sketches", http.StatusInternalServerError)
		return
	}

	// Prepare the data for the template
	templateData := SketchListerPageData{
		PageData: *pageData,
		Members:  membersData, // Renamed from MembersData
	}

	err = tmpl.ExecuteTemplate(w, "page-sketch-lister", templateData)
	if err != nil {
		log.Printf("Error executing page-sketch-lister template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
