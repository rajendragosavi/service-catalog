package service

import (
	"context"

	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/repository"
)

type ServiceCatalog struct {
	repo repository.Repository
}

func NewServiceCatalog(r *repository.Repository) *ServiceCatalog {
	return &ServiceCatalog{
		repo: *r,
	}
}

type ServiceCatalogService interface {
	Create(ctx context.Context, params CreateParams) (string, error)
	List(ctx context.Context) ([]*model.ServiceCatalog, error)
	Get(ctx context.Context, name string) (*model.ServiceCatalog, error)
}
