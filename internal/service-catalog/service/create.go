package service

import (
	"context"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
)

type CreateParams struct {
	Name        string   `valid:"required"`
	Description string   `valid:"required"`
	Versions    []string `valid:"required"`
	//Versions interface{}  `valid:"required"`
	Status model.Status `valid:"required"`
}

func (s *ServiceCatalog) Create(ctx context.Context, params CreateParams) (string, error) {
	fmt.Println("create sql opeartion started")
	if _, err := govalidator.ValidateStruct(params); err != nil {
		fmt.Printf("err from validator - %+v \n", err)
		return "", err // TODO error handling
	}
	fmt.Println("create sql txn started")
	tx, err := s.repo.Db.BeginTxx(ctx, nil)
	if err != nil {
		fmt.Printf("error in begin txn - %+v , error - %+v \n", tx, err)
		return "", err // TODO error handling
	}
	// Defer rollback in case of failure/error
	defer tx.Rollback()

	//x := pq.Array(params)
	//	fmt.Printf("x.Value = %+v \n", x)
	//	fmt.Printf("create Params - %+v \n", params)
	obj := model.ServiceCatalog{
		Name:        params.Name,
		Description: params.Description,
		Versions:    params.Versions,
		CreatedOn:   time.Now().UTC(),
	}
	fmt.Printf("sql obj for create service entry in the db - %+v \n", obj.Versions)
	err = s.repo.Create(ctx, &obj)
	if err != nil {
		fmt.Printf("ERRROR  - %+v \n", err)
		return "", err // TODO error handling
	}
	err = tx.Commit()
	return obj.ID, err
}
