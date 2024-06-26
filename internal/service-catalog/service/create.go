package service

import (
	"context"
	"fmt"
	"time"

	// "github.com/asaskevich/govalidator"
	//	"github.com/go-playground/validator/v10"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
)

type CreateParams struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"required"`
	Versions    []string `valid:"required"`
}

func (s *ServiceCatalog) Create(ctx context.Context, params CreateParams) (string, error) {

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
		fmt.Printf("error from repo - %+v \n", err)
		return "", err // TODO error handling
	}
	obj.ID = id
	err = tx.Commit()
	return obj.ID, err
}
