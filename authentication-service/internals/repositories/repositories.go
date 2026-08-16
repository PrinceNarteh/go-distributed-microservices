package repositories

import (
	"database/sql"
	"time"
)

const dbTimeout = 5 * time.Second

type Repository struct {
	User UserRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User: &userRepository{DB: db},
	}
}
