package repository

import "errors"

// known errors that might be returned by this repository's functions
var (
	ErrNotFound = errors.New("data not found")
)
