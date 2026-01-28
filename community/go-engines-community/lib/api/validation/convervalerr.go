package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ConvertValErr(err error, v any) error {
	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		return NewError(valErr, v)
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
