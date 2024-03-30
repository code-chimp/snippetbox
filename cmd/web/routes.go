package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.getSnippets)
	mux.HandleFunc("GET /snippet/view/{id}", app.getSnippet)
	mux.HandleFunc("GET /snippet/create", app.getSnippetForm)
	mux.HandleFunc("POST /snippet/create", app.postSnippetForm)

	return mux
}
