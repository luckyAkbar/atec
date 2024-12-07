package usecase

import "errors"

// UsecaseError custom error type wrapper returned by all the usecase function if there is an error. Compliance with error interface
//
//nolint:revive // allow stutters
type UsecaseError struct {
	ErrType error
	Message string
}

// Error returns the error message
func (ue UsecaseError) Error() string {
	return ue.Message
}

// known usecase error type
var (
	ErrBadRequest = errors.New("bad request")

	ErrInternal = errors.New("internal server error")

	ErrUnauthorized = errors.New("unauthorized")

	ErrNotFound = errors.New("not found")
)
