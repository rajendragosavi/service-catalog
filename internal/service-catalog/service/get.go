package service

import (
	"context"

	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
)

func (s *Service) Get(ctx context.Context, id int) (*model.ServiceCatalog, error) {
	svcCatalog, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	return svcCatalog, nil
}
