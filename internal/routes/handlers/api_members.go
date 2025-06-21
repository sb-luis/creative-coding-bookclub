package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// MemberResponse represents a member in API responses
type MemberResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetMembersHandler handles GET requests to return a list of all members (IDs and Names)
func GetMembersHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type for JSON response
		w.Header().Set("Content-Type", "application/json")

		if services == nil {
			log.Printf("Services not initialized")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Get all members from service
		members, err := services.Member.GetAllMembers()
		if err != nil {
			log.Printf("Error getting all members: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert to response format (exclude password hash and other sensitive data)
		var memberResponses []MemberResponse
		for _, member := range members {
			memberResponses = append(memberResponses, MemberResponse{
				ID:   member.ID,
				Name: member.Name,
			})
		}

		// Return JSON response
		if err := json.NewEncoder(w).Encode(memberResponses); err != nil {
			log.Printf("Error encoding JSON response for members: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Printf("Served %d members via API", len(memberResponses))
	}
}

// GetCurrentMemberHandler handles GET requests to return the current authenticated member
func GetCurrentMemberHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type for JSON response
		w.Header().Set("Content-Type", "application/json")

		// Get authenticated member ID from context
		memberID, ok := r.Context().Value("authenticated_member_id").(int)
		if !ok {
			http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
			return
		}

		// Get member details
		member, err := services.Member.GetMemberByID(memberID)
		if err != nil {
			log.Printf("Error getting member by ID %d: %v", memberID, err)
			http.Error(w, `{"error":"Member not found"}`, http.StatusNotFound)
			return
		}

		// Return member response
		response := MemberResponse{
			ID:   member.ID,
			Name: member.Name,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding current member response: %v", err)
			http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}

		log.Printf("Served current member information for member %d (%s)", member.ID, member.Name)
	}
}

// UpdatePasswordRequest represents the request body for password updates
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// UpdatePasswordHandler handles PATCH requests to update the authenticated member's password
func UpdatePasswordHandler(services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type for JSON response
		w.Header().Set("Content-Type", "application/json")

		// Get authenticated member ID from context
		memberID, ok := r.Context().Value("authenticated_member_id").(int)
		if !ok {
			http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
			return
		}

		// Parse request body
		var req UpdatePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error decoding password update request: %v", err)
			http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
			return
		}

		// Validate input
		if req.CurrentPassword == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
			http.Error(w, `{"error":"All password fields are required"}`, http.StatusBadRequest)
			return
		}

		if req.NewPassword != req.ConfirmPassword {
			http.Error(w, `{"error":"New passwords do not match"}`, http.StatusBadRequest)
			return
		}

		if len(req.NewPassword) < 6 {
			http.Error(w, `{"error":"New password must be at least 6 characters long"}`, http.StatusBadRequest)
			return
		}

		// Hash passwords
		currentPasswordHash := utils.HashPassword(req.CurrentPassword)
		newPasswordHash := utils.HashPassword(req.NewPassword)

		// Update password using the service
		err := services.Member.UpdatePassword(memberID, currentPasswordHash, newPasswordHash)
		if err != nil {
			log.Printf("Error updating password for member ID %d: %v", memberID, err)

			// Return appropriate error messages
			switch err.Error() {
			case "only verified members can update their password":
				http.Error(w, `{"error":"Only verified members can update their password"}`, http.StatusForbidden)
			case "current password is incorrect":
				http.Error(w, `{"error":"Current password is incorrect"}`, http.StatusBadRequest)
			default:
				http.Error(w, `{"error":"Failed to update password. Please try again."}`, http.StatusInternalServerError)
			}
			return
		}

		// Success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "message": "Password updated successfully"}`))

		log.Printf("Successfully updated password for member ID %d", memberID)
	}
}
