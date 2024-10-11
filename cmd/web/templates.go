package main

import (
	"github.com/code-chimp/snippetbox/internal/models"
	"html/template"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// humanDate returns a human readable string representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// functions is a map of functions that can be used in templates.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// newTemplateCache creates a new template cache.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// parse base template to create a new template set
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.go.tmpl")
		if err != nil {
			return nil, err
		}

		// parse all partials and add them to the template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// finally parse the page template and add it to the template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
