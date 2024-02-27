package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// create an application wide struct for logging
type application struct {
	logger *slog.Logger
}

func main() {
	// add better configuration management
	addr := flag.String("addr", ":4000", "HttP network address")
	flag.Parse()

	// add new logging functionality
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// initialize application instance of our struct with the dependencies
	app := &application{
		logger: logger,
	}

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

	// Print a log a message to say that the server is starting.
	logger.Info("starting server", "addr", *addr)

	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address
	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
