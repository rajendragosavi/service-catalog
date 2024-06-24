package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rajendragosavi/service-catalog/internal/service-catalog/model"
	err "github.com/rajendragosavi/service-catalog/pkg/errors"
)

func (s service) Get() http.HandlerFunc {
	fmt.Println("GET handler running")
	type response struct {
		ID              string       `json:"service_id"`
		Name            string       `json:"service_name"`
		Description     string       `json:"description"`
		Status          model.Status `json:"status"`
		CreatedTime     time.Time    `json:"created_time"`
		LastUpdatedTime *time.Time   `json:"last_updated_time,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		name := vars["name"]
		if name == "" {

			s.respond(w, err.ErrorArgument{
				Wrapped: errors.New("valid name must provide in path"),
			}, 0)
			return
		}

		getResponse, err := s.serviceCatalog.Get(r.Context(), name)
		if err != nil {
			fmt.Printf("error in get -  %+v , response - %+v \n", err, getResponse)
			s.respond(w, err, 0)
			return
		}
		s.respond(w, response{
			ID:              getResponse.ID,
			Name:            getResponse.Name,
			Description:     getResponse.Description,
			Status:          getResponse.Status,
			CreatedTime:     getResponse.CreatedOn,
			LastUpdatedTime: getResponse.UpdatedOn,
		}, http.StatusOK)
	}
}
