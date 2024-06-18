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
		Insert(book *Book, tx *sql.Tx) error
		Get(id int64) (*Book, error)
		Update(book *Book, tx *sql.Tx) error
		Delete(id int64, tx *sql.Tx) error
		GetAll() ([]*Book, error)
	}
	Authors interface {
		Insert(book *Author, tx *sql.Tx) error
		Get(id int64) (*Author, error)
		Update(book *Author, tx *sql.Tx) error
		Delete(id int64, tx *sql.Tx) error
		GetAll() ([]*Author, error)
	}
	Translations interface {
		Create() (*sql.Tx, error)
		Commit(*sql.Tx) error
		Rollback(*sql.Tx) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Books:        BookModel{DB: db},
		Authors:      AuthorModel{DB: db},
		Translations: Transactions{DB: db},
	}
}
