package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox.ajigherighe.net/internals/models"
	"strconv"
)

// Define a home handler function
// Further define it as a method against *application type from main
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path that exactly matches "/".
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// Initialize a slice containing the paths to the two files
	/*
		files := []string{
			"./ui/html/pages/base.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
			"./ui/html/pages/home.tmpl.html",
		}
	*/

	// use template home page
	/*
			ts, err := template.ParseFiles(files...)
			if err != nil {
				app.serverError(w, r, err) // server error helper
				return
			}
			err = ts.ExecuteTemplate(w, "base", nil)

		if err != nil {
			app.serverError(w, r, err) // server error helper
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	*/
	// w.Write([]byte("Hello from Snippets"))
}

// Add a snippetView handler function
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // not found helper function
		return
	}

	// Use SnippetModel's Get to retrieve data from the database
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", snippet)
}

// Add a snippetCreate handler function.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	/*if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		//w.WriteHeader(405)
		//w.Write([]byte("Method Not Allowed"))
		http.Error(w, "Method Not Allowed", 405)
		return
	}*/
	// this now uses constants and helper functions to be more idomatic Go code
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed) // client error helper
		return
	}

	// testing with dummy data
	title := "O snail"
	content := "O snail\nClimb Mount Fuji, \nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	// Insert into the database using SnippetModel.Insert()
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Redirect to the relevant snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)

	w.Write([]byte("Create a new snippet..."))
}
