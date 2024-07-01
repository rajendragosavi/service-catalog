package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	catalog "github.com/rajendragosavi/service-catalog/internal/service-catalog/service"
	customErr "github.com/rajendragosavi/service-catalog/pkg/errors"
	"github.com/sirupsen/logrus"
)

type request struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Versions    []string `json:"versions"`
}

type Response struct {
	ID string `json:"serice_id"`
}

// Create godoc
// @Summary Create a new service
// @Description Create a new service with the provided details
// @Tags services
// @Accept  json
// @Produce  json
// @Param   service  body  Service  true  "Service to create"
// @Success 201 {object} Service
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/services [post]
func (s Service) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		var ok bool
		s.logger, ok = r.Context().Value("logger").(*logrus.Entry)
		if !ok {
			http.Error(w, "Logger not found in context", http.StatusInternalServerError)
			return
		}
		s.logger.Debugln("CREATE service http handler")
		// If there is an error, respond to the client with the error message and a 400 status code.
		err := s.decode(r, &req)
		if err != nil {
			s.logger.Errorf("Invalid input data. error : %+v \n", err)
			s.respond(w, err, 0)
			return
		}
		// Validate the request
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			s.logger.Errorf("invalid request data, name field is mandatory : %+v \n", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ErrorResponse{ErrorMessage: "Invalid input data. name field is mandatory"})
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		id, err := s.serviceCatalog.Create(r.Context(), catalog.CreateParams{
			Name:        req.Name,
			Description: req.Description,
			Versions:    req.Versions,
		})
		if err != nil {
			s.logger.Errorf("could not create service. error : %+v \n", err)
			if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
				s.respond(w, customErr.DuplicateKeyError{
					Wrapped: errors.New("service name already exist"),
				}, 0)
			} else {
				s.respond(w, err, 0)
			}
			return
		}
		s.respond(w, Response{ID: id}, http.StatusOK)
	}
}
