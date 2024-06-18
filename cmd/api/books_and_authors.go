package main

import (
	"encoding/json"
	"fmt"
	"github.com/am-silex/go_library/internal/data"
	"net/http"
	"strconv"
)

func (app *application) updateBookAndAuthorHandler(w http.ResponseWriter, r *http.Request) {

	// Assume that book & author are items
	var inputData struct {
		Book   data.Book   `json:"book"`
		Author data.Author `json:"author"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputData)
	if err != nil {
		app.logger.Println(err)
		return
	}

	authorId, err := strconv.Atoi(r.PathValue("author_id"))
	if err != nil || authorId < 1 {
		return
	}

	bookId, err := strconv.Atoi(r.PathValue("book_id"))
	if err != nil || bookId < 1 {
		return
	}

	tx, err := app.models.Translations.Create()
	defer func() {
		if err == nil {
			err = app.models.Translations.Commit(tx)
			if err != nil {
				app.logger.Println(err)
			}
		}
		if err != nil {
			err = app.models.Translations.Rollback(tx)
			if err != nil {
				app.logger.Println(err)
			}
		}
	}()

	author := &data.Author{
		ID:          authorId,
		FirstName:   inputData.Author.FirstName,
		LastName:    inputData.Author.LastName,
		Bio:         inputData.Author.Bio,
		DateOfBirth: inputData.Author.DateOfBirth,
	}
	err = app.models.Authors.Update(author, tx)
	if err != nil {
		app.logger.Println(err)
		return
	}

	book := &data.Book{
		ID:       bookId,
		Title:    inputData.Book.Title,
		AuthorID: inputData.Book.AuthorID,
		Year:     inputData.Book.Year,
		ISBN:     inputData.Book.ISBN,
	}
	err = app.models.Books.Update(book, tx)
	if err != nil {
		app.logger.Println(err)
		return
	}

	// When sending an HTTP response, we want to include a Location header to let
	// the client know which URL they can find the newly-created resource at. We
	// make an empty http.Header map and then use the Set() method to add a new
	// Location header, interpolating the system-generated ID for our new book
	// in the URL.
	w.Header().Set("Location", fmt.Sprintf("/books/%d/authors/%d", bookId, authorId))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
