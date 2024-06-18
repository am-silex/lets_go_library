package main

import (
	"encoding/json"
	"fmt"
	"github.com/am-silex/go_library/internal/data"
	"net/http"
	"strconv"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {

	var inputData data.Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputData)
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book := &data.Book{
		Title:    inputData.Title,
		AuthorID: inputData.AuthorID,
		Year:     inputData.Year,
		ISBN:     inputData.ISBN,
	}

	err = app.models.Books.Insert(book, nil)
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("book wasn't created"))
		return
	}
	// When sending an HTTP response, we want to include a Location header to let
	// the client know which URL they can find the newly-created resource at. We
	// make an empty http.Header map and then use the Set() method to add a new
	// Location header, interpolating the system-generated ID for our new book
	// in the URL.
	w.Header().Set("Location", fmt.Sprintf("/books/%d", book.ID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)

}

func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {

	var inputData data.Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputData)
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book := &data.Book{
		ID:       inputData.ID,
		Title:    inputData.Title,
		AuthorID: inputData.AuthorID,
		Year:     inputData.Year,
		ISBN:     inputData.ISBN,
	}

	err = app.models.Books.Update(book, nil)
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("book wasn't updated"))
		return
	}
	// When sending an HTTP response, we want to include a Location header to let
	// the client know which URL they can find the newly-created resource at. We
	// make an empty http.Header map and then use the Set() method to add a new
	// Location header, interpolating the system-generated ID for our new book
	// in the URL.
	w.Header().Set("Location", fmt.Sprintf("/books/%d", book.ID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = app.models.Books.Delete(id, nil)
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("book wasn't deleted"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("book successfully deleted"))
}

func (app *application) getBookHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	book, err := app.models.Books.Get(id)
	if err != nil {
		app.logger.Println(err)
		w.Write([]byte("book not found"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)

}

func (app *application) listBooksHandler(w http.ResponseWriter, r *http.Request) {

	books, err := app.models.Books.GetAll()
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}
