package validation

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

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

func Bind(c *gin.Context, v any) error {
	if err := c.ShouldBind(v); err != nil {
		return convertValErr(err, v)
	}

	return nil
}

func BindQuery(c *gin.Context, v any) error {
	if err := c.ShouldBindQuery(v); err != nil {
		return convertValErr(err, v)
	}

	return nil
}

func ValidateStruct(v any) error {
	if err := binding.Validator.ValidateStruct(v); err != nil {
		return convertValErr(err, v)
	}

	return nil
}

func convertValErr(err error, v any) error {
	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		return NewError(valErr, v)
	}

	return NewInvalidRequestBodyError(err)
}
