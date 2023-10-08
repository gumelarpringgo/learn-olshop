package common

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrFailCreateData = errors.New("failed create data")
	ErrUnauthorized   = errors.New("user unauthorized")
)
