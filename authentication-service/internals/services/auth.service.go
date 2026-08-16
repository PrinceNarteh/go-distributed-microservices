package services

import (
	"authentication/internals/helpers"
	"authentication/internals/models"
	"authentication/internals/repositories"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

var _ AuthService = (*authService)(nil)

type AuthService interface {
	Login(data *models.LoginRequest) (*models.AuthResponse, error)
	Register(data *models.RegisterRequest) (*models.AuthResponse, error)
}

type authService struct {
	repo *repositories.Repository
}

func (s *authService) Login(data *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.repo.User.GetByEmail(data.Email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, fmt.Errorf("invalid credentials")
		default:
			return nil, fmt.Errorf("something went wrong")
		}
	}

	if isMatch := user.ComparePassword(data.Password); !isMatch {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := helpers.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *authService) Register(data *models.RegisterRequest) (*models.AuthResponse, error) {
	now := time.Now()
	user := &models.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  data.Password,
		IsActive:  true,
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// hash password
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	// create user in database
	if err := s.repo.User.Insert(user); err != nil {
		return nil, err
	}

	token, err := helpers.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}
