package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// RegisterPageData holds data for the register page
type RegisterPageData struct {
	utils.PageData
	Error   string
	Success string
}

// RegisterGetHandler shows the registration form
func RegisterGetHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	pageData.Title = utils.Translate(pageData.Lang, "pages.register.meta.title")
	pageData.Description = utils.Translate(pageData.Lang, "pages.register.meta.description")

	templateData := RegisterPageData{
		PageData: *pageData,
	}

	tmplClone, err := tmpl.Clone()
	if err != nil {
		http.Error(w, "Error cloning template", http.StatusInternalServerError)
		log.Printf("Error cloning template for register: %v", err)
		return
	}

	if err := tmplClone.ExecuteTemplate(w, "page-register", templateData); err != nil {
		http.Error(w, "Error rendering page-register template", http.StatusInternalServerError)
		log.Printf("Error rendering page-register template: %v", err)
	}
}

// RegisterPostHandler processes the registration form
func RegisterPostHandler(services *services.Services) func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	return func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
		pageData.Title = utils.Translate(pageData.Lang, "pages.register.meta.title")
		pageData.Description = utils.Translate(pageData.Lang, "pages.register.meta.description")

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		templateData := RegisterPageData{
			PageData: *pageData,
		}

		// Validation
		if name == "" || password == "" {
			templateData.Error = "Name and password are required"
		} else if password != confirmPassword {
			templateData.Error = "Passwords do not match"
		} else {
			// Try to create the account
			if services == nil {
				log.Printf("Services not initialized")
				templateData.Error = "Unable to create account. Please try again."
			} else {
				// Hash the password
				passwordHash := utils.HashPassword(password)

				// Create the member
				member, err := services.Member.CreateMember(name, passwordHash)
				if err != nil {
					// Log the actual error for debugging
					log.Printf("Failed to create member account for name '%s': %v", name, err)

					// Show generic error message to user
					templateData.Error = "Unable to create account. Please try again."
				} else {
					// Account created successfully, sign in the user
					session, err := services.Session.CreateSession(member.ID)
					if err != nil {
						log.Printf("Failed to create session for new member %d: %v", member.ID, err)
						templateData.Error = "Account created but unable to sign in. Please try signing in manually."
					} else {
						// Set session cookie and redirect to homepage
						utils.SetSessionCookie(w, session.ID)
						log.Printf("Member %s registered and automatically signed in", member.Name)
						http.Redirect(w, r, "/", http.StatusSeeOther)
						return
					}
				}
			}
		}

		tmplClone, err := tmpl.Clone()
		if err != nil {
			http.Error(w, "Error cloning template", http.StatusInternalServerError)
			log.Printf("Error cloning template for register: %v", err)
			return
		}

		if err := tmplClone.ExecuteTemplate(w, "page-register", templateData); err != nil {
			http.Error(w, "Error rendering page-register template", http.StatusInternalServerError)
			log.Printf("Error rendering page-register template: %v", err)
		}
	}
}
