package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LanguageInfo holds information about a supported language.
type LanguageInfo struct {
	Code       string
	NameNative string
	NameEn     string
}

// translations stores the unmarshalled JSON data, allowing for nested structures.
var translations = make(map[string]interface{})

// supportedLanguagesInfo stores LanguageInfo structs
var supportedLanguagesInfo = []LanguageInfo{
	{Code: "en", NameNative: "English", NameEn: "English"},
	// {Code: "es", NameNative: "Español", NameEn: "Spanish"},
}
var defaultLanguage = "en"

const languageCookieName = "locale"

// I18nInit initializes the translation strings by loading them from JSON files.
func I18nInit() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}

	// Extract just the codes for loading files
	var langCodes []string
	for _, langInfo := range supportedLanguagesInfo {
		langCodes = append(langCodes, langInfo.Code)
	}

	for _, lang := range langCodes {
		filePath := filepath.Join(wd, "web", "locales", lang+".json")
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Warning: Could not read translation file %s: %v. Skipping language.", filePath, err)
			continue
		}

		var langTranslations interface{}
		if err := json.Unmarshal(fileBytes, &langTranslations); err != nil {
			log.Printf("Warning: Could not parse translation file %s: %v. Skipping language.", filePath, err)
			continue
		}
		translations[lang] = langTranslations
		log.Printf("Successfully loaded translations for language: %s from %s", lang, filePath)
	}

	if len(translations) == 0 {
		log.Fatalf("No translations loaded. Please check 'web/locales/' directory and file contents.")
	}
	// Ensure default language has translations loaded
	if _, ok := translations[defaultLanguage]; !ok {
		log.Fatalf("Default language '%s' translations not found. Please ensure '%s.json' exists and is valid.", defaultLanguage, defaultLanguage)
	}
}

// GetDefaultLanguage returns the default language code.
func GetDefaultLanguage() string {
	return defaultLanguage
}

// Translate returns the translated string for the given language and key (e.g., "parent.child.key").
// It uses Sprintf for basic variable substitution if args are provided.
func Translate(lang, key string, args ...interface{}) string {
	originalKey := key

	getValue := func(targetLang string) (string, bool) {
		if langData, ok := translations[targetLang]; ok {
			parts := strings.Split(key, ".")
			current := langData
			for i, part := range parts {
				if m, ok := current.(map[string]interface{}); ok {
					if val, found := m[part]; found {
						current = val
						if i == len(parts)-1 { // Last part
							if strVal, okStr := current.(string); okStr {
								return strVal, true
							}
							// Found a value, but it's not a string (e.g. a sub-map because key was too short)
							return "", false
						}
					} else {
						return "", false // Part of the path not found
					}
				} else {
					return "", false // Current element is not a map, cannot go deeper
				}
			}
		}
		return "", false
	}

	// Try the requested language
	if value, found := getValue(lang); found {
		if len(args) > 0 {
			return fmt.Sprintf(value, args...)
		}
		return value
	}

	// Try the default language if different from the requested one
	if lang != defaultLanguage {
		if value, found := getValue(defaultLanguage); found {
			log.Printf("Warning: Translation for key '%s' not found in '%s'. Using default '%s'.", originalKey, lang, defaultLanguage)
			if len(args) > 0 {
				return fmt.Sprintf(value, args...)
			}
			return value
		}
	}

	// Fallback to the original key if value not found in both requested and default languages
	log.Printf("Warning: Translation for key '%s' not found in '%s' or default '%s'. Returning key itself.", originalKey, lang, defaultLanguage)
	return originalKey
}

// GetCurrentLanguage determines the language from cookie, then Accept-Language header, then default.
func GetCurrentLanguage(r *http.Request) string {
	// Check for language cookie
	cookie, err := r.Cookie(languageCookieName)
	if err == nil {
		normalizedCookieLang := strings.ToLower(cookie.Value)
		for _, supportedLangInfo := range supportedLanguagesInfo {
			if normalizedCookieLang == supportedLangInfo.Code {
				return supportedLangInfo.Code
			}
		}
		log.Printf("Warning: Invalid language code '%s' in cookie. Ignoring.", normalizedCookieLang)
	}

	// Fallback to Accept-Language header (simplified check)
	acceptLanguage := r.Header.Get("Accept-Language")
	if acceptLanguage != "" {
		langs := strings.Split(acceptLanguage, ",")
		for _, langTag := range langs {
			mainLang := strings.Split(strings.TrimSpace(langTag), ";")[0]
			normalizedAcceptLang := strings.ToLower(mainLang)
			// Check for full match (e.g., "en-US" -> "en") or prefix match
			for _, supportedLangInfo := range supportedLanguagesInfo {
				if strings.HasPrefix(normalizedAcceptLang, supportedLangInfo.Code) || normalizedAcceptLang == supportedLangInfo.Code {
					return supportedLangInfo.Code
				}
			}
		}
	}

	// Fallback to default language
	return defaultLanguage
}

// GetSupportedLanguages returns the list of supported language details.
func GetSupportedLanguages() []LanguageInfo {
	return supportedLanguagesInfo
}

// SetLanguageCookie sets the language preference in a cookie.
func SetLanguageCookie(w http.ResponseWriter, langCode string) {
	validLang := false
	for _, supportedLang := range supportedLanguagesInfo {
		if langCode == supportedLang.Code {
			validLang = true
			break
		}
	}

	if !validLang {
		log.Printf("Warning: Attempted to set invalid language cookie: %s. Defaulting to %s.", langCode, defaultLanguage)
		langCode = defaultLanguage
	}

	http.SetCookie(w, &http.Cookie{
		Name:     languageCookieName,
		Value:    langCode,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour), // 1 year
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true, // Assuming HTTPS
	})
	log.Printf("Language cookie set to: %s", langCode)
}
