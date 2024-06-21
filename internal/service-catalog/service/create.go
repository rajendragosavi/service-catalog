package service

import (
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
)

type CreateParams struct {
	Name        string   `valid:"required"`
	Description string   `valid:"required"`
	Versions    []string `valid:"required"`
}

func (s *Service) Create(ctx context.Context, params CreateParams) (int, error) {
	if _, err := govalidator.ValidateStruct(params); err != nil {
		return 0, err // TODO error handling
	}
	tx, err := s.repo.Db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err // TODO error handling
	}
	// Defer rollback in case of failure/error
	defer tx.Rollback()
	obj := model.ServiceCatalog{
		Name:        params.Name,
		Description: params.Description,
		Versions:    params.Versions,
	}
	err = s.repo.Create(ctx, &obj)
	if err != nil {
		return 0, err // TODO error handling
	}
	err = tx.Commit()
	return obj.ID, nil
}
