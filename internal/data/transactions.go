package data

import "database/sql"

type Transactions struct {
	DB *sql.DB
}

func (service Transactions) Create() (tx *sql.Tx, err error) {

	tx, err = service.DB.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (service Transactions) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (service Transactions) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}
