package service

import "github.com/rajendragosavi/service-catalog/internal/service-catalog/repository"

type ServiceCatalog struct {
	repo repository.Repository
}

func NewService(r repository.Repository) ServiceCatalog {
	return ServiceCatalog{
		repo: r,
	}
}
