package httperror

import "errors"

var ErrNotFound = errors.New("not found")
var ErrUnauthorized = errors.New("unauthorized")
var ErrMethodNotAllowed = errors.New("method not allowed")
var ErrRequestEntityTooLarge = errors.New("request entity too large")
var ErrServiceUnavailable = errors.New("service unavailable")
var ErrMaintenance = errors.New("maintenance")

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func NewConflictError(msg string) *ConflictError {
	return &ConflictError{Message: msg}
}

func NewForbiddenError(msg string) *ForbiddenError {
	return &ForbiddenError{Message: msg}
}
