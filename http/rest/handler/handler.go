package handler

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	catalogRepo "github.com/rajendragosavi/service-catalog/internal/service-catalog/repository"
	catalog "github.com/rajendragosavi/service-catalog/internal/service-catalog/service"
	"github.com/sirupsen/logrus"
)

type service struct {
	logger         *logrus.Logger
	router         *mux.Router
	serviceCatalog catalog.ServiceCatalog
}

func newHandler(lg *logrus.Logger, db *sqlx.DB) service {
	return service{
		logger:         lg,
		serviceCatalog: catalog.NewService(catalogRepo.NewRepository(db)),
	}
}
