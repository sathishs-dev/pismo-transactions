package repository

import "github.com/jmoiron/sqlx"

type (
	pismoRepo struct {
		db *sqlx.DB
	}

	PismoRepo interface{}
)

func NewPismoRepo(db *sqlx.DB) PismoRepo {
	return &pismoRepo{
		db,
	}
}
