package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rajendragosavi/service-catalog/pkg/errors"
)

// Don’t have to repeat yourself every time you respond to user, instead you can use some helper functions.

func (s Service) respond(w http.ResponseWriter, data interface{}, status int) {
	var respData interface{}
	switch v := data.(type) {
	case nil:
	case errors.ErrorArgument:
		status = http.StatusBadRequest
		respData = ErrorResponse{ErrorMessage: v.Unwrap().Error()}
	case errors.DuplicateKeyError:
		status = http.StatusBadRequest
		respData = ErrorResponse{ErrorMessage: v.Unwrap().Error()}
	case error:
		if http.StatusText(status) == "" {
			status = http.StatusInternalServerError
		} else {
			respData = ErrorResponse{ErrorMessage: v.Error()}
		}
	default:
		respData = data
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if data != nil {
		err := json.NewEncoder(w).Encode(respData)
		if err != nil {
			http.Error(w, "Could not encode in json", http.StatusBadRequest)
			return
		}
	}
}

// it does not read to the memory, instead it will read it to the given 'v' interface.
func (s Service) decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
