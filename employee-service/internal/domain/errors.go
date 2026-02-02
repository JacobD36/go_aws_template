package domain

import "errors"

var (
	ErrInvalidName     = errors.New("invalid employee name")
	ErrInvalidEmail    = errors.New("invalid employee email")
	ErrInvalidPassword = errors.New("invalid password: must be at least 8 characters with at least one uppercase letter, one number, and one special character")
	ErrNotFound        = errors.New("employee not found")
)
