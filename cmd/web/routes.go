package main

import (
	"github.com/justinas/alice"
	"net/http"
)

// routes returns a http.Handler after passing requests and responses through middleware
func (app *application) routes() http.Handler {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern
	// fmt.Println("Hello, world!")
	mux := http.NewServeMux()
	// serve files from the static directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// handle routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Wrap the existing chain with the recoverPanic middleware.
	// Wrap the existing chain with the logRequest
	// Pass the mux to the middleware for added security to be added before returning it
	// Pre-using justinas/alice package
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	// Create a middleware chain containing our 'standard' middleware which will bre used for every request the app receives
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Return the 'standard' middleware chain followed by the servemux
	return standard.Then(mux)
}

/*
func (app *application) routes() *http.ServeMux {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern
	// fmt.Println("Hello, world!")
	mux := http.NewServeMux()
	// serve files from the static directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// handle routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
*/
