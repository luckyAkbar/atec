package repository

import "errors"

var (
	ErrNotFound error = errors.New("data not found")
)
