package api

import (
	"encoding/json"
	"net/http"

	"libdb.so/onlygithub/onlygithub"
)

// RespondError responds to the client with an error message.
func RespondError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(onlygithub.ErrorResponse{Message: err.Error()})
}

// Respond responds to the client with data. If data is nil, the response will
// be a 204 No Content response.
func Respond(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if data == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(data)
}
