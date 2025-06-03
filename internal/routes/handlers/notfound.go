package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	w.WriteHeader(http.StatusNotFound)

	pageData.Title = utils.Translate(pageData.Lang, "pages.notFound.meta.title")
	pageData.Description = utils.Translate(pageData.Lang, "pages.notFound.meta.description")
	pageData.Keywords = utils.Translate(pageData.Lang, "pages.notFound.meta.keywords")
	pageData.Icon = "🤷"

	// Clone the master template for this request.
	clonedMasterTmpl, err := tmpl.Clone()
	if err != nil {
		log.Printf("Error cloning template in NotFoundHandler: %v", err)
		http.Error(w, utils.Translate(pageData.Lang, "pages.notFound.title")+" - Error preparing page", http.StatusInternalServerError)
		return
	}

	// Execute the specific page template, passing pageData directly.
	err = clonedMasterTmpl.ExecuteTemplate(w, "not-found.html", pageData)
	if err != nil {
		log.Printf("Error executing not-found.html template: %v", err)
		// Fallback to a simpler error if the 404 template itself fails
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
