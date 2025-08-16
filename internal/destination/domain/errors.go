package domain

import "errors"

var (
	ErrDestinationNotFound      = errors.New("destination not found")
	ErrDestinationAlreadyExists = errors.New("destination already exists")
)
