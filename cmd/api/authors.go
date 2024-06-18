package main

import (
	"encoding/json"
	"fmt"
	"github.com/am-silex/go_library/internal/data"
	"net/http"
	"strconv"
)

func (app *application) createAuthorHandler(w http.ResponseWriter, r *http.Request) {

	var inputData data.Author
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputData)
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	author := &data.Author{
		FirstName:   inputData.FirstName,
		LastName:    inputData.LastName,
		Bio:         inputData.Bio,
		DateOfBirth: inputData.DateOfBirth,
	}

	err = app.models.Authors.Insert(author, nil)
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// When sending an HTTP response, we want to include a Location header to let
	// the client know which URL they can find the newly-created resource at. We
	// make an empty http.Header map and then use the Set() method to add a new
	// Location header, interpolating the system-generated ID for our new author
	// in the URL.
	w.Header().Set("Location", fmt.Sprintf("/authors/%d", author.ID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(author)

}

func (app *application) updateAuthorHandler(w http.ResponseWriter, r *http.Request) {

	var inputData data.Author
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputData)
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	author := &data.Author{
		ID:          inputData.ID,
		LastName:    inputData.LastName,
		FirstName:   inputData.FirstName,
		Bio:         inputData.Bio,
		DateOfBirth: inputData.DateOfBirth,
	}

	err = app.models.Authors.Update(author, nil)
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// When sending an HTTP response, we want to include a Location header to let
	// the client know which URL they can find the newly-created resource at. We
	// make an empty http.Header map and then use the Set() method to add a new
	// Location header, interpolating the system-generated ID for our new author
	// in the URL.
	w.Header().Set("Location", fmt.Sprintf("/authors/%d", author.ID))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(author)
}

func (app *application) deleteAuthorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return
	}
	err = app.models.Authors.Delete(id, nil)
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("author successfully deleted"))
}

func (app *application) getAuthorHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return
	}
	author, err := app.models.Authors.Get(id)
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(author)

}

func (app *application) listAuthorsHandler(w http.ResponseWriter, r *http.Request) {

	authors, err := app.models.Authors.GetAll()
	if err != nil {
		app.logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authors)
}
