package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) getSnippets(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.gohtml",
		"./ui/html/partials/nav.gohtml",
		"./ui/html/pages/getSnippets.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) getSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display snippet %d...", id)
}

func (app *application) getSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET Snippet Form"))
}

func (app *application) postSnippetForm(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("POST Snippet Form"))
}
