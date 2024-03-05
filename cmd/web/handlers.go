package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"snippetbox.ajigherighe.net/internal/models"
	"strconv"
)

// Define a home handler function
// Further define it as a method against *application type from main
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path that exactly matches "/".
	/*if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	*/

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Call newTemplateData to add the current year to data
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Use the render helper function.
	app.render(w, r, http.StatusOK, "home.tmpl.html", data)

	/*
		files := []string{
			"./ui/html/base.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
			"./ui/html/pages/home.tmpl.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		// Create a templateData instance
		data := templateData{
			Snippets: snippets,
		}

		// Pass in the templateData struct when executing template
		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		for _, snippet := range snippets {
			fmt.Fprintf(w, "%+v\n", snippet)
		}*/
}

// Add a snippetView handler function
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	/*id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // not found helper function
		return
	}
	*/

	// When httprouter is parsing a request, the values of any named parameters will be stored
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
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

	// Repeat process from home here
	data := app.newTemplateData(r)
	data.Snippet = snippet

	// User the new render helper.
	app.render(w, r, http.StatusOK, "view.tmpl.html", data)

}

// Add a snippetCreate handler function.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Allow form to create a snippet
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
	// w.Write([]byte("Display the form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Call r.ParseForm() to add any data in POST request body to the r.PostForm map. This works for PUT and PATCH
	// requests.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use r.PostForm.Get to grab desired information from r.PostForm map
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// r.PostForm returns strings but expiration data are numbers. So, must convert.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Insert into the database using SnippetModel.Insert()
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Redirect to the relevant snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
