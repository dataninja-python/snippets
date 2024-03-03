package main

import (
	"html/template"
	"path/filepath"
	"snippetbox.ajigherighe.net/internal/models"
)

// The template path constants.
const baseTemplatePath string = "./ui/html/base.tmpl.html"
const allPartialsPath string = "./ui/html/partials/*.tmpl.html"
const allPagesPath string = "./ui/html/pages/*.tmpl.html"

type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map that caches template data
	cache := map[string]*template.Template{}

	// Use filepath.Glob() function to get slice of all filepaths that match the pattern
	// "./ui/html/pages/*.tmpl.html".
	pages, err := filepath.Glob(allPagesPath)
	if err != nil {
		return nil, err
	}

	// Loop through page filepaths
	for _, page := range pages {
		// extract the file name (home.tmpl.html) from the full path
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		ts, err := template.ParseFiles(baseTemplatePath)
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob(allPartialsPath)
		if err != nil {
			return nil, err
		}

		// Call ParseFiles *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map as normal
		cache[name] = ts

		/*
			// Create a slice of base template, partials, and the page
			files := []string{
				"./ui/html/base.tmpl.html",
				"./ui/html/partials/nav.tmpl.html",
				page,
			}

			// Parse the files into a template set.
			ts, err := template.ParseFiles(files...)
			if err != nil {
				return nil, err
			}

			// Add the template set to the map, using the name of the page
			// (like 'home.tmpl.html') as the key.
			cache[name] = ts
		*/
	}

	// Return the map.
	return cache, nil
}
