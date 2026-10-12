package response

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../../mocks/lib/api/ui/response/responder.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response ErrorResponder

import (
	"errors"
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
)

type ErrorResponder interface {
	// Respond calls c.Abort() and writes response code and body based on error.
	Respond(c *gin.Context, err error)
	// GetResponse returns response code and body based on error.
	GetResponse(c *gin.Context, err error) (code int, body *fastjson.Value)
}

type responder struct {
	trans  validation.ErrorTranslator
	logger zerolog.Logger
}

func NewErrorResponder(trans validation.ErrorTranslator, logger zerolog.Logger) ErrorResponder {
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
	var logLvl zerolog.Level
	if errors.Is(err, authctx.ErrNotFound) {
		code = http.StatusUnauthorized
		logLvl = zerolog.DebugLevel
	} else {
		code, msg, logLvl = httperror.MapError(err)
	}

	if code == http.StatusInternalServerError {
		// c.Error panics if a nil error is passed — log stack trace to debug GetResponse call with nil error
		err = c.Error(err)
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
