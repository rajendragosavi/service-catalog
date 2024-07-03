package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/rajendragosavi/service-catalog/docs"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel"
)

// @title Service API
// @version 1.0
// @description This is a sample service API.
// @host localhost:80
// @BasePath /api/v1
func Register(r *mux.Router, lg *logrus.Entry, db *sqlx.DB) {
	shutdown := initTracer()
	defer shutdown()

	tracer = otel.Tracer("service-catalog")

	handler := newHandler(lg, db)
	var api = r.PathPrefix("/api").Subrouter()
	api.Use(loggingMiddleware(lg))
	api.Use(tracingMiddleware)
	v1 := api.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/services/{name}", handler.Get()).Methods(http.MethodGet)
	v1.HandleFunc("/services", handler.Create()).Methods(http.MethodPost)
	v1.HandleFunc("/services", handler.List()).Methods(http.MethodGet)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	// Expose the registered metrics via HTTP
	http.Handle("/metrics", promhttp.Handler())
}
