package routes

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sb-luis/creative-coding-bookclub/internal/routes/handlers"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// Helper function to prepare template data and handle theme
func preparePageData(r *http.Request, w http.ResponseWriter, currentLang string) *utils.PageData {
	theme := utils.GetResolvedTheme(r)
	utils.SetThemeCookie(w, theme)

	pageData := utils.GetDefaultPageData(r.URL.Path, currentLang, theme, r.RequestURI)
	pageData.SupportedLanguages = utils.GetSupportedLanguages()

	// Check authentication status
	if member, err := utils.GetCurrentMember(r); err == nil {
		pageData.IsAuthenticated = true
		pageData.MemberName = member.Name
	} else {
		pageData.IsAuthenticated = false
		pageData.MemberName = ""
	}

	return pageData
}

// renderNotFound renders the custom 404 page.
func renderNotFound(w http.ResponseWriter, r *http.Request, masterTmpl *template.Template, pageData *utils.PageData) {
	handlers.NotFoundHandler(w, r, masterTmpl, pageData)
}

// RegisterRoutes registers all the route handlers to the provided custom Router.
func RegisterRoutes(router *utils.Router) {
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

	// Preference handlers
	router.HandleFunc("/preferences/theme", handlers.ThemePreferencesPostHandler, "POST")
	router.HandleFunc("/preferences/locale", handlers.LocalePreferencesPostHandler, "POST")

	// Authentication routes - GET handlers
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for register: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.RegisterGetHandler(w, r, tmpl, pageData)
	}, "GET")

	router.HandleFunc("/sign-in", func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for sign-in: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.SignInGetHandler(w, r, tmpl, pageData)
	}, "GET")

	// Authentication routes - POST handlers
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for register: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.RegisterPostHandler(w, r, tmpl, pageData)
	}, "POST")

	router.HandleFunc("/sign-in", func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for sign-in: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.SignInPostHandler(w, r, tmpl, pageData)
	}, "POST")

	// Logout route
	router.HandleFunc("/logout", handlers.LogoutHandler, "POST")

	// Member profile route (authenticated only)
	router.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for me: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.ProfileHandler(w, r, tmpl, pageData)
	}, "GET")

	// Sign out route (authenticated only)
	router.HandleFunc("/sign-out", handlers.SignOutHandler, "GET")

	// Member's sketch page
	router.HandleFunc("/members/{memberName}/{sketchName}", func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for sketch page: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.SketchPageGetHandler(w, r, tmpl, pageData)
	}, "GET")

	// Route for the sketch lister page
	router.HandleFunc("/sketches", func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning master template for sketch lister: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.SketchListerPageHandler(w, r, tmpl, pageData)
	}, "GET")

	// Homepage
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		currentLang := utils.GetCurrentLanguage(r)
		pageData := preparePageData(r, w, currentLang)
		tmpl, err := masterTmpl.Clone()
		if err != nil {
			log.Printf("Error cloning template for /: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		handlers.HomePageGetHandler(w, r, tmpl, pageData)
	}, "GET")

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
