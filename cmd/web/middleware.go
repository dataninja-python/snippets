package main

import (
	"fmt"
	"net/http"
)

const (
	cspName           string = "Content-Security-Policy"
	cspInstructions   string = "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	rpName            string = "Referer-Policy"
	rpInstructions    string = "origin-when-cross-origin"
	xctoName          string = "X-Content-Type-Options"
	xctoInstructions  string = "nosniff"
	xfoName           string = "X-Frame-Options"
	xfoInstructions   string = "deny"
	xxsspName         string = "X-XSS-Protection"
	xxsspInstructions string = "0"
)

type SecureHeader struct {
	CSPName           string
	CSPInstructions   string
	RPName            string
	RPInstructions    string
	XCTOName          string
	XCTOInstructions  string
	XFOName           string
	XFOInstructions   string
	XXSSPName         string
	XXSSPInstructions string
}

type RequestLog struct {
	IP     string
	Proto  string
	Method string
	URI    string
}

func NewSecureHeader() *SecureHeader {
	h := &SecureHeader{
		CSPName:           cspName,
		CSPInstructions:   cspInstructions,
		RPName:            rpName,
		RPInstructions:    rpInstructions,
		XCTOName:          xctoName,
		XCTOInstructions:  xctoInstructions,
		XFOName:           xfoName,
		XFOInstructions:   xfoInstructions,
		XXSSPName:         xxsspName,
		XXSSPInstructions: xxsspInstructions,
	}

	return h
}

func secureHeaders(next http.Handler) http.Handler {
	// Initialize a secure header object
	// sHeader collects all the standard header information defined above into one object
	headerInfo := NewSecureHeader()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerInfo.CSPName, headerInfo.CSPInstructions)
		w.Header().Set(headerInfo.RPName, headerInfo.RPInstructions)
		w.Header().Set(headerInfo.XCTOName, headerInfo.XCTOInstructions)
		w.Header().Set(headerInfo.XFOName, headerInfo.XFOInstructions)
		w.Header().Set(headerInfo.XXSSPName, headerInfo.XXSSPInstructions)

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var l RequestLog
		l.IP = r.RemoteAddr
		l.Proto = r.Proto
		l.Method = r.Method
		l.URI = r.URL.RequestURI()

		app.logger.Info("received request", "ip", l.IP, "proto", l.Proto, "method", l.Method,
			"uri", l.URI)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function that always runs in the event of a panic
		// as Go unwinds the stack
		defer func() {
			// Use the builtin recover function to check if there has been a panic. if there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500 Internal Server response.
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
