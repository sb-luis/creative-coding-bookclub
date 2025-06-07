package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// EmptyIframeHandler handles requests for the minimal empty iframe page
func EmptyIframeHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	pageData.Title = utils.Translate(pageData.Lang, "pages.emptyIframe.meta.title")
	pageData.Description = utils.Translate(pageData.Lang, "pages.emptyIframe.meta.description")

	tmplClone, err := tmpl.Clone()
	if err != nil {
		http.Error(w, "Error cloning template for empty iframe", http.StatusInternalServerError)
		log.Printf("Error cloning template for empty iframe: %v", err)
		return
	}

	if err := tmplClone.ExecuteTemplate(w, "page-empty-iframe", pageData); err != nil {
		http.Error(w, "Error rendering page-empty-iframe template", http.StatusInternalServerError)
		log.Printf("Error rendering page-empty-iframe template: %v", err)
	}
}
