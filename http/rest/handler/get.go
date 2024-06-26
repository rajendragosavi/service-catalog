package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	err "github.com/rajendragosavi/service-catalog/pkg/errors"
)

// Get godoc
// @Summary Get service details
// @Description Get details of the service specified by name
// @Tags services
// @Produce  json
// @Param   name  path  string  true  "Service Name"
// @Success 200 {object} Service
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/services/{name} [get]
func (s Service) Get() http.HandlerFunc {
	type response struct {
		ID              string     `json:"service_id"`
		Name            string     `json:"service_name"`
		Description     string     `json:"description"`
		CreatedTime     time.Time  `json:"created_time"`
		LastUpdatedTime *time.Time `json:"last_updated_time,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		fmt.Printf("name value - %+v \n", name)
		if name == "" {
			s.respond(w, err.ErrorArgument{
				Wrapped: errors.New("valid name must provide in path"),
			}, 0)
			return
		}
		getResponse, err := s.serviceCatalog.Get(r.Context(), name)
		if err != nil {
			s.logger.Errorf("could not get service. error : %+v \n", err)
			s.respond(w, err, 0)
			return
		}
		s.respond(w, response{
			ID:              getResponse.ID,
			Name:            getResponse.Name,
			Description:     getResponse.Description,
			CreatedTime:     getResponse.CreatedOn,
			LastUpdatedTime: getResponse.UpdatedOn,
		}, http.StatusOK)
	}
}
