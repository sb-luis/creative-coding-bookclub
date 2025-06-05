package handlers

import (
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// ThemePreferencesPostHandler handles POST requests to update theme preferences.
func ThemePreferencesPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	theme := r.FormValue("theme")
	redirectURL := r.FormValue("redirectURL")

	// Validate theme value
	validThemes := map[string]bool{"light": true, "dark": true, "system": true}
	if !validThemes[theme] {
		log.Printf("Invalid theme value received: %s. Defaulting to 'system'.", theme)
		theme = "system" // Default to system if invalid
	}

	utils.SetThemeCookie(w, theme)

	if redirectURL == "" {
		log.Printf("RedirectURL not provided. Defaulting to '/'")
		redirectURL = "/" // Default redirect
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// LocalePreferencesPostHandler handles POST requests to update locale preferences.
func LocalePreferencesPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	locale := r.FormValue("locale")
	redirectURL := r.FormValue("redirectURL")

	utils.SetLanguageCookie(w, locale)

	if redirectURL == "" {
		log.Printf("RedirectURL not provided for locale change. Defaulting to '/'")
		redirectURL = "/" // Default redirect
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
