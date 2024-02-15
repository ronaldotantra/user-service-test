// This file contains the repository implementation layer.
package repository

import (
	"database/sql"
)

type repository struct {
	Db *sql.DB
}

type NewRepositoryOptions struct {
	Db *sql.DB
}

func NewRepository(opts NewRepositoryOptions) RepositoryInterface {
	return &repository{
		Db: opts.Db,
	}
}
