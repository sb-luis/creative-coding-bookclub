package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

func HomePageGetHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	tmplClone, err := tmpl.Clone()
	if err != nil {
		http.Error(w, "Error cloning template for index", http.StatusInternalServerError)
		log.Printf("Error cloning template for index: %v", err)
		return
	}

	if err := tmplClone.ExecuteTemplate(w, "page-index", pageData); err != nil {
		http.Error(w, "Error rendering page-index template", http.StatusInternalServerError)
		log.Printf("Error rendering page-index template: %v", err)
	}
}
