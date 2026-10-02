package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ConvertValErr(err error, v any) error {
	// Keep errors already built by a custom decoder: their namespace is not bound to v.
	if valErr, ok := errors.AsType[*Error](err); ok {
		return valErr
	}

	if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
		return NewError(ve, v)
	}

	return NewInvalidRequestBodyError(err)
}

type InvalidRequestBodyError struct {
	Err error
}

func (e *InvalidRequestBodyError) Error() string {
	return fmt.Sprintf("invalid request body: %+v", e.Err)
}

func (e *InvalidRequestBodyError) Unwrap() error {
	return e.Err
}

func NewInvalidRequestBodyError(err error) *InvalidRequestBodyError {
	return &InvalidRequestBodyError{Err: err}
}
