package routes

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sb-luis/creative-coding-bookclub/internal/routes/handlers"
	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// Helper function to prepare template data and handle theme
func preparePageData(r *http.Request, w http.ResponseWriter, currentLang string, services *services.Services) *utils.PageData {
	theme := utils.GetResolvedTheme(r)
	utils.SetThemeCookie(w, theme)

	pageData := utils.GetDefaultPageData(r.URL.Path, currentLang, theme, r.RequestURI)
	pageData.SupportedLanguages = utils.GetSupportedLanguages()

	// Check authentication status
	if sessionID, err := utils.GetSessionFromRequest(r); err == nil {
		if memberID, err := services.Session.GetMemberIDFromSession(sessionID); err == nil {
			if member, err := services.Member.GetMemberByID(memberID); err == nil {
				pageData.IsAuthenticated = true
				pageData.MemberName = member.Name
			} else {
				pageData.IsAuthenticated = false
				pageData.MemberName = ""
			}
		} else {
			pageData.IsAuthenticated = false
			pageData.MemberName = ""
		}
	} else {
		pageData.IsAuthenticated = false
		pageData.MemberName = ""
	}

	return pageData
}

// apiMiddleware sets appropriate headers for API endpoints
func apiMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Route-Type", "api")
		handler(w, r)
	}
}

// authMiddleware ensures that a request is authenticated
func authMiddleware(handler http.HandlerFunc, services *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Route-Type", "api")

		// Check authentication
		sessionID, err := utils.GetSessionFromRequest(r)
		if err != nil {
			http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
			return
		}

		memberID, err := services.Session.GetMemberIDFromSession(sessionID)
		if err != nil {
			http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
			return
		}

		// Store authenticated member ID in request context
		ctx := context.WithValue(r.Context(), "authenticated_member_id", memberID)
		handler(w, r.WithContext(ctx))
	}
}

// webMiddleware sets appropriate headers for HTML page endpoints
func webMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Route-Type", "web")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		handler(w, r)
	}
}

// renderNotFound renders the custom 404 page.
func renderNotFound(w http.ResponseWriter, r *http.Request, masterTmpl *template.Template, pageData *utils.PageData) {
	handlers.NotFoundHandler(w, r, masterTmpl, pageData)
}

