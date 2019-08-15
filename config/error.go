package config

import "errors"

var (
	ErrDefault = errors.New("System error occurred, Please contact us.")
	ErrNotFound = errors.New("Resource not found")
)
