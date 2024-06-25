package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type response struct {
	ID          string `json:"service_id"`
	Name        string `json:"service_name"`
	Description string `json:"description"`
	// Status          model.Status `json:"status"`
	Stats           int        `json:"status"`
	CreatedTime     time.Time  `json:"created_time"`
	LastUpdatedTime *time.Time `json:"last_updated_time,omitempty"`
}
type listResponse struct {
	listResponse []response
}

// List godoc
// @Summary List all services
// @Description List all services available
// @Tags services
// @Produce  json
// @Success 200 {array} Service
// @Failure 500 {object} ErrorResponse
// @Router /services [get]
func (s service) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO we will add more filters here
		if len(mux.Vars(r)) == 0 {
			listResponse, err := s.serviceCatalog.List(r.Context())
			if err != nil {
				s.logger.Errorf("error listing services - %+v \n", err)
				s.respond(w, err, 0)
				return
			}
			s.respond(w, listResponse, http.StatusOK)
		}
	}
}
