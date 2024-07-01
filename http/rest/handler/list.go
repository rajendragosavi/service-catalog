package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

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

// List godoc
// @Summary List all services
// @Description List all services available
// @Tags services
// @Produce  json
// @Success 200 {array} Service
// @Failure 500 {object} ErrorResponse
// @Router /v1/services [get]
func (s Service) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO we will add more filters here
		var ok bool
		s.logger, ok = r.Context().Value("logger").(*logrus.Entry)
		if !ok {
			http.Error(w, "Logger not found in context", http.StatusInternalServerError)
			return
		}
		s.logger.Debugln("LIST service http handler")
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
				s.logger.Errorf("Invalid UUID string: %v\n", err)
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
