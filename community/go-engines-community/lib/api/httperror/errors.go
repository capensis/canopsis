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
	var invalidReqErr *validation.InvalidRequestBodyError
	var forbiddenErr *ForbiddenError
	var conflictErr *ConflictError
	logLvl = zerolog.DebugLevel
	if errors.Is(err, ErrNotFound) {
		httpCode = http.StatusNotFound
	} else if errors.Is(err, ErrUnauthorized) {
		httpCode = http.StatusUnauthorized
	} else if errors.Is(err, ErrMethodNotAllowed) {
		httpCode = http.StatusMethodNotAllowed
	} else if errors.Is(err, ErrServiceUnavailable) {
		httpCode = http.StatusServiceUnavailable
	} else if errors.Is(err, ErrMaintenance) {
		httpCode = http.StatusServiceUnavailable
		httpMsg = "The service is under maintenance. Please try again later."
	} else if errors.Is(err, ErrRequestEntityTooLarge) {
		httpCode = http.StatusRequestEntityTooLarge
	} else if errors.As(err, &forbiddenErr) {
		httpCode = http.StatusForbidden
		httpMsg = forbiddenErr.Message
	} else if errors.As(err, &conflictErr) {
		httpCode = http.StatusConflict
		httpMsg = conflictErr.Message
	} else if errors.As(err, &invalidReqErr) {
		httpCode = http.StatusBadRequest
		logLvl = zerolog.WarnLevel
	} else if errors.Is(err, context.Canceled) {
		httpCode = http.StatusRequestTimeout
		logLvl = zerolog.WarnLevel
	} else if errors.Is(err, context.DeadlineExceeded) || mongodriver.IsTimeout(err) {
		httpCode = http.StatusRequestTimeout
		logLvl = zerolog.ErrorLevel
	} else {
		httpCode = http.StatusInternalServerError
		logLvl = zerolog.ErrorLevel
	}

	return httpCode, httpMsg, logLvl
}
