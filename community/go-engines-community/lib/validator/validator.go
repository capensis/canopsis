package validator

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const InvalidIDChars = "/?.$"

var trans ut.Translator

// RegisterTranslations defines custom validation error messages.
func RegisterTranslations(v *validator.Validate) {
	enTrans := en.New()
	uniTrans := ut.New(enTrans, enTrans)
	trans, _ = uniTrans.GetTranslator("en")

	_ = v.RegisterTranslation("required_if", trans, func(ut ut.Translator) error {
		return ut.Add("required_if", "{0} is required when {1} is defined.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required_if", fe.StructField(), fe.Param())
		return t
	})
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
	_ = v.RegisterTranslation("required_without", trans, func(ut ut.Translator) error {
		return ut.Add("required_without", "{0} is required when {1} is not present.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required_without", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("required_not_both", trans, func(ut ut.Translator) error {
		return ut.Add("required_not_both", "Can't be present both {0} and {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required_not_both", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("required_or", trans, func(ut ut.Translator) error {
		return ut.Add("required_or", "{0} or {1} is required.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required_or", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("required_not_both", trans, func(ut ut.Translator) error {
		return ut.Add("required_not_both", "Can't be present both {0} and {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required_not_both", fe.StructField(), fe.Param())
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
	_ = v.RegisterTranslation("gtefield", trans, func(ut ut.Translator) error {
		return ut.Add("gtefield", "{0} should be greater or equal than {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gtefield", fe.StructField(), fe.Param())
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
		t, _ := ut.T("gte", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("ltfield", trans, func(ut ut.Translator) error {
		return ut.Add("ltfield", "{0} should be less than {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ltfield", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("ltefield", trans, func(ut ut.Translator) error {
		return ut.Add("ltefield", "{0} should be less or equal than {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ltefield", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("lt", trans, func(ut ut.Translator) error {
		return ut.Add("lt", "{0} should be less than {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("lt", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("lte", trans, func(ut ut.Translator) error {
		return ut.Add("lte", "{0} should be less or equal than {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("lte", fe.StructField(), fe.Param())
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
		return ut.Add("not_exist", "{0} doesn't exist.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("not_exist", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("not_approver", trans, func(ut ut.Translator) error {
		return ut.Add("not_approver", "{0} doesn't have approve rights or doesn't exist.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("not_approver", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("unique", trans, func(ut ut.Translator) error {
		return ut.Add("unique", "{0} already exists.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("oneof", trans, func(ut ut.Translator) error {
		return ut.Add("oneof", "{0} must be one of [{1}].", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("oneof", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("oneoforempty", trans, func(ut ut.Translator) error {
		return ut.Add("oneoforempty", "{0} must be one of [{1}] or empty.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("oneoforempty", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("url", trans, func(ut ut.Translator) error {
		return ut.Add("url", "{0} is not an url.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("url", fe.Value().(string))
		return t
	})
	_ = v.RegisterTranslation("json", trans, func(ut ut.Translator) error {
		return ut.Add("json", "{0} is not a valid json.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("json", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("template", trans, func(ut ut.Translator) error {
		return ut.Add("template", "{0} is not a valid template.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("template", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("has_duplicates", trans, func(ut ut.Translator) error {
		return ut.Add("has_duplicates", "{0} contains duplicate values.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("has_duplicates", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("has_duplicates_with", trans, func(ut ut.Translator) error {
		return ut.Add("has_duplicates_with", "{0} contains duplicate values with {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("has_duplicates_with", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("id", trans, func(ut ut.Translator) error {
		return ut.Add("id", fmt.Sprintf("{0} cannot contain %s characters.", InvalidIDChars), true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("id", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("must_be_empty", trans, func(ut ut.Translator) error {
		return ut.Add("must_be_empty", "{0} is not empty.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("must_be_empty", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("ltefield", trans, func(ut ut.Translator) error {
		return ut.Add("ltefield", "{0} should be less or equal than {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ltefield", fe.StructField(), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} is not a valid email address.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("filemask", trans, func(ut ut.Translator) error {
		return ut.Add("filemask", "{0} is not a valid file mask.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("filemask", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("iscolor", trans, func(ut ut.Translator) error {
		return ut.Add("iscolor", "{0} is not valid.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("iscolor", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("iscolororempty", trans, func(ut ut.Translator) error {
		return ut.Add("iscolororempty", "{0} is not valid.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("iscolororempty", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("metaalarm_config_invalid", trans, func(ut ut.Translator) error {
		return ut.Add("metaalarm_config_invalid", "Config is not a dict.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("metaalarm_config_invalid", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("metaalarm_config_invalid_with_type", trans, func(ut ut.Translator) error {
		return ut.Add("metaalarm_config_invalid_with_type", "{0} config can not be in type {1}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("metaalarm_config_invalid_with_type", fe.Value().(string), fe.Param())
		return t
	})
	_ = v.RegisterTranslation("time_format", trans, func(ut ut.Translator) error {
		return ut.Add("time_format", "{0} is invalid time format.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("time_format", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("multi_sort_invalid", trans, func(ut ut.Translator) error {
		return ut.Add("multi_sort_invalid", "Invalid multi_sort value.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("multi_sort_invalid", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("regexp", trans, func(ut ut.Translator) error {
		return ut.Add("regexp", "{0} is invalid regexp.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("regexp", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("event_pattern", trans, func(ut ut.Translator) error {
		return ut.Add("event_pattern", "{0} is invalid event pattern.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("event_pattern", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("alarm_pattern", trans, func(ut ut.Translator) error {
		return ut.Add("alarm_pattern", "{0} is invalid alarm pattern.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alarm_pattern", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("entity_pattern", trans, func(ut ut.Translator) error {
		return ut.Add("entity_pattern", "{0} is invalid entity pattern.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("entity_pattern", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("pbehavior_pattern", trans, func(ut ut.Translator) error {
		return ut.Add("pbehavior_pattern", "{0} is invalid pbehavior pattern.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("pbehavior_pattern", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("weather_service_pattern", trans, func(ut ut.Translator) error {
		return ut.Add("weather_service_pattern", "{0} is invalid weather service pattern.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("weather_service_pattern", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("info_value", trans, func(ut ut.Translator) error {
		return ut.Add("info_value", types.ErrInvalidInfoType.Error(), true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("info_value", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("latitude", trans, func(ut ut.Translator) error {
		return ut.Add("latitude", "{0} must contain valid latitude coordinates.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("latitude", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("longitude", trans, func(ut ut.Translator) error {
		return ut.Add("longitude", "{0} must contain valid longitude coordinates.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("longitude", fe.StructField())
		return t
	})
	_ = v.RegisterTranslation("invalid", trans, func(ut ut.Translator) error {
		return ut.Add("invalid", "{0} is invalid.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("invalid", fe.StructField())
		return t
	})
}

// TranslateError returns custom validation error message.
func TranslateError(fe validator.FieldError) string {
	return fe.Translate(trans)
}
