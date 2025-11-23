package auth

import "errors"

var (
	ErrEmailExists = errors.New("email already registered")
)
