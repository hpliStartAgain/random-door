package service

import "errors"

var (
	ErrInvalidParam = errors.New("invalid parameter")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrPermission   = errors.New("permission denied")
)

type clientError struct {
	kind    error
	message string
}

func (e *clientError) Error() string {
	return e.message
}

func (e *clientError) Unwrap() error {
	return e.kind
}

func invalidParam(message string) error {
	return &clientError{kind: ErrInvalidParam, message: message}
}

func notFound(message string) error {
	return &clientError{kind: ErrNotFound, message: message}
}

func conflict(message string) error {
	return &clientError{kind: ErrConflict, message: message}
}

func permissionDenied(message string) error {
	return &clientError{kind: ErrPermission, message: message}
}

// ClientMessage returns the stable public message for classified service errors.
func ClientMessage(err error) string {
	var clientErr *clientError
	if errors.As(err, &clientErr) {
		return clientErr.message
	}
	return "internal server error"
}
