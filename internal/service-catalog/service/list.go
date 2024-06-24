package service

import (
	"context"
	"fmt"

	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
)

func (s *ServiceCatalog) List(ctx context.Context) ([]model.ServiceCatalog, error) {
	fmt.Printf("LIST Running....")
	svcCatalog, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return svcCatalog, nil
}
