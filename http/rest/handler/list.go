package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

type response struct {
	ID              string     `json:"service_id"`
	Name            string     `json:"service_name"`
	Description     string     `json:"description"`
	CreatedTime     time.Time  `json:"created_time"`
	LastUpdatedTime *time.Time `json:"last_updated_time,omitempty"`
}
type listResponse struct {
	listResponse []response
}

func (s service) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO we will add more filters here
		if len(mux.Vars(r)) == 0 {
			uuID := r.Header.Get("userID")
			if uuID == "" {
				err := errors.New("userID header is missing")
				s.logger.Errorf("error listing services - %+v \n", err)
				s.respond(w, err, 0)
				return
			}

			exist, err := s.serviceCatalog.CheckUserExists(r.Context(), uuID)
			if err != nil {
				s.logger.Errorf("error listing services - %+v \n", err)
				s.respond(w, err, 0)
				return
			}
			if !exist {
				s.logger.Warn("user dont have access to the service")
				s.respond(w, err, 0)
				return
			}
			uID, err := uuid.Parse(uuID)
			if err != nil {
				fmt.Printf("Invalid UUID string: %v\n", err)
				return
			}
			listResponse, err := s.serviceCatalog.ListServicesForUser(r.Context(), uID)
			if err != nil {
				s.logger.Errorf("error listing services - %+v \n", err)
				s.respond(w, err, 0)
				return
			}
			s.respond(w, listResponse, http.StatusOK)
		}
	}
}
