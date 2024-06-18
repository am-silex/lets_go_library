package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	AuthorID int    `json:"author_id"`
	Year     int    `json:"year,omitempty"`
	ISBN     string `json:"isbn"`
}

// BookModel Define a struct type which wraps a sql.DB connection pool.
type BookModel struct {
	DB *sql.DB
}

// Insert The method accepts a pointer to a book struct, which should contain
// the data for the new record.
func (m BookModel) Insert(book *Book, tx *sql.Tx) error {
	query := `
		INSERT INTO public.books (title, authorid, year, isbn)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	args := []interface{}{book.Title, book.AuthorID, book.Year, book.ISBN}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	switch tx {
	case nil:
		return m.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID)
	default:
		return tx.QueryRowContext(ctx, query, args...).Scan(&book.ID)
	}

}

// Get fetches a specific record from the books table.
func (m BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, title, authorid, year, isbn
		FROM public.books
		WHERE id = $1`

	var book Book

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.AuthorID,
		&book.Year,
		&book.ISBN,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &book, nil
}

// Update updates a specific record in the books table.
func (m BookModel) Update(book *Book, tx *sql.Tx) error {
	query := `
        UPDATE public.books
        SET title = $1, year = $2, authorid = $3, isbn = $4
        WHERE id = $5
        RETURNING id`

	args := []interface{}{
		book.Title,
		book.Year,
		book.AuthorID,
		book.ISBN,
		book.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	switch tx {
	case nil:
		err := m.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID)
		if err != nil {
			return err
		}
	default:
		err := tx.QueryRowContext(ctx, query, args...).Scan(&book.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes a specific record from the books table.
func (m BookModel) Delete(id int64, tx *sql.Tx) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM public.books
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	switch tx {
	case nil:
		result, err := m.DB.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return ErrRecordNotFound
		}
	default:
		result, err := tx.ExecContext(ctx, query, id)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return ErrRecordNotFound
		}
	}

	return nil
}

// GetAll method returns a slice of books.
func (m BookModel) GetAll() ([]*Book, error) {
	query := `
		SELECT id, title, authorid, year, isbn
		FROM public.books
		ORDER BY title ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := []*Book{}

	for rows.Next() {
		var book Book

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.AuthorID,
			&book.Year,
			&book.ISBN,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
