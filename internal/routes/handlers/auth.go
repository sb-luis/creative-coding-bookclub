package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// RegisterPageData holds data for the register page
type RegisterPageData struct {
	utils.PageData
	Error string
}

// SignInPageData holds data for the sign-in page
type SignInPageData struct {
	utils.PageData
	Error string
}

// ProfilePageData holds data for the profile page
type ProfilePageData struct {
	utils.PageData
	Name     string
	MemberID int
}

// RegisterGetHandler shows the registration form
func RegisterGetHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	pageData.Title = "Register"
	pageData.Description = "Create a new member account"

	templateData := RegisterPageData{
		PageData: *pageData,
	}

	tmplClone, err := tmpl.Clone()
	if err != nil {
		http.Error(w, "Error cloning template", http.StatusInternalServerError)
		log.Printf("Error cloning template for register: %v", err)
		return
	}

	if err := tmplClone.ExecuteTemplate(w, "register.html", templateData); err != nil {
		http.Error(w, "Error rendering register template", http.StatusInternalServerError)
		log.Printf("Error rendering register template: %v", err)
	}
}

// SignInGetHandler shows the sign-in form
func SignInGetHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	pageData.Title = "Sign in"
	pageData.Description = "Sign in to your account"

	templateData := SignInPageData{
		PageData: *pageData,
	}

	tmplClone, err := tmpl.Clone()
	if err != nil {
		http.Error(w, "Error cloning template", http.StatusInternalServerError)
		log.Printf("Error cloning template for sign-in: %v", err)
		return
	}

	if err := tmplClone.ExecuteTemplate(w, "sign-in.html", templateData); err != nil {
		http.Error(w, "Error rendering sign-in template", http.StatusInternalServerError)
		log.Printf("Error rendering sign-in template: %v", err)
	}
}
