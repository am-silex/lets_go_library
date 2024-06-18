package main

import (
	"fmt"
	"net/http"
)

func (app *application) Serve() error {

	// log init
	// mux
	// listenAndServe

	mux := http.NewServeMux()
	mux.HandleFunc("POST /books", app.createBookHandler)
	mux.HandleFunc("GET /books", app.listBooksHandler)
	mux.HandleFunc("GET /books/{id}", app.getBookHandler)
	mux.HandleFunc("PUT /books/{id}", app.updateBookHandler)
	mux.HandleFunc("DELETE /books/{id}", app.deleteBookHandler)

	mux.HandleFunc("POST /authors", app.createAuthorHandler)
	mux.HandleFunc("GET /authors", app.listAuthorsHandler)
	mux.HandleFunc("GET /authors/{id}", app.getAuthorHandler)
	mux.HandleFunc("PUT /authors/{id}", app.updateAuthorHandler)
	mux.HandleFunc("DELETE /authors/{id}", app.deleteAuthorHandler)

	mux.HandleFunc("PUT /books/{book_id}/authors/{author_id}", app.updateBookAndAuthorHandler)

	httpServer := &http.Server{Addr: ":8080", Handler: app.authHandler(mux)}
	if err := httpServer.ListenAndServe(); err != nil {
		app.logger.Fatalln(fmt.Errorf("fatal error: %w", err))
	}
	return nil

}
