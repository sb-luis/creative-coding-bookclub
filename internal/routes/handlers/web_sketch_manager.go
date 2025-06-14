package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// SketchManagerPageData holds all data for the sketch manager page template.
type SketchManagerPageData struct {
	utils.PageData
	MemberID   int
	MemberName string
}

// SketchManagerPageHandler handles requests to display the sketch manager page
func SketchManagerPageHandler(services *services.Services) func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	return func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
		// Check if user is authenticated
		sessionID, err := utils.GetSessionFromRequest(r)
		if err != nil {
			// Redirect to sign-in if not authenticated
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		if services == nil {
			log.Printf("Services not initialized")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Get member ID from session
		memberID, err := services.Session.GetMemberIDFromSession(sessionID)
		if err != nil {
			// Redirect to sign-in if session is invalid
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		// Get member details
		member, err := services.Member.GetMemberByID(memberID)
		if err != nil {
			log.Printf("Member not found for ID %d: %v", memberID, err)
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		// Override metadata fields specific to the sketch manager page
		pageData.Title = utils.Translate(pageData.Lang, "pages.sketchManager.meta.title")
		pageData.Description = utils.Translate(pageData.Lang, "pages.sketchManager.meta.description")
		pageData.Keywords = utils.Translate(pageData.Lang, "pages.sketchManager.meta.keywords")

		templateData := SketchManagerPageData{
			PageData:   *pageData,
			MemberID:   memberID,
			MemberName: member.Name,
		}

		err = tmpl.ExecuteTemplate(w, "page-sketch-manager", templateData)
		if err != nil {
			log.Printf("Error executing page-sketch-manager template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
