package api

import (
	"encoding/json"
	"net/http"

	"github.com/diamondburned/hrt"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"libdb.so/onlygithub"
)

// RespondError responds to the client with an error message.
func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	err = WrapError(r, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(hrt.ErrorHTTPStatus(err, 500))
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
	json.NewEncoder(w).Encode(data)
}

var errInternal = hrt.NewHTTPError(500, "Internal Server Error")

// ExtractError extracts an error code and message from an error.
func ExtractError(r *http.Request, err error) (int, string) {
	err = WrapError(r, err)
	return hrt.ErrorHTTPStatus(err, 500), err.Error()
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// WrapError wraps an error into an HTTP status code. If the error is not
// recognized, it will be wrapped into a 500 Internal Server Error.
func WrapError(r *http.Request, err error) error {
	log := hclog.FromContext(r.Context())

	var internalError *onlygithub.InternalError
	if errors.As(err, &internalError) {
		args := []any{"error", internalError.Unwrap()}
		if st, ok := internalError.Unwrap().(stackTracer); ok {
			args = append(args, "stack", st.StackTrace())
		}

		log.Warn("internal error:", args...)
		return errInternal
	}

	if hrt.ErrorHTTPStatus(err, 0) == 0 {
		switch {
		case errors.Is(err, onlygithub.ErrNotFound):
			err = hrt.WrapHTTPError(404, err)
		case errors.Is(err, onlygithub.ErrUnauthorized):
			err = hrt.WrapHTTPError(401, err)
		}
	}

	if hrt.ErrorHTTPStatus(err, 0) >= 500 {
		log.Warn("internal error:", "error", err)
		err = errInternal
	}

	return err
}
