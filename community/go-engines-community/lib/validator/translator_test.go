package validator_test

import (
	"errors"
	"strings"
	"testing"

	uivalidation "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	libvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/validator"
	"github.com/go-playground/validator/v10"
)

func TestValidationTagTranslationsAndFallback(t *testing.T) {
	validate := validator.New()
	for _, tag := range []string{"invalid", "value_string", "value_string_array", "untranslated_test_tag"} {
		if err := validate.RegisterValidation(tag, func(validator.FieldLevel) bool { return false }); err != nil {
			t.Fatalf("RegisterValidation(%q) error = %v", tag, err)
		}
	}
	trans, err := libvalidator.NewTranslator(validate, uivalidation.InvalidIDChars)
	if err != nil {
		t.Fatalf("NewTranslator() error = %v", err)
	}

	for _, locale := range []string{"en", "fr"} {
		translator, ok := trans.GetTranslator(locale)
		if !ok {
			t.Fatalf("translator for %q not found", locale)
		}
		for _, tag := range []string{"invalid", "value_string", "value_string_array"} {
			t.Run(locale+"/"+tag, func(t *testing.T) {
				fieldErr := validationFieldError(t, validate, tag)
				assertTranslated(t, "validator FieldError "+tag, fieldErr.Translate(translator), fieldErr.Error())
				customFieldErr := validation.NewFieldError(tag, "Formats", "formats")
				assertTranslated(t, "api validation fieldError "+tag, customFieldErr.Translate(translator), customFieldErr.Error())
			})
		}

		t.Run(locale+"/unregistered-tag-falls-back", func(t *testing.T) {
			fieldErr := validationFieldError(t, validate, "untranslated_test_tag")
			translated := fieldErr.Translate(translator)
			if translated != fieldErr.Error() {
				t.Fatalf("unregistered tag translated to %q, want fallback %q", translated, fieldErr.Error())
			}
			customFieldErr := validation.NewFieldError("untranslated_test_tag", "Formats", "formats")
			if got := customFieldErr.Translate(translator); got != customFieldErr.Error() {
				t.Fatalf("custom unregistered tag translated to %q, want fallback %q", got, customFieldErr.Error())
			}
		})
	}
}

func assertTranslated(t *testing.T, label, translated, fallback string) {
	t.Helper()
	if strings.TrimSpace(translated) == "" || translated == fallback {
		t.Fatalf("%s = %q, want non-empty registered translation instead of fallback %q", label, translated, fallback)
	}
}

func validationFieldError(t *testing.T, validate *validator.Validate, tag string) validator.FieldError {
	t.Helper()
	err := validate.Var("sample", tag)
	if err == nil {
		t.Fatalf("validate.Var() unexpectedly accepted tag %q", tag)
	}
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) || len(errs) != 1 {
		t.Fatalf("validate.Var() error = %v, want one validation error", err)
	}
	return errs[0]
}
