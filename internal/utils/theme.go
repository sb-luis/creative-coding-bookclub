package utils

import (
	"net/http"
	"time"
)

const themeCookieName = "theme"

// GetResolvedTheme determines the theme based on cookies or system default.
// It returns the theme name ("light", "dark", "system").
func GetResolvedTheme(r *http.Request) string {
	validThemes := map[string]bool{"light": true, "dark": true, "system": true}

	cookie, err := r.Cookie(themeCookieName)
	if err == nil && validThemes[cookie.Value] {
		return cookie.Value // Theme from cookie
	}

	return "system" // Default theme
}

// SetThemeCookie sets or clears the theme cookie based on the provided theme.
// If theme is "system", the cookie is cleared. Otherwise, it's set.
func SetThemeCookie(w http.ResponseWriter, theme string) {
	validThemes := map[string]bool{"light": true, "dark": true, "system": true}
	if !validThemes[theme] {
		theme = "system" // Default to system if an invalid theme is passed
	}

	if theme == "system" {
		// Clear the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     themeCookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0), // Expire immediately
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   true, // Assuming HTTPS
		})
	} else {
		// Set the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     themeCookieName,
			Value:    theme,
			Path:     "/",
			Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 year
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   true, // Assuming HTTPS
		})
	}
}
