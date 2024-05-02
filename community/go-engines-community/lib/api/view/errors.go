package view

import "errors"

type ValidationError struct {
	field string
	error error
}

func (v ValidationError) Error() string {
	return v.error.Error()
}

var (
	ErrViewsNotFound  = errors.New("views not found")
	ErrValueIsMissing = errors.New("value is missing")
)
