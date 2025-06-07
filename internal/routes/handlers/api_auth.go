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
func RegisterPostHandler(services *services.Services) func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	return func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
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
			if services == nil {
				log.Printf("Services not initialized")
				templateData.Error = "Unable to create account. Please try again."
			} else {
				// Hash the password
				passwordHash := utils.HashPassword(password)

				// Create the member
				_, err := services.Member.CreateMember(name, passwordHash)
				if err != nil {
					// Log the actual error for debugging
					log.Printf("Failed to create member account for name '%s': %v", name, err)

					// Show generic error message to user
					templateData.Error = "Unable to create account. Please try again."
				} else {
					templateData.Success = "Account created successfully! Please sign in."
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
func SignInPostHandler(services *services.Services) func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	return func(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
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

// LogoutHandler handles member logout
func LogoutHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := utils.GetSessionFromRequest(r)
		if err == nil {
			// Delete the session from database
			if services != nil {
				if err := services.Session.DeleteSession(sessionID); err != nil {
					log.Printf("Error logging out member: %v", err)
				}
			}
		}

		// Clear the session cookie
		utils.ClearSessionCookie(w)

		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
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
}

// SignOutHandler handles member sign out (GET request)
func SignOutHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow authenticated members to sign out
		sessionID, err := utils.GetSessionFromRequest(r)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Delete the session from database
		if services != nil {
			if err := services.Session.DeleteSession(sessionID); err != nil {
				log.Printf("Error signing out member: %v", err)
			}
		}

		// Clear the session cookie
		utils.ClearSessionCookie(w)

		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// UpdatePasswordHandler handles password update requests for verified members
func UpdatePasswordHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get the current authenticated member
		sessionID, err := utils.GetSessionFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		memberID, err := services.Session.GetMemberIDFromSession(sessionID)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		currentPassword := r.FormValue("current_password")
		newPassword := r.FormValue("new_password")
		confirmPassword := r.FormValue("confirm_password")

		// Validate input
		if currentPassword == "" || newPassword == "" || confirmPassword == "" {
			http.Error(w, "All password fields are required", http.StatusBadRequest)
			return
		}

		if newPassword != confirmPassword {
			http.Error(w, "New passwords do not match", http.StatusBadRequest)
			return
		}

		if len(newPassword) < 6 {
			http.Error(w, "New password must be at least 6 characters long", http.StatusBadRequest)
			return
		}

		// Hash passwords
		currentPasswordHash := utils.HashPassword(currentPassword)
		newPasswordHash := utils.HashPassword(newPassword)

		// Update password using the service
		err = services.Member.UpdatePassword(memberID, currentPasswordHash, newPasswordHash)
		if err != nil {
			log.Printf("Error updating password for member ID %d: %v", memberID, err)

			// Return appropriate error messages
			switch err.Error() {
			case "only verified members can update their password":
				http.Error(w, "Only verified members can update their password", http.StatusForbidden)
			case "current password is incorrect":
				http.Error(w, "Current password is incorrect", http.StatusBadRequest)
			default:
				http.Error(w, "Failed to update password. Please try again.", http.StatusInternalServerError)
			}
			return
		}

		// Success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "message": "Password updated successfully"}`))
	}
}
