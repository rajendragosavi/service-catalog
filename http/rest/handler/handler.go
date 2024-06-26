package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/repository"
	catalog "github.com/rajendragosavi/service-catalog/internal/service-catalog/service"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger         *logrus.Entry
	serviceCatalog *catalog.ServiceCatalog
}

func newHandler(lg *logrus.Entry, db *sqlx.DB) Service {
	catalogRepo := repository.NewRepository(db)
	serviceCatalog := catalog.NewServiceCatalog(&catalogRepo) // service catalog has access to repo. Repo has methods to talk to DB model to do operations
	return Service{
		logger:         lg,
		serviceCatalog: serviceCatalog,
	}
}
