package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/services"
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
