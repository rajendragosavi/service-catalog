package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/rajendragosavi/service-catalog/docs"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Service API
// @version 1.0
// @description This is a sample service API.
// @host localhost:80
// @BasePath /api/v1

func Register(r *mux.Router, lg *logrus.Entry, db *sqlx.DB) {
	handler := newHandler(lg, db)
	var api = r.PathPrefix("/api").Subrouter()
	api.Use(loggingMiddleware(lg))
	v1 := api.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/services/{name}", handler.Get()).Methods(http.MethodGet)
	v1.HandleFunc("/services", handler.Create()).Methods(http.MethodPost)
	v1.HandleFunc("/services", handler.List()).Methods(http.MethodGet)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
