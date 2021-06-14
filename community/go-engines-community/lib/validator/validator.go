package validator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var trans ut.Translator

// RegisterTranslations defines custom validation error messages.
func RegisterTranslations(v *validator.Validate) {
	enTrans := en.New()
	uniTrans := ut.New(enTrans, enTrans)
	trans, _ = uniTrans.GetTranslator("en")

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is missing.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("required_with", trans, func(ut ut.Translator) error {
		return ut.Add("required_with", "{0} is required when {1} is present.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required_with", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("required_or", trans, func(ut ut.Translator) error {
		return ut.Add("required_or", "{0} or {1} is required.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required_or", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("notblank", trans, func(ut ut.Translator) error {
		return ut.Add("notblank", "{0} should not be blank.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("notblank", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("gtfield", trans, func(ut ut.Translator) error {
		return ut.Add("gtfield", "{0} should be greater than {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gtfield", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("gt", trans, func(ut ut.Translator) error {
		return ut.Add("gt", "{0} should be greater than {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gt", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("gte", trans, func(ut ut.Translator) error {
		return ut.Add("gte", "{0} should be greater or equal than {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gte", fe.Field(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("lte", trans, func(ut ut.Translator) error {
		return ut.Add("lte", "{0} should be less or equal than {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("lte", fe.Field(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("max", trans, func(ut ut.Translator) error {
		return ut.Add("max", "{0} should be {1} or less.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("max", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("min", trans, func(ut ut.Translator) error {
		return ut.Add("min", "{0} should be {1} or more.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("rrule", trans, func(ut ut.Translator) error {
		return ut.Add("rrule", "{0} is invalid recurrence rule.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("rrule", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("not_exist", trans, func(ut ut.Translator) error {
		return ut.Add("not_exist", "{0} doesn't exist", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("not_exist", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("unique", trans, func(ut ut.Translator) error {
		return ut.Add("unique", "{0} already exists", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("filter", trans, func(ut ut.Translator) error {
		return ut.Add("filter", "{0} is not valid mongo query", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("filter", fe.Field())
		return t
	})
	_ = v.RegisterTranslation("oneof", trans, func(ut ut.Translator) error {
		return ut.Add("oneof", "{0} must be one of [{1}]", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("oneof", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("oneoforempty", trans, func(ut ut.Translator) error {
		return ut.Add("oneoforempty", "{0} must be one of [{1}]", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("oneoforempty", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("entitypattern_invalid", trans, func(ut ut.Translator) error {
		return ut.Add("entitypattern_invalid", "Invalid entity pattern list.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("entitypattern_invalid", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("entitypattern_contains_empty", trans, func(ut ut.Translator) error {
		return ut.Add("entitypattern_contains_empty", "entity pattern list contains an empty pattern.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("entitypattern_contains_empty", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("alarmpattern_invalid", trans, func(ut ut.Translator) error {
		return ut.Add("alarmpattern_invalid", "Invalid alarm pattern list.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alarmpattern_invalid", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("alarmpattern_contains_empty", trans, func(ut ut.Translator) error {
		return ut.Add("alarmpattern_contains_empty", "alarm pattern list contains an empty pattern.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alarmpattern_contains_empty", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("entityfilter", trans, func(ut ut.Translator) error {
		return ut.Add("entityfilter", "{0} is invalid entity filter.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("entityfilter", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("url", trans, func(ut ut.Translator) error {
		return ut.Add("url", "{0} is not an url", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("url", fe.Value().(string))
		return t
	})
	_ = v.RegisterTranslation("json", trans, func(ut ut.Translator) error {
		return ut.Add("json", "{0} is not a valid json", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("json", fe.Field())
		return t
	})
	_ = v.RegisterTranslation("template", trans, func(ut ut.Translator) error {
		return ut.Add("template", "{0} is not a valid template", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("template", fe.Field())
		return t
	})
}

// TranslateError returns custom validation error message.
func TranslateError(fe validator.FieldError) string {
	return fe.Translate(trans)
}
