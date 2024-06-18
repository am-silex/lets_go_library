package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Author struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Bio         string `json:"bio,omitempty"`
	DateOfBirth int    `json:"date_of_birth"`
}

// AuthorModel Define a struct type which wraps a sql.DB connection pool.
type AuthorModel struct {
	DB *sql.DB
}

// Insert The method accepts a pointer to a author struct, which should contain
// the data for the new record.
func (m AuthorModel) Insert(author *Author, tx *sql.Tx) error {
	query := `
		INSERT INTO public.authors (first_name, last_name, bio, date_of_birth)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	args := []interface{}{author.FirstName, author.LastName, author.Bio, author.DateOfBirth}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	switch tx {
	case nil:
		return m.DB.QueryRowContext(ctx, query, args...).Scan(&author.ID)
	default:
		return tx.QueryRowContext(ctx, query, args...).Scan(&author.ID)
	}

}

// Get fetches a specific record from the authors table.
func (m AuthorModel) Get(id int64) (*Author, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, first_name, last_name, bio, date_of_birth
		FROM public.authors
		WHERE id = $1`

	var author Author

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&author.ID,
		&author.FirstName,
		&author.LastName,
		&author.Bio,
		&author.DateOfBirth,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &author, nil
}

// Update updates a specific record in the authors table.
func (m AuthorModel) Update(author *Author, tx *sql.Tx) error {
	query := `
        UPDATE public.authors
        SET first_name = $1, last_name = $2, bio = $3, date_of_birth = $4
        WHERE id = $5
        RETURNING id`

	args := []interface{}{
		author.FirstName,
		author.LastName,
		author.Bio,
		author.DateOfBirth,
		author.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	switch tx {
	case nil:
		err := m.DB.QueryRowContext(ctx, query, args...).Scan(&author.ID)
		if err != nil {
			return err
		}
	default:
		err := tx.QueryRowContext(ctx, query, args...).Scan(&author.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes a specific record from the authors table.
func (m AuthorModel) Delete(id int64, tx *sql.Tx) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM public.authors
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

// GetAll method returns a slice of authors.
func (m AuthorModel) GetAll() ([]*Author, error) {
	query := `
		SELECT id, first_name, last_name, bio, date_of_birth
		FROM public.authors
		ORDER BY last_name ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	authors := []*Author{}

	for rows.Next() {
		var author Author

		err := rows.Scan(
			&author.ID,
			&author.FirstName,
			&author.LastName,
			&author.Bio,
			&author.DateOfBirth,
		)
		if err != nil {
			return nil, err
		}

		authors = append(authors, &author)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}
