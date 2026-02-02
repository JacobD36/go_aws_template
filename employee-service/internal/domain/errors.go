package domain

import "errors"

var (
	ErrInvalidName  = errors.New("invalid employee name")
	ErrInvalidEmail = errors.New("invalid employee email")
	ErrNotFound     = errors.New("employee not found")
)
