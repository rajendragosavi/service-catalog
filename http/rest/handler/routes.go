package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func Register(r *mux.Router, lg *logrus.Logger, db *sqlx.DB) {
	handler := newHandler(lg, db)
	r.HandleFunc("/services", handler.Create()).Methods(http.MethodPost)
	//r.HandleFunc("/services/{name}", handler.Get()).Methods(http.MethodGet)
	//r.HandleFunc("/services", handler.Get()).Methods(http.MethodGet)
	// r.HandleFunc("/services", handler.Create()).Methods(http.MethodPost)
	// r.HandleFunc("/services", handler.Create()).Methods(http.MethodPost)

}
