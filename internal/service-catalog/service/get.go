package service

import (
	"context"
	"fmt"

	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
)

func (s *ServiceCatalog) Get(ctx context.Context, name string) (*model.ServiceCatalog, error) {
	fmt.Printf("GET Running....")
	svcCatalog, err := s.repo.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return svcCatalog, nil
}
