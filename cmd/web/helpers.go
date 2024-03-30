package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError logs the error and sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends a specific status code and corresponding description to the client.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// render is a helper that renders a template with the base template and partials.
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, name string, data templateData) {
	ts, ok := app.templates[name]
	if !ok {
		app.serverError(w, r, fmt.Errorf("the template %s does not exist", name))
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
