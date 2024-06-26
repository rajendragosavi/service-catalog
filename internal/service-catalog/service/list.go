package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
)

func (s *ServiceCatalog) ListServicesForUser(ctx context.Context, userID uuid.UUID) ([]model.ServiceCatalog, error) {
	svcCatalog, err := s.repo.ListServicesForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return svcCatalog, nil
}

func (s *ServiceCatalog) ListUsersWithAccess(ctx context.Context, serviceID uuid.UUID) ([]model.User, error) {
	users, err := s.repo.ListUsersWithAccess(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *ServiceCatalog) CheckUserExists(ctx context.Context, userID string) (bool, error) {

	exists, err := s.repo.CheckUserExists(ctx, userID)
	if err != nil {
		return false, err
	}
	return exists, nil
}
