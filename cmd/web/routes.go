package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.getSnippets))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.getSnippet))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.getSnippetForm))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.postSnippetForm))

	baseMiddlewares := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return baseMiddlewares.Then(mux)
}
