package httperror

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/rs/zerolog"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

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

func MapError(err error) (httpCode int, httpMsg string, logLvl zerolog.Level) {
	if fe, ok := errors.AsType[*ForbiddenError](err); ok {
		return http.StatusForbidden, fe.Message, zerolog.DebugLevel
	}

	if ce, ok := errors.AsType[*ConflictError](err); ok {
		return http.StatusConflict, ce.Message, zerolog.DebugLevel
	}

	if _, ok := errors.AsType[*validation.InvalidRequestBodyError](err); ok {
		return http.StatusBadRequest, "", zerolog.WarnLevel
	}

	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, "", zerolog.DebugLevel
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, "", zerolog.DebugLevel
	case errors.Is(err, ErrMethodNotAllowed):
		return http.StatusMethodNotAllowed, "", zerolog.DebugLevel
	case errors.Is(err, ErrServiceUnavailable):
		return http.StatusServiceUnavailable, "", zerolog.DebugLevel
	case errors.Is(err, ErrMaintenance):
		return http.StatusServiceUnavailable, "The service is under maintenance. Please try again later.", zerolog.DebugLevel
	case errors.Is(err, ErrRequestEntityTooLarge):
		return http.StatusRequestEntityTooLarge, "", zerolog.DebugLevel
	case errors.Is(err, context.Canceled):
		return http.StatusRequestTimeout, "", zerolog.WarnLevel
	case errors.Is(err, context.DeadlineExceeded) || mongodriver.IsTimeout(err):
		return http.StatusRequestTimeout, "", zerolog.ErrorLevel
	default:
		return http.StatusInternalServerError, "", zerolog.ErrorLevel
	}
}
