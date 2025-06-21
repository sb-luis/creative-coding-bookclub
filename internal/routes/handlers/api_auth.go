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

// SignOutHandler handles member sign out and returns JSON response
func SignOutHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := utils.GetSessionFromRequest(r)
		if err == nil {
			// Delete the session from database
			if services != nil {
				if err := services.Session.DeleteSession(sessionID); err != nil {
					log.Printf("Error signing out member: %v", err)
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
