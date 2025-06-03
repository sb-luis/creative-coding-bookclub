package utils

import (
	"os"
	"strings"
)

const (
	ProdBaseURL = "https://creativecodingbook.club"
	DevBaseURL  = "http://localhost:8000"
)

// GetBaseURL returns the base URL based on the APP_ENV environment variable.
// It defaults to DevBaseURL if APP_ENV is not "production".
func GetBaseURL() string {
	appEnv := strings.ToLower(os.Getenv("APP_ENV"))
	if appEnv == "production" {
		return ProdBaseURL
	}
	return DevBaseURL
}

// GetFullURL constructs a full URL by joining the base URL with a given path.
// It ensures that there's exactly one slash between the base URL and the path.
func GetFullURL(path string) string {
	baseURL := GetBaseURL()
	// Ensure path starts with a slash if it's not empty
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	// Remove trailing slash from baseURL if present, to avoid double slashes
	// unless the path is empty or just "/"
	if strings.HasSuffix(baseURL, "/") && (path != "" && path != "/") {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}

	// If path is "/" and baseURL ends with "/", avoid double slash
	if path == "/" && strings.HasSuffix(baseURL, "/") {
		return baseURL
	}

	return baseURL + path
}
