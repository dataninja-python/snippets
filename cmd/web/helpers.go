package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// This helper function sends a specific status code and description to the user.
// Deals with 400 Bad Request errors right now
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// This helper deals with 404 Not Found errors right now.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Retrieve the appropriate template set from teh cache based on the page name
	// (like 'home.tmpl.html'). If no entry exists in the cache with the provided name, then create
	// a new error and call the serverError() helper method that we made earlier and return.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)
	// Write the template to the buffer, instead of straight to the http.ResponseWriter. If there's an error, call our
	// serverError() helper and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// If the template is written to the buffer without an error we can proceed.
	// Write out the provided HTTP status code ('200 OK', '400 Bad Request' etc).
	w.WriteHeader(status)

	// Use the result from the buffer to the http.ResponseWriter
	buf.WriteTo(w)
	/*
		// Execute the template set and write the response body. Again, if there is any error we call the serverError()
		// helper.
		err := ts.ExecuteTemplate(w, "base", data)
		if err != nil {
			app.serverError(w, r, err)
		}
	*/
}

// Create a newTemplateData() helper, which returns a pointer to a templateData
func (app *application) newTemplateData(r *http.Request) templateData {
	var year int
	year = time.Now().Year()
	// fmt.Println("Year:", year)
	return templateData{
		CurrentYear: year,
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call Parse/form() on requests like the createSnippetpost handler.
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}

	return nil
}
