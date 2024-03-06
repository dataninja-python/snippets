package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"snippetbox.ajigherighe.net/internal/models"
	"snippetbox.ajigherighe.net/internal/validator"
	"strconv"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

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

	// Initialize a new createSnippetForm instance and pass it to the template.
	// TODO: set default values as it makes sense
	data.Form = snippetCreateForm{
		Expires: 365,
	}
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
	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")

	// r.PostForm returns strings but expiration data are numbers. So, must convert.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	// Because the Validator struct is embedded by the snippetCreateForm struct, we can call CheckField() directly on
	// it to execute our validation checks. CheckField() will add the provided key and error message to the FieldErrors
	// map if the check does not evaluate to true. For example, in the first line here we "check that the form.Title
	// field is not blank." In the second, we "check that the form.Title field has a maximum character length of 100"
	// and so on.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than "+
		"100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires",
		"This field must equal 1, 7, or 365")

	// Use the Valid() method to see if any of the checks failed. If they did, then re-render the template passing
	// in the form in the same way as before.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
