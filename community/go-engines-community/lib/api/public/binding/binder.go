package binding

import (
	"context"
	"errors"
	"io"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type Binder interface {
	Bind(*http.Request, any) error
	BindForm(r *http.Request, v any) error
	URLParam(r *http.Request, k string) string
	ValidateStruct(ctx context.Context, v any) error
}

func NewBinder(decoder encoding.Decoder, validator *validator.Validate) Binder {
	b := &binder{
		decoder:    decoder,
		validator:  validator,
		urlDecoder: schema.NewDecoder(),
	}
	b.urlDecoder.SetAliasTag("form")

	return b
}

type binder struct {
	decoder    encoding.Decoder
	validator  *validator.Validate
	urlDecoder *schema.Decoder
}

func (b *binder) Bind(r *http.Request, v any) error {
	emptyErr := validation.NewInvalidRequestBodyError(errors.New("body is empty"))
	if r.Body == nil {
		return emptyErr
	}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if len(buf) == 0 {
		return emptyErr
	}

	err = b.decoder.Decode(buf, v)
	if err != nil {
		return err
	}

	err = b.validator.StructCtx(r.Context(), v)
	if err != nil {
		return validation.ConvertValErr(err, v)
	}

	return nil
}

func (b *binder) BindForm(r *http.Request, v any) error {
	err := b.urlDecoder.Decode(v, r.URL.Query())
	if err != nil {
		var schemaErr schema.MultiError
		if errors.As(err, &schemaErr) {
			return b.transformSchemaErr(schemaErr)
		}

		return err
	}

	err = b.validator.StructCtx(r.Context(), v)
	if err != nil {
		return validation.ConvertValErr(err, v)
	}

	return nil
}

func (b *binder) URLParam(r *http.Request, k string) string {
	return chi.URLParam(r, k)
}

func (b *binder) ValidateStruct(ctx context.Context, v any) error {
	err := b.validator.StructCtx(ctx, v)
	if err != nil {
		return validation.ConvertValErr(err, v)
	}

	return nil
}

func (b *binder) transformSchemaErr(schemaErr schema.MultiError) error {
	var emptyErr schema.EmptyFieldError
	var unknownErr schema.UnknownKeyError
	valErrs := make(validator.ValidationErrors, 0, len(schemaErr))
	for f, fe := range schemaErr {
		var vallErr validator.FieldError
		if errors.As(fe, &emptyErr) {
			vallErr = validation.NewFieldError("required", f, f)
		} else if errors.As(fe, &unknownErr) {
			vallErr = validation.NewFieldError("unsupported_field", f, f)
		} else {
			vallErr = validation.NewFieldError("invalid", f, f)
		}

		valErrs = append(valErrs, vallErr)
	}

	return validation.NewError(valErrs, nil)
}
