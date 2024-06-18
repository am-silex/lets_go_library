package data

import (
	"database/sql"
	"errors"
)

// ErrRecordNotFound Define a custom ErrRecordNotFound error.
var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Books interface {
		Insert(book *Book) error
		Get(id int64) (*Book, error)
		Update(book *Book) error
		Delete(id int64) error
		GetAll() ([]*Book, error)
	}
	Authors interface {
		Insert(book *Author) error
		Get(id int64) (*Author, error)
		Update(book *Author) error
		Delete(id int64) error
		GetAll() ([]*Author, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Books:   BookModel{DB: db},
		Authors: AuthorModel{DB: db},
	}
}
