package onlygithub

import "github.com/pkg/errors"

var (
	// ErrUnauthorized is returned when a request is unauthorized.
	ErrUnauthorized = errors.New("user is unauthorized")
	// ErrNotFound is returned when a resource is not found.
	ErrNotFound = errors.New("resource not found")
)

// InternalError is an internal error returned from the database.
// These errors should not be exposed to the user.
type InternalError struct {
	err error
}

// WrapInternalError wraps an error in an InternalError.
func WrapInternalError(err error) error {
	return &InternalError{errors.WithStack(err)}
}

func (e *InternalError) Error() string {
	return "internal error"
}

func (e *InternalError) Unwrap() error {
	return e.err
}

// ErrorResponse is a response that contains an error message.
type ErrorResponse struct {
	Message string `json:"message"`
}
