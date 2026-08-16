package models

import (
	"database/sql"
)

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{}
}

type Models struct {
}
