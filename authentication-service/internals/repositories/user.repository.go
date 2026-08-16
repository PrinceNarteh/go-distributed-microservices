package repositories

import (
	"authentication/internals/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var _ UserRepository = (*userRepository)(nil)

type UserRepository interface {
	Insert(user *models.User) error
	GetAll() ([]*models.User, error)
	GetOne(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(userID int, updates map[string]any) error
	Delete(userID int) error
	ResetPassword(userID int, newPassword string) error
}

type userRepository struct {
	DB *sql.DB
}

func (r *userRepository) Insert(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	var userID int
	stmt := `INSERT INTO users (first_name, last_name, email, password, is_active, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	if err := r.DB.QueryRowContext(ctx, stmt, user.FirstName, user.LastName, user.Email, user.Password, user.IsActive, user.CreatedAt, user.UpdatedAt).Scan(userID); err != nil {
		return err
	}
	user.ID = userID

	return nil
}

func (r *userRepository) GetAll() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, first_name, last_name, email, password, is_active, created_at, updated_at FROM users ORDER BY last_name`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := new(models.User)
		if err := rows.Scan(user.ID, user.Email, user.FirstName, user.LastName, user.IsActive, user.CreatedAt, user.UpdatedAt); err != nil {
			log.Println("Error scanning user: ", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return users, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, first_name, last_name, email, password, is_active, created_at, updated_at FROM users WHERE email = $1`
	row := r.DB.QueryRowContext(ctx, query, email)

	user := new(models.User)
	if err := row.Scan(user.ID, user.Email, user.FirstName, user.LastName, user.IsActive, user.CreatedAt, user.UpdatedAt); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetOne(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, first_name, last_name, email, password, is_active, created_at, updated_at FROM users WHERE id = $1`
	row := r.DB.QueryRowContext(ctx, query, id)

	user := new(models.User)
	if err := row.Scan(user.ID, user.Email, user.FirstName, user.LastName, user.IsActive, user.CreatedAt, user.UpdatedAt); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(userID int, updates map[string]any) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	if len(updates) == 0 {
		return nil
	}

	var columns []string
	var args []any
	argId := 1

	for key, value := range updates {
		columns = append(columns, fmt.Sprintf("%s = $%d", key, argId))
		args = append(args, value)
		argId++
	}

	args = append(args, userID)
	query := fmt.Sprintf(`UPDATE users SET %s WHERE id = %d`, strings.Join(columns, ", "), argId)

	_, err := r.DB.ExecContext(ctx, query, args...)
	return err
}

func (r *userRepository) Delete(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `DELETE FROM users WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, stmt, userID)

	return err
}

func (r *userRepository) ResetPassword(userID int, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `UPDATE users SET password = $1 WHERE id = $2`
	_, err = r.DB.ExecContext(ctx, stmt, string(hash), userID)

	return err
}
