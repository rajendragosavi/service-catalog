package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
)

type response struct {
	ID              string       `json:"service_id"`
	Name            string       `json:"service_name"`
	Description     string       `json:"description"`
	Status          model.Status `json:"status"`
	CreatedTime     time.Time    `json:"created_time"`
	LastUpdatedTime *time.Time   `json:"last_updated_time,omitempty"`
}
type listResponse struct {
	listResponse []response
}

func (s service) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(mux.Vars(r)) == 0 {
			listResponse, err := s.serviceCatalog.List(r.Context())
			if err != nil {
				fmt.Printf("error in list -  %+v , response - %+v \n", err, listResponse)
				s.respond(w, err, 0)
				return
			}
			s.respond(w, listResponse, http.StatusOK)
		}
	}
}
