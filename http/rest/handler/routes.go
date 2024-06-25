package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func Register(r *mux.Router, lg *logrus.Logger, db *sqlx.DB) {
	handler := newHandler(lg, db)
	var api = r.PathPrefix("/api").Subrouter()
	v1 := api.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/services", handler.Create()).Methods(http.MethodPost)
	v1.HandleFunc("/services/{name}", handler.Get()).Methods(http.MethodGet)
	v1.HandleFunc("/services", handler.List()).Methods(http.MethodGet)
}
