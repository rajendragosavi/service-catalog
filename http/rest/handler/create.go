package handler

import (
	"net/http"

	catalog "github.com/rajendragosavi/service-catalog/internal/service-catalog/service"
)

type request struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Versions    []string `json:"versions"`
	//Versions interface{}  `json:"versions"`
	Status int `json:"status"`
	// Status model.Status `json:"status"`
}

type Response struct {
	ID string `json:"id"`
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
	s.logger.Debugln("running create service http handler")
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		// If there is an error, respond to the client with the error message and a 400 status code.
		err := s.decode(r, &req)
		if err != nil {
			s.logger.Errorf("Invalid input data. error : %+v \n", err)
			s.respond(w, err, 0)
			return
		}
		id, err := s.serviceCatalog.Create(r.Context(), catalog.CreateParams{
			Name:        req.Name,
			Description: req.Description,
			Versions:    req.Versions,
			Status:      req.Status,
		})
		if err != nil {
			s.logger.Errorf("could not create service. error : %+v \n", err)
			s.respond(w, err, 0)
			return
		}
		s.respond(w, Response{ID: id}, http.StatusOK)
	}
}
