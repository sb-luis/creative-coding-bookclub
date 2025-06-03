package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sb-luis/creative-coding-bookclub/internal/model"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// SketchPageData holds all data for the sketch page template.
type SketchPageData struct {
	utils.PageData
	MemberName   string
	SketchName   string
	SketchJsPath string
}

// SketchPageGetHandler handles requests to display a sketch page.
func SketchPageGetHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	memberName := utils.PathVariable(r, "memberName")
	sketchName := utils.PathVariable(r, "sketchName")

	baseSketchDir := filepath.Join("web", "assets", "js", "sketches")

	sourceFileName := sketchName + ".js"
	sketchJsonFileName := sketchName + ".json"
	sketchJsonFilePath := filepath.Join(baseSketchDir, memberName, sketchJsonFileName)
	actualSketchJsFilePath := filepath.Join(baseSketchDir, memberName, sourceFileName)

	// Check if the JavaScript file exists.
	fileInfoJs, errJs := os.Stat(actualSketchJsFilePath)
	if os.IsNotExist(errJs) || (errJs == nil && fileInfoJs.IsDir()) {
		if os.IsNotExist(errJs) {
			log.Printf("Sketch JavaScript file not found: %s", actualSketchJsFilePath)
		} else {
			log.Printf("Sketch JavaScript path is a directory, not a file: %s", actualSketchJsFilePath)
		}
		NotFoundHandler(w, r, tmpl, pageData) // Pass original pageData for 404
		return
	} else if errJs != nil {
		log.Printf("Error checking sketch JavaScript file %s: %v", actualSketchJsFilePath, errJs)
		http.Error(w, "Internal Server Error checking sketch JS file", http.StatusInternalServerError)
		return
	}

	// Read the associated JSON for metadata.
	jsonData, errJsonRead := os.ReadFile(sketchJsonFilePath)
	if errJsonRead == nil {
		var sketchInfoData model.SketchInfo
		if jsonErr := json.Unmarshal(jsonData, &sketchInfoData); jsonErr == nil {
			// Populate pageData with metadata from JSON
			if sketchInfoData.Title != nil && *sketchInfoData.Title != "" {
				pageData.Title = *sketchInfoData.Title
			}
			if sketchInfoData.Description != nil && *sketchInfoData.Description != "" {
				pageData.Description = *sketchInfoData.Description
			}
			if sketchInfoData.Keywords != nil && *sketchInfoData.Keywords != "" {
				pageData.Keywords = *sketchInfoData.Keywords
			}
		} else {
			log.Printf("Error unmarshalling sketch JSON data from %s: %v. Proceeding without JSON metadata.", sketchJsonFilePath, jsonErr)
		}
	} else {
		// Log if JSON is not found or unreadable. This is not a fatal error for the page.
		log.Printf("Sketch JSON file %s not found or unreadable (error: %v). Proceeding without JSON metadata.", sketchJsonFilePath, errJsonRead)
	}

	sketchJsPath := "/assets/js/sketches/" + memberName + "/" + sourceFileName

	templateData := SketchPageData{
		PageData:     *pageData,
		MemberName:   memberName,
		SketchName:   sketchName,
		SketchJsPath: sketchJsPath,
	}

	log.Printf("Rendering sketch page for member: %s, sketch: %s (JS file: %s, JSON data file: %s)",
		memberName, sketchName, templateData.SketchJsPath, sketchJsonFilePath)

	err := tmpl.ExecuteTemplate(w, "page-sketch", templateData)
	if err != nil {
		log.Printf("Error executing page-sketch template: %v", err)
		http.Error(w, "Internal Server Error executing template", http.StatusInternalServerError)
	}
}
