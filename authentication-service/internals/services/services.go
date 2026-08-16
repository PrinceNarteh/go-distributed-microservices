package services

import "authentication/internals/repositories"

type Services struct{}

func NewService(repo *repositories.Repository) *Services {
	return &Services{}
}
