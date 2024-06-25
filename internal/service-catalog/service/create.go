package service

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
)

type CreateParams struct {
	Name        string   `valid:"required"`
	Description string   `valid:"required"`
	Versions    []string `valid:"required"`
	// Status      model.Status `valid:"required"`
	Status int `valid:"required"`
}

func (s *ServiceCatalog) Create(ctx context.Context, params CreateParams) (string, error) {
	if _, err := govalidator.ValidateStruct(params); err != nil {
		return "", err // TODO error handling
	}
	tx, err := s.repo.BeginTxx(ctx, nil)
	if err != nil {
		return "", err // TODO error handling
	}
	// Defer rollback in case of failure/error
	defer tx.Rollback()

	obj := model.ServiceCatalog{
		Name:        params.Name,
		Description: params.Description,
		Status:      1,
		Versions:    params.Versions,
		CreatedOn:   time.Now().UTC(),
	}
	id, err := s.repo.Create(ctx, &obj)
	if err != nil {
		return "", err // TODO error handling
	}
	obj.ID = id
	err = tx.Commit()
	return obj.ID, err
}
