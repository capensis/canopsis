package httperror

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/api/httperror/httperror.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror Responder

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
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

type Responder interface {
	// Respond calls c.Abort() and writes response code and body based on error.
	Respond(c *gin.Context, err error)
	// GetResponse returns response code and body based on error.
	GetResponse(c *gin.Context, err error) (code int, body *fastjson.Value)
}

type responder struct {
	trans  validation.ErrorTranslator
	logger zerolog.Logger
}

func NewConflictError(msg string) *ConflictError {
	return &ConflictError{Message: msg}
}

func NewForbiddenError(msg string) *ForbiddenError {
	return &ForbiddenError{Message: msg}
}

func NewResponder(trans validation.ErrorTranslator, logger zerolog.Logger) Responder {
	return &responder{
		trans:  trans,
		logger: logger,
	}
}

func (r *responder) Respond(c *gin.Context, err error) {
	code, res := r.buildResponse(c, err)
	c.Abort()
	c.Data(code, gin.MIMEJSON, res.MarshalTo(nil))
}

func (r *responder) GetResponse(c *gin.Context, err error) (int, *fastjson.Value) {
	// delegate to private method to adjust CallerSkipFrame(2)
	return r.buildResponse(c, err)
}

func (r *responder) buildResponse(c *gin.Context, err error) (int, *fastjson.Value) {
	var res *fastjson.Value
	valErr := &validation.Error{}
	if errors.As(err, &valErr) {
		res, err = r.getValErrResponse(c, valErr)
		if err == nil {
			return http.StatusBadRequest, res
		}
	}

	var code int
	var msg string
	var invalidReqErr *validation.InvalidRequestBodyError
	var forbiddenErr *ForbiddenError
	var conflictErr *ConflictError
	logLvl := zerolog.DebugLevel
	if errors.Is(err, ErrNotFound) {
		code = http.StatusNotFound
	} else if errors.Is(err, ErrUnauthorized) {
		code = http.StatusUnauthorized
	} else if errors.Is(err, ErrMethodNotAllowed) {
		code = http.StatusMethodNotAllowed
	} else if errors.Is(err, ErrServiceUnavailable) {
		code = http.StatusServiceUnavailable
	} else if errors.Is(err, ErrMaintenance) {
		code = http.StatusServiceUnavailable
		msg = "The service is under maintenance. Please try again later."
	} else if errors.Is(err, ErrRequestEntityTooLarge) {
		code = http.StatusRequestEntityTooLarge
	} else if errors.As(err, &forbiddenErr) {
		code = http.StatusForbidden
		msg = forbiddenErr.Message
	} else if errors.As(err, &conflictErr) {
		code = http.StatusConflict
		msg = conflictErr.Message
	} else if errors.As(err, &invalidReqErr) {
		code = http.StatusBadRequest
		logLvl = zerolog.WarnLevel
	} else if errors.Is(err, authctx.ErrNotFound) {
		code = http.StatusUnauthorized
	} else if errors.Is(err, context.Canceled) {
		code = http.StatusRequestTimeout
		logLvl = zerolog.WarnLevel
	} else if errors.Is(err, context.DeadlineExceeded) || mongodriver.IsTimeout(err) {
		code = http.StatusRequestTimeout
		logLvl = zerolog.ErrorLevel
	} else {
		code = http.StatusInternalServerError
		// c.Error panics if a nil error is passed — log stack trace to debug GetResponse call with nil error
		err = c.Error(err)
		logLvl = zerolog.ErrorLevel
	}

	r.logger.WithLevel(logLvl).
		CallerSkipFrame(2).
		Err(err).
		Str("uri", c.Request.RequestURI).
		Str("method", c.Request.Method).
		Int("response_code", code).
		Msg("unexpected error during http request")

	ar := fastjson.Arena{}
	res = ar.NewObject()
	res.Set("error", ar.NewString(http.StatusText(code)))
	if msg != "" {
		res.Set("message", ar.NewString(msg))
	}

	return code, res
}

// getValErrResponse translates validation errors and returns response body.
func (r *responder) getValErrResponse(c *gin.Context, valErr *validation.Error) (*fastjson.Value, error) {
	locale, err := authctx.GetLocale(c)
	if err != nil {
		r.logger.Debug().Err(err).Msg("cannot get locale from context, use default locale")
	}

	errTrans, err := r.trans.Translate(locale, valErr)
	if err != nil {
		return nil, fmt.Errorf("cannot translate validation error: %w", err)
	}

	ar := fastjson.Arena{}
	errsObj := ar.NewObject()
	for k, v := range errTrans {
		errsObj.Set(k, ar.NewString(v))
	}

	res := ar.NewObject()
	res.Set("errors", errsObj)

	return res, nil
}
