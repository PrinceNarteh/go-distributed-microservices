package services

import "authentication/internals/repositories"

type Services struct {
	Auth AuthService
}

func NewService(repo *repositories.Repository) *Services {
	return &Services{
		Auth: &authService{repo: repo},
	}
}
