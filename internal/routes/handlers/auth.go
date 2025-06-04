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
	Error   string
	Success string
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

	if err := tmplClone.ExecuteTemplate(w, "page-register", templateData); err != nil {
		http.Error(w, "Error rendering page-register template", http.StatusInternalServerError)
		log.Printf("Error rendering page-register template: %v", err)
	}
}

// RegisterPostHandler processes the registration form
func RegisterPostHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	pageData.Title = "Register"
	pageData.Description = "Create a new member account"

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
		_, err := utils.CreateMemberAccount(name, password)
		if err != nil {
			// Log the actual error for debugging
			log.Printf("Failed to create member account for name '%s': %v", name, err)

			// Show generic error message to user
			templateData.Error = "Unable to create account. Please try again."
		} else {
			templateData.Success = "Account created successfully! Please sign in."
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

	if err := tmplClone.ExecuteTemplate(w, "page-sign-in", templateData); err != nil {
		http.Error(w, "Error rendering page-sign-in template", http.StatusInternalServerError)
		log.Printf("Error rendering page-sign-in template: %v", err)
	}
}

// SignInPostHandler processes the sign-in form
func SignInPostHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	pageData.Title = "Sign in"
	pageData.Description = "Sign in to your account"

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
		member, err := utils.AuthenticateMember(name, password)
		if err != nil {
			templateData.Error = "Invalid name or password"
		} else {
			// Create session
			session, err := utils.CreateMemberSession(member.ID)
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

// LogoutHandler handles member logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := utils.GetSessionFromRequest(r)
	if err == nil {
		// Delete the session from database
		if err := utils.LogoutMember(sessionID); err != nil {
			log.Printf("Error logging out member: %v", err)
		}
	}

	// Clear the session cookie
	utils.ClearSessionCookie(w)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ProfileHandler shows the authenticated member's profile
func ProfileHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	// Get the current member
	member, err := utils.GetCurrentMember(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	pageData.Title = "My Profile"
	pageData.Description = "View your member profile"

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

// SignOutHandler handles member sign out (GET request)
func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow authenticated members to sign out
	sessionID, err := utils.GetSessionFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Delete the session from database
	if err := utils.LogoutMember(sessionID); err != nil {
		log.Printf("Error signing out member: %v", err)
	}

	// Clear the session cookie
	utils.ClearSessionCookie(w)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
