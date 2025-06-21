package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// SignInPageData holds data for the sign-in page
type SignInPageData struct {
	utils.PageData
	Error string
}

// SignInGetHandler shows the sign-in form
func SignInGetHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	pageData.Title = utils.Translate(pageData.Lang, "pages.signIn.meta.title")
	pageData.Description = utils.Translate(pageData.Lang, "pages.signIn.meta.description")

	templateData := SignInPageData{
		PageData: *pageData,
	}

	tmplClone, err := tmpl.Clone()
	if err != nil {
		http.Error(w, "Error cloning template", http.StatusInternalServerError)
		log.Printf("Error cloning template for sign-in: %v", err)
		return
	}

	if err := tmplClone.ExecuteTemplate(w, "page-sign-in", templateData); err != nil {
		http.Error(w, "Error rendering page-sign-in template", http.StatusInternalServerError)
		log.Printf("Error rendering page-sign-in template: %v", err)
	}
}

// SignInPostHandler processes the sign-in form
func SignInPostHandler(services *services.Services) func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	return func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
		pageData.Title = utils.Translate(pageData.Lang, "pages.signIn.meta.title")
		pageData.Description = utils.Translate(pageData.Lang, "pages.signIn.meta.description")

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		password := r.FormValue("password")

		templateData := SignInPageData{
			PageData: *pageData,
		}

		if name == "" || password == "" {
			templateData.Error = "Name and password are required"
		} else {
			// Try to authenticate
			if services == nil {
				templateData.Error = "Service unavailable"
			} else {
				// Get member by name
				member, err := services.Member.GetMemberByName(name)
				if err != nil {
					templateData.Error = "Invalid name or password"
				} else if !utils.VerifyPassword(password, member.PasswordHash) {
					templateData.Error = "Invalid name or password"
				} else {
					// Create session
					session, err := services.Session.CreateSession(member.ID)
					if err != nil {
						templateData.Error = "Error creating session"
						log.Printf("Error creating session for member %d: %v", member.ID, err)
					} else {
						// Set session cookie and redirect
						utils.SetSessionCookie(w, session.ID)
						http.Redirect(w, r, "/", http.StatusSeeOther)
						return
					}
				}
			}
		}

		tmplClone, err := tmpl.Clone()
		if err != nil {
			http.Error(w, "Error cloning template", http.StatusInternalServerError)
			log.Printf("Error cloning template for page-sign-in: %v", err)
			return
		}

		if err := tmplClone.ExecuteTemplate(w, "page-sign-in", templateData); err != nil {
			http.Error(w, "Error rendering page-sign-in template", http.StatusInternalServerError)
			log.Printf("Error rendering page-sign-in template: %v", err)
		}
	}
}
