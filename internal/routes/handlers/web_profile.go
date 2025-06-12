package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// ProfilePageData holds data for the profile page
type ProfilePageData struct {
	utils.PageData
	Name     string
	MemberID int
}

// ProfileHandler shows the authenticated member's profile
func ProfileHandler(services *services.Services) func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	return func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
		// Get the current member
		sessionID, err := utils.GetSessionFromRequest(r)
		if err != nil {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		memberID, err := services.Session.GetMemberIDFromSession(sessionID)
		if err != nil {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		member, err := services.Member.GetMemberByID(memberID)
		if err != nil {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		pageData.Title = utils.Translate(pageData.Lang, "pages.profile.meta.title")
		pageData.Description = utils.Translate(pageData.Lang, "pages.profile.meta.description")

		templateData := ProfilePageData{
			PageData: *pageData,
			Name:     member.Name,
			MemberID: member.ID,
		}

		tmplClone, err := tmpl.Clone()
		if err != nil {
			http.Error(w, "Error cloning template", http.StatusInternalServerError)
			log.Printf("Error cloning template for profile: %v", err)
			return
		}

		if err := tmplClone.ExecuteTemplate(w, "page-profile", templateData); err != nil {
			http.Error(w, "Error rendering page-profile template", http.StatusInternalServerError)
			log.Printf("Error rendering page-profile template: %v", err)
		}
	}
}
