package response

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../../mocks/lib/api/public/response/responder.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response Responder

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
)

const (
	MIMEJSON = "application/json"
)

type Responder interface {
	// Respond writes response code and body.
	Respond(w http.ResponseWriter, r *http.Request, code int, b any)
	// RespondError writes response code and body based on error.
	RespondError(w http.ResponseWriter, r *http.Request, err error)
	// ErrorResponse returns response code and body based on error.
	ErrorResponse(r *http.Request, err error) (code int, body *fastjson.Value)
}

type responder struct {
	encoder encoding.Encoder
	trans   validation.ErrorTranslator
	logger  zerolog.Logger
}

func NewResponder(encoder encoding.Encoder, trans validation.ErrorTranslator, logger zerolog.Logger) Responder {
	return &responder{
		encoder: encoder,
		trans:   trans,
		logger:  logger,
	}
}

// Respond writes response code and body.
func (r *responder) Respond(w http.ResponseWriter, req *http.Request, code int, body any) {
	var b []byte
	var err error
	if body != nil {
		b, err = r.encoder.Encode(body)
		if err != nil {
			var res *fastjson.Value
			code, res = r.buildErrorResponse(req, err)
			b = res.MarshalTo(nil)
		}
	}

	w.Header().Set("Content-Type", MIMEJSON)
	w.WriteHeader(code)
	if len(b) > 0 {
		_, err = w.Write(b)
		if err != nil {
			r.logger.Err(err).Msg("failed to write response")
		}
	}
}

func (r *responder) RespondError(w http.ResponseWriter, req *http.Request, handlerErr error) {
	code, res := r.buildErrorResponse(req, handlerErr)
	w.Header().Set("Content-Type", MIMEJSON)
	w.WriteHeader(code)
	_, err := w.Write(res.MarshalTo(nil))
	if err != nil {
		r.logger.Err(err).Msg("failed to write response")
	}
}

func (r *responder) ErrorResponse(req *http.Request, err error) (int, *fastjson.Value) {
	// delegate to private method to adjust CallerSkipFrame(2)
	return r.buildErrorResponse(req, err)
}

func (r *responder) buildErrorResponse(req *http.Request, err error) (int, *fastjson.Value) {
	var res *fastjson.Value
	valErr := &validation.Error{}
	if errors.As(err, &valErr) {
		res, err = r.getValErrResponse(req.Context(), valErr)
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

	r.logger.WithLevel(logLvl).
		CallerSkipFrame(2).
		Err(err).
		Str("uri", req.RequestURI).
		Str("method", req.Method).
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
func (r *responder) getValErrResponse(ctx context.Context, valErr *validation.Error) (*fastjson.Value, error) {
	locale, err := authctx.GetLocale(ctx)
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
