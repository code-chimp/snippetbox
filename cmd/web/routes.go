package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a new middleware chain for dynamic requests
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Serve static files from the "./ui/static/" directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.getSnippets))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.getSnippet))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.getSnippetForm))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.postSnippetForm))

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.getSignupForm))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.postSignupForm))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.getLoginForm))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.postLoginForm))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(app.postLogout))

	// Create a base middleware chain for all requests
	baseMiddlewares := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Return the middleware chain with the ServeMux as the final handler
	return baseMiddlewares.Then(mux)
}