// RegisterRoutes registers all the route handlers to the provided custom Router.
func RegisterRoutes(router *utils.Router, services *services.Services) {
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}

	htmlTemplates, err := utils.GetTemplateFiles()
	if err != nil {
		log.Fatalf("Error retrieving template files: %v", err)
	}

	masterTmpl, err := template.New("").Funcs(template.FuncMap{
		"i18nText": utils.Translate,
		"i18nHtml": func(lang string, key string, args ...interface{}) template.HTML {
			translatedText := utils.Translate(lang, key, args...)
			return template.HTML(utils.InnerMarkToHTML(translatedText))
		},
	}).ParseFiles(htmlTemplates...)
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Serve static files from "web/assets" under "/assets/"
	staticAssetsDir := filepath.Join(baseDir, "web/assets")

	router.PathPrefix("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(staticAssetsDir))))

	// =============================================================================
	// API ROUTES - Backend data endpoints
	// =============================================================================

	// Public Member API endpoints
	router.HandleFunc("/api/members", apiMiddleware(handlers.GetMembersHandler(services)), "GET")
	router.HandleFunc("/api/members/me", authMiddleware(handlers.GetCurrentMemberHandler(services), services), "GET")

	// Public Preference API endpoints
	router.HandleFunc("/api/preferences/theme", apiMiddleware(handlers.ThemePreferencesPostHandler), "POST")
	router.HandleFunc("/api/preferences/locale", apiMiddleware(handlers.LocalePreferencesPostHandler), "POST")

	// Public Authentication API endpoints
	router.HandleFunc("/api/auth/logout", apiMiddleware(handlers.LogoutHandler(services)), "POST")
	router.HandleFunc("/api/auth/sign-out", apiMiddleware(handlers.SignOutHandler(services)), "GET")
	router.HandleFunc("/api/auth/update-password", apiMiddleware(handlers.UpdatePasswordHandler(services)), "POST")
	router.HandleFunc("/api/auth/register", apiMiddleware(func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang, services)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for register: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.RegisterPostHandler(services)(w, r, tmpl, pageData)
	}), "POST")
	router.HandleFunc("/api/auth/sign-in", apiMiddleware(func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang, services)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for sign-in: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.SignInPostHandler(services)(w, r, tmpl, pageData)
	}), "POST")

	// Public Sketch API endpoints
	router.HandleFunc("/api/sketches/{memberName}", apiMiddleware(handlers.GetMemberSketchesHandler(services)), "GET")
	router.HandleFunc("/api/sketches/{memberName}/{sketchSlug}", apiMiddleware(handlers.SketchCodeHandler(services)), "GET")

	// Protected Sketch API endpoints (require authentication)
	router.HandleFunc("/api/sketches/{memberName}/{sketchSlug}", authMiddleware(handlers.CreateSketchHandler(services), services), "POST")
	router.HandleFunc("/api/sketches/{memberName}/{sketchSlug}", authMiddleware(handlers.UpdateSketchHandler(services), services), "PUT")           // source code only
	router.HandleFunc("/api/sketches/{memberName}/{sketchSlug}", authMiddleware(handlers.UpdateSketchMetadataHandler(services), services), "PATCH") // metadata only
	router.HandleFunc("/api/sketches/{memberName}/{sketchSlug}", authMiddleware(handlers.DeleteSketchHandler(services), services), "DELETE")

	// =============================================================================
	// WEB ROUTES - Frontend HTML page rendering
	// =============================================================================

	// Register page
	router.HandleFunc("/register", webMiddleware(func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang, services)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for register: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.RegisterGetHandler(w, r, tmpl, pageData)
	}), "GET")

	// Sign-in page
	router.HandleFunc("/sign-in", webMiddleware(func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang, services)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for sign-in: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.SignInGetHandler(w, r, tmpl, pageData)
	}), "GET")

	// Member's profile (requires authentication)
	router.HandleFunc("/me", webMiddleware(func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang, services)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for me: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.ProfileHandler(services)(w, r, tmpl, pageData)
	}), "GET")

	// Clean sketch view page (for viewing only, no editor)
	router.HandleFunc("/sketches/{memberName}/{sketchSlug}", webMiddleware(func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang, services)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for sketch view: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.SketchViewerPageHandler(services)(w, r, tmpl, pageData)
	}), "GET")

	// Sketch iframe content (for sandboxed execution)
	router.HandleFunc("/sketches/{memberName}/{sketchSlug}/iframe", webMiddleware(func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang, services)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for sketch iframe: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.SketchIframeContentHandler(services)(w, r, tmpl, pageData)
	}), "GET")

	// Member's sketch page (with editor capabilities)
	router.HandleFunc("/sketches/{memberName}/{sketchSlug}/edit", webMiddleware(func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang, services)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for sketch page: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.SketchEditorPageHandler(services)(w, r, tmpl, pageData)
	}), "GET")

	// Sketch lister page
	router.HandleFunc("/sketches", webMiddleware(func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang, services)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for sketch lister: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.SketchListerPageHandler(services)(w, r, tmpl, pageData)
	}), "GET")

	// Sketch Manager page (requires authentication)
	router.HandleFunc("/sketch-manager", webMiddleware(func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang, services)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for sketch manager: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.SketchManagerPageHandler(services)(w, r, tmpl, pageData)
	}), "GET")

	// Empty iframe page for iframe initialization
	router.HandleFunc("/empty-iframe", webMiddleware(func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang, services)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for empty iframe: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.EmptyIframeHandler(w, r, tmpl, pageData)
	}), "GET")

	// Homepage
	router.HandleFunc("/", webMiddleware(func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang, services)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning template for /: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.HomePageGetHandler(w, r, tmpl, pageData)
	}), "GET")

	// Set the NotFoundHandler on the router
	router.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		theme := utils.GetResolvedTheme(r)
		utils.SetThemeCookie(w, theme) // Ensure cookie is set

		pageData := utils.GetDefaultPageData(r.URL.Path, currentLang, theme, r.RequestURI)
		pageData.SupportedLanguages = utils.GetSupportedLanguages()
		renderNotFound(w, r, masterTmpl, pageData)
	}
}
