package httperror

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/api/httperror/httperror.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror Responder

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrNotFound = errors.New("not found")

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

func NewResponder(trans validation.ErrorTranslator, logger zerolog.Logger) Responder {
	return &responder{
		trans:  trans,
		logger: logger,
	}
}

func (r *responder) Respond(c *gin.Context, err error) {
	code, res := r.GetResponse(c, err)
	c.Abort()
	c.Data(code, gin.MIMEJSON, res.MarshalTo(nil))
}

func (r *responder) GetResponse(c *gin.Context, err error) (int, *fastjson.Value) {
	valErr := &validation.Error{}
	if errors.As(err, &valErr) {
		return http.StatusBadRequest, r.getValErrResponse(c, valErr)
	}

	var code int
	logLvl := zerolog.DebugLevel
	if errors.Is(err, ErrNotFound) {
		code = http.StatusNotFound
	} else if errors.Is(err, validation.ErrInvalidRequestBody) {
		code = http.StatusBadRequest
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
		Err(err).
		Str("uri", c.Request.RequestURI).
		Str("method", c.Request.Method).
		Int("response_code", code).
		Msg("unexpected error during http request")

	ar := fastjson.Arena{}
	res := ar.NewObject()
	res.Set("error", ar.NewString(http.StatusText(code)))

	return code, res
}

// getValErrResponse translates validation errors and returns response body.
func (r *responder) getValErrResponse(c *gin.Context, valErr *validation.Error) *fastjson.Value {
	locale, err := authctx.GetLocale(c)
	if err != nil {
		r.logger.Debug().Err(err).Msg("cannot get locale from context, use default locale")
	}

	errTrans := r.trans.Translate(locale, valErr)
	ar := fastjson.Arena{}
	errsObj := ar.NewObject()
	for k, v := range errTrans {
		errsObj.Set(k, ar.NewString(v))
	}

	res := ar.NewObject()
	res.Set("errors", errsObj)

	return res
}
