package models

import (
	"database/sql"
	"time"
)

const dbTimeout = 3 * time.Second

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{}
}

type Models struct {
}
