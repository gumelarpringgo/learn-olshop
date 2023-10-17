package common

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrFailedCreateData = errors.New("failed create data")
	ErrFailedUpdateData = errors.New("failed update data")
	ErrUnauthorized     = errors.New("user unauthorized")
	ErrUploadFile       = errors.New("failed upload data")
	ErrDeleteData       = errors.New("failed delete data")
	ErrMustHavePrimary  = errors.New("must have primary")
)
