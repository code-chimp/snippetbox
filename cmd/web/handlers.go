package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox.code-chimp.net/internal/models"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (app *application) getSnippets(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.gohtml", data)
}

func (app *application) getSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.gohtml", data)
}

type snippetcreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) getSnippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetcreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.gohtml", data)
}

func (app *application) postSnippetForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := snippetcreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: make(map[string]string),
	}

	// validation
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "Title is required"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "Title must be less than 100 characters"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "Content is required"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "Expiry must be 1, 7, or 365 days"
	}

	// send them back if we found errors
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusBadRequest, "create.gohtml", data)
		return
	}

	// persist
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
