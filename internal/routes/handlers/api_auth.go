package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// API authentication handlers
// These are pure API handlers that return JSON responses or perform redirects
// They do not render HTML templates directly

// LogoutHandler handles member logout (POST request)
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

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
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

		// Redirect to homepage
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
