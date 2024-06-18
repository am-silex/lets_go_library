package main

import (
	"net/http"
)

func (app *application) authHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Println(r.Method, r.RequestURI)
		// Auth checks goes here... Bypassing for now.
		h.ServeHTTP(w, r)
	})
}
