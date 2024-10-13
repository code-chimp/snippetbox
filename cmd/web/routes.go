package main

import (
	"github.com/code-chimp/snippetbox/ui"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a new middleware chain for dynamic requests
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	// Add the requireAuthentication middleware to the dynamic middleware chain
	protected := dynamic.Append(app.requireAuthentication)

	// Register the static file server
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// Snippets
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.getSnippets))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.getSnippet))
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.getSnippetForm))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.postSnippetForm))

	// User authentication
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.getSignupForm))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.postSignupForm))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.getLoginForm))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.postLoginForm))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.postLogout))

	// Create a base middleware chain for all requests
	baseMiddlewares := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Return the middleware chain with the ServeMux as the final handler
	return baseMiddlewares.Then(mux)
}
