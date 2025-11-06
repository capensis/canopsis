package common

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/validator"
	"github.com/go-playground/validator/v10"
	"github.com/valyala/fastjson"
)

const LimitLinkedRules = 11

// ErrorResponse is base failed response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewErrorResponse
// @Failure 500 {object} ErrorResponse
func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Error: err.Error()}
}

var NotFoundResponse = ErrorResponse{Error: "Not found"}
var MethodNotAllowedResponse = ErrorResponse{Error: http.StatusText(http.StatusMethodNotAllowed)}
var UnauthorizedResponse = ErrorResponse{Error: http.StatusText(http.StatusUnauthorized)}
var InternalServerErrorResponse = ErrorResponse{Error: "Internal server error"}
var ForbiddenResponse = ErrorResponse{Error: http.StatusText(http.StatusForbidden)}
var RequestTimeoutResponse = ErrorResponse{Error: http.StatusText(http.StatusRequestTimeout)}
var CanopsisUnderMaintenanceResponse = ErrorResponse{Error: "canopsis is under maintenance"}

// ValidationErrorResponse is response for failed validation.
type ValidationErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

// NewValidationErrorResponse creates response by validation errors.
// @Failure 400 {object} ValidationErrorResponse
func NewValidationErrorResponse(err error, request interface{}) interface{} {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		return TransformValidationErrors(errs, request)
	}

	return ErrorResponse{Error: "request has invalid structure"}
}

func TransformValidationErrors(errs validator.ValidationErrors, request interface{}) ValidationErrorResponse {
	var res ValidationErrorResponse
	res.Errors = make(map[string]string)
	for _, fe := range errs {
		field := transformNamespace(fe.StructNamespace(), request)
		res.Errors[field] = libvalidator.TranslateError(fe)
	}

	return res
}

func NewValidationErrorFastJsonValue(ar *fastjson.Arena, err error, request interface{}) *fastjson.Value {
	var validatorErrs validator.ValidationErrors
	if errors.As(err, &validatorErrs) {
		value := ar.NewObject()
		for _, fe := range validatorErrs {
			field := transformNamespace(fe.StructNamespace(), request)
			value.Set(field, ar.NewString(libvalidator.TranslateError(fe)))
		}

		return value
	}

	var commonValidatorErrs ValidationError
	if errors.As(err, &commonValidatorErrs) {
		value := ar.NewObject()
		for k, v := range commonValidatorErrs.ValidationErrorResponse().Errors {
			value.Set(k, ar.NewString(v))
		}

		return value
	}

	return ar.NewString("request has invalid structure")
}

// transformNamespace prepares field namespace for response.
// for example:
// - Username -> username
// - Items[0] -> items.0
// - Items[0].Name -> items.0.name
func transformNamespace(namespace string, request interface{}) string {
	re := regexp.MustCompile(`(\.*)\[([^\]]+)\](\.*)`)
	// remove brackets
	namespace = re.ReplaceAllStringFunc(namespace, func(s string) string {
		s = strings.ReplaceAll(s, "[", "")
		s = strings.ReplaceAll(s, "]", "")
		if s[0] != '.' {
			s = "." + s
		}
		if s[len(s)-1] != '.' {
			s = s + "."
		}

		return s
	})
	// replace name to json tag name
	path := strings.Split(namespace, ".")
	path = path[1:]
	val := reflect.ValueOf(request)
loop:
	for i := range path {
		k := val.Kind()

		switch k {
		case reflect.Interface, reflect.Ptr:
			val = val.Elem()
			k = val.Kind()
		}

		switch k {
		case reflect.Struct:
			if f, ok := val.Type().FieldByName(path[i]); ok {
				tag := f.Tag.Get("json")
				if tag == "" {
					tag = f.Tag.Get("form")
				}

				tags := strings.Split(tag, ",")
				if len(tags) > 1 && tags[len(tags)-1] == "omitempty" {
					tag = strings.Join(tags[:len(tags)-1], ",")
				}
				if tag == "-" {
					tag = strings.ToLower(path[i])
				}
				val = val.FieldByName(path[i])
				path[i] = tag
			}
		case reflect.Slice, reflect.Array:
			index, err := strconv.Atoi(path[i])
			if err != nil {
				break loop
			}
			val = val.Index(index)
		case reflect.Map:
			nameVal := reflect.ValueOf(path[i])
			if !nameVal.Type().AssignableTo(val.Type().Key()) {
				break loop
			}
			val = val.MapIndex(nameVal)
		}
	}

	res := make([]string, 0)
	for _, p := range path {
		if p != "" {
			res = append(res, p)
		}
	}

	return strings.Join(res, ".")
}

type ValidationError struct {
	errMsgs map[string]string
}

func NewValidationError(field, errMsg string) ValidationError {
	return ValidationError{errMsgs: map[string]string{field: errMsg}}
}

func NewValidationErrors(errMsgs map[string]string) ValidationError {
	return ValidationError{errMsgs: errMsgs}
}

func (v ValidationError) Error() string {
	b := strings.Builder{}
	i := 0
	for field, errMsg := range v.errMsgs {
		b.WriteString(field)
		b.WriteString(": ")
		b.WriteString(errMsg)
		if i < len(v.errMsgs)-1 {
			b.WriteString("; ")
		}

		i++
	}

	return b.String()
}

func (v ValidationError) AddFieldPrefix(p string) ValidationError {
	errMsgs := make(map[string]string, len(v.errMsgs))
	for f, m := range v.errMsgs {
		errMsgs[p+"."+f] = m
	}

	return ValidationError{errMsgs: errMsgs}
}

func (v ValidationError) ValidationErrorResponse() ValidationErrorResponse {
	return ValidationErrorResponse{
		Errors: v.errMsgs,
	}
}

type AlarmStep struct {
	Type      string            `bson:"_t" json:"_t"`
	Timestamp *datetime.CpsTime `bson:"t" json:"t" swaggertype:"integer"`
	Author    string            `bson:"a" json:"a"`
	UserID    string            `bson:"user_id,omitempty" json:"user_id"`
	Message   string            `bson:"m" json:"m"`
	Value     types.CpsNumber   `bson:"val" json:"val"`
	Initiator string            `bson:"initiator" json:"initiator"`
	Execution string            `bson:"exec,omitempty" json:"-"`
	IconName  string            `bson:"icon_name,omitempty" json:"icon_name,omitempty"`
	Color     string            `bson:"color,omitempty" json:"color,omitempty"`

	// Ticket related fields
	types.TicketInfo `bson:",inline"`

	InPbehaviorInterval bool `bson:"in_pbh,omitempty" json:"in_pbh,omitempty"`
}
