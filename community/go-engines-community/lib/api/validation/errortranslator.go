package validation

import (
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/rs/zerolog"
)

type ErrorTranslator interface {
	Translate(locale string, err *Error) ErrorsTranslations
}

type ErrorsTranslations map[string]string

func NewErrorTranslator(trans *ut.UniversalTranslator, logger zerolog.Logger) ErrorTranslator {
	return &errorTranslator{
		trans:  trans,
		logger: logger,
	}
}

type errorTranslator struct {
	trans  *ut.UniversalTranslator
	logger zerolog.Logger
}

func (t *errorTranslator) Translate(locale string, err *Error) ErrorsTranslations {
	if err == nil {
		return nil
	}

	var trans ut.Translator
	if locale != "" {
		var ok bool
		if trans, ok = t.trans.GetTranslator(locale); !ok {
			t.logger.Error().Str("locale", locale).Msg("validation error translator not found, use default locale")
		}
	}

	if trans == nil {
		trans = t.trans.GetFallback()
	}

	res := make(ErrorsTranslations, len(err.errors))
	for _, fe := range err.errors {
		ns := err.TransformNamespace(fe.StructNamespace())
		m := fe.Translate(trans)
		// remove field name from the beginning of the error manually because the lib doesn't allow to control it
		m = strings.TrimPrefix(m, fe.Field()+" ")
		res[ns] = m
	}

	return res
}
