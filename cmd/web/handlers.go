package main

import (
	"errors"
	"fmt"
	"github.com/code-chimp/snippetbox/internal/models"
	"github.com/code-chimp/snippetbox/internal/validator"
	"net/http"
	"strconv"
)

// getSnippets displays the most recent snippets.
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

// getSnippet displays a specific snippet based on its ID.
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

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

// getSnippetForm displays the snippet form.
func (app *application) getSnippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.gohtml", data)
}

// postSnippetForm handles the submission of the snippet form.
func (app *application) postSnippetForm(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validation
	form.CheckField(validator.NotBlank(form.Title), "title", "Title is required")
	form.CheckField(validator.MaxLength(form.Title, 100), "title", "Title cannot be longer than 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "Content is required")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "Expiry must be 1, 7, or 365 days")

	// send them back if we found errors
	if !form.Valid() {
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

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
