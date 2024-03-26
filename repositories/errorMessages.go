package repositories

import "errors"

var (
	ErrNotFound           = errors.New("record does not exist")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
