package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

// routes returns a http.Handler after passing requests and responses through middleware
func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	// Update the patter for the router for the static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Create the routes using the appropriate methods, patterns, and handlers
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Create a middleware chain containing our 'standard' middleware which will bre used for every request the app receives
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Return the 'standard' middleware chain with the router
	return standard.Then(router)
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
