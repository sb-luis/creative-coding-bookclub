package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sb-luis/creative-coding-bookclub/internal/model"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

// SketchListerPageData holds all data for the sketch lister page template.
type SketchListerPageData struct {
	utils.PageData
	Members []model.MemberSketchInfo
}

// SketchListerPageHandler handles requests to display the list of all sketches.
func SketchListerPageHandler(w http.ResponseWriter, r *http.Request, tmpl *template.Template, pageData *utils.PageData) {
	baseSketchesPath := filepath.Join("web", "assets", "js", "sketches")
	var membersData []model.MemberSketchInfo

	// Override metadata fields specific to the sketch lister page
	pageData.Title = utils.Translate(pageData.Lang, "pages.sketchLister.meta.title")
	pageData.Description = utils.Translate(pageData.Lang, "pages.sketchLister.meta.description")
	pageData.Keywords = utils.Translate(pageData.Lang, "pages.sketchLister.meta.keywords")

	memberDirs, err := os.ReadDir(baseSketchesPath)
	if err != nil {
		log.Printf("Error reading sketches directory %s: %v", baseSketchesPath, err)
		http.Error(w, "Failed to list sketches", http.StatusInternalServerError)
		return
	}

	for _, memberDir := range memberDirs {
		if memberDir.IsDir() && !strings.HasPrefix(memberDir.Name(), "_") { // Ignore directories like _example
			memberName := memberDir.Name()
			memberPath := filepath.Join(baseSketchesPath, memberName)
			var sketches []model.SketchInfo

			sketchFiles, err := os.ReadDir(memberPath)
			if err != nil {
				log.Printf("Error reading member sketch directory %s: %v", memberPath, err)
				continue // Skip this member on error
			}

			for _, sketchFile := range sketchFiles {
				if !sketchFile.IsDir() && strings.HasSuffix(sketchFile.Name(), ".json") {
					sketchBaseName := strings.TrimSuffix(sketchFile.Name(), ".json")

					// Check for the corresponding JavaScript file
					jsFileName := sketchBaseName + ".js"
					jsFilePath := filepath.Join(memberPath, jsFileName)
					jsFileInfo, errJs := os.Stat(jsFilePath)

					if errJs != nil {
						if os.IsNotExist(errJs) {
							log.Printf("Sketch Lister: JS file %s not found for sketch %s by member %s. Skipping.", jsFilePath, sketchBaseName, memberName)
						} else {
							log.Printf("Sketch Lister: Error stating JS file %s for sketch %s by member %s: %v. Skipping.", jsFilePath, sketchBaseName, memberName, errJs)
						}
						continue // Skip this sketch
					}
					if jsFileInfo.IsDir() {
						log.Printf("Sketch Lister: JS path %s for sketch %s by member %s is a directory, not a file. Skipping.", jsFilePath, sketchBaseName, memberName)
						continue // Skip this sketch
					}

					// Proceed if JS file exists and is a file
					sketchURL := filepath.Join("/members", memberName, sketchBaseName)
					sketchURL = strings.ReplaceAll(sketchURL, "\\\\", "/") // Ensure forward slashes for URL

					var sketchInfoData model.SketchInfo
					sketchInfoData.Name = sketchBaseName
					sketchInfoData.URL = sketchURL
					sketchInfoData.Alias = memberName

					// Read metadata from JSON file
					metadataFilePath := filepath.Join(memberPath, sketchFile.Name())
					metadataBytes, err := os.ReadFile(metadataFilePath)
					if err != nil {
						log.Printf("Error reading sketch metadata file %s: %v", metadataFilePath, err)
					} else {
						if err := json.Unmarshal(metadataBytes, &sketchInfoData); err != nil {
							log.Printf("Error unmarshalling sketch JSON data into SketchInfo from %s: %v", metadataFilePath, err)
						}
					}

					sketches = append(sketches, sketchInfoData)
				}
			}

			if len(sketches) > 0 {
				membersData = append(membersData, model.MemberSketchInfo{Name: memberName, Sketches: sketches})
			}
		}
	}

	// Prepare the data for the template
	templateData := SketchListerPageData{
		PageData: *pageData,
		Members:  membersData, // Renamed from MembersData
	}

	err = tmpl.ExecuteTemplate(w, "page-sketch-lister", templateData)
	if err != nil {
		log.Printf("Error executing page-sketch-lister template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
