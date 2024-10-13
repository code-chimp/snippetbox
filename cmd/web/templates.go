package main

import (
	"github.com/code-chimp/snippetbox/internal/models"
	"github.com/code-chimp/snippetbox/ui"
	"html/template"
	"io/fs"
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

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.go.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		// parse base template to create a new template set
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
