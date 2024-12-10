package common

import (
	validate "github.com/go-playground/validator/v10"
)

// Validator is singleton to be used to mainly validating struct values
var Validator = validate.New()
