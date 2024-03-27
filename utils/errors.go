package utils

import "errors"

var (
	ErrNotFound                = errors.New("record does not exist")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrInvalidRangeForPriority = errors.New("priority must be between Low [0] and High [2]")
	ErrInvalidTypeForPriority  = errors.New("priority must be an integer")
	ErrInvalidValueForPriority = errors.New("invalid value")
)
