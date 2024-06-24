package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/repository"
	catalog "github.com/rajendragosavi/service-catalog/internal/service-catalog/service"
	"github.com/sirupsen/logrus"
)

type service struct {
	logger         *logrus.Logger
	serviceCatalog *catalog.ServiceCatalog
}

func newHandler(lg *logrus.Logger, db *sqlx.DB) service {
	catalogRepo := repository.NewRepository(db)
	serviceCatalog := catalog.NewServiceCatalog(&catalogRepo) // service catalog has access to repo. Repo has methods to talk to DB model to do operations
	return service{
		logger:         lg,
		serviceCatalog: serviceCatalog,
	}
}
