package service

import "github.com/rajendragosavi/service-catalog/internal/service-catalog/repository"

type Service struct {
	repo repository.Repository
}

func NewService(r repository.Repository) Service {
	return Service{
		repo: r,
	}
}
