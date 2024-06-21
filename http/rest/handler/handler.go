package handler

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type service struct {
	logger *logrus.Logger
	router *mux.Router
}

func newHandler(lg *logrus.Logger, db *sqlx.DB) service {
	return service{
		logger: lg,
	}
}
