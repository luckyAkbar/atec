package usecase

import "errors"

type UsecaseError struct {
	ErrType error
	Message string
}

func (ue UsecaseError) Error() string {
	return ue.Message
}

var (
	ErrBadRequest error = errors.New("bad request")

	ErrInternal = errors.New("internal server error")
)
