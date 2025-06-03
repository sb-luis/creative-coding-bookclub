package utils

import (
	"io/fs"
	"log"
	"path/filepath"
)

func GetTemplateFiles() ([]string, error) {
	// Collect all HTML templates
	var htmlTemplates []string

	// Get all HTML templates from 'web/pages'
	err := filepath.WalkDir("web/pages", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Add only .html files to the list
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			htmlTemplates = append(htmlTemplates, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through 'web/pages/_static': %v", err)
	}

	// Get all HTML templates from 'web/components'
	err = filepath.WalkDir("web/components", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Add only .html files to the list
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			htmlTemplates = append(htmlTemplates, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through 'web/components': %v", err)
	}

	return htmlTemplates, nil
}
