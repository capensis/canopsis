package common

import (
	"errors"
	"math"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/validator"
	"github.com/go-playground/validator/v10"
	"github.com/valyala/fastjson"
)

// PaginatedMeta is meta for paginated list data.
type PaginatedMeta struct {
	Page       int64 `json:"page"`
	PerPage    int64 `json:"per_page"`
	PageCount  int64 `json:"page_count"`
	TotalCount int64 `json:"total_count"`
}

// PaginatedListResponse is response for paginated list data.
type PaginatedListResponse struct {
	Data interface{}   `json:"data"`
	Meta PaginatedMeta `json:"meta"`
}

// PaginatedData provides access to inner data and total count
type PaginatedData interface {
	GetData() interface{}
	GetTotal() int64
}

func NewPaginatedResponse(q pagination.Query, d PaginatedData) (PaginatedListResponse, error) {
	meta, err := NewPaginatedMeta(q, d.GetTotal())
	if err != nil {
		return PaginatedListResponse{}, err
	}

	data := d.GetData()
	if data == nil {
		data = []interface{}{}
	}

	return PaginatedListResponse{
		Data: data,
		Meta: meta,
	}, nil
}

func NewPaginatedMeta(q pagination.Query, total int64) (PaginatedMeta, error) {
	if !q.Paginate {
		q.Limit = total
	}

	var pageCount int64
	if q.Limit > 0 {
		pageCount = int64(math.Ceil(float64(total) / float64(q.Limit)))
	}
	if pageCount == 0 {
		pageCount = 1
	}

	return PaginatedMeta{
		Page:       q.Page,
		PerPage:    q.Limit,
		PageCount:  pageCount,
		TotalCount: total,
	}, nil
}

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
		field := transformNamespace(fe.Namespace(), request)
		res.Errors[field] = libvalidator.TranslateError(fe)
	}

	return res
}

func NewValidationErrorFastJsonValue(ar *fastjson.Arena, err error, request interface{}) *fastjson.Value {
	var validatorErrs validator.ValidationErrors
	if errors.As(err, &validatorErrs) {
		value := ar.NewObject()
		for _, fe := range validatorErrs {
			field := transformNamespace(fe.Namespace(), request)
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
	field  string
	errMsg string
}

func NewValidationError(field, errMsg string) ValidationError {
	return ValidationError{field: field, errMsg: errMsg}
}

func (v ValidationError) Error() string {
	return v.errMsg
}

func (v ValidationError) ValidationErrorResponse() ValidationErrorResponse {
	return ValidationErrorResponse{
		Errors: map[string]string{v.field: v.Error()},
	}
}

type AlarmStep struct {
	Type         string             `bson:"_t" json:"_t"`
	Timestamp    *datetime.CpsTime  `bson:"t" json:"t" swaggertype:"integer"`
	Author       string             `bson:"a" json:"a"`
	UserID       string             `bson:"user_id,omitempty" json:"user_id"`
	Message      string             `bson:"m" json:"m"`
	Value        types.CpsNumber    `bson:"val" json:"val"`
	Initiator    string             `bson:"initiator" json:"initiator"`
	Execution    string             `bson:"exec,omitempty" json:"-"`
	StateCounter *types.CropCounter `bson:"statecounter,omitempty" json:"statecounter,omitempty"`

	// Ticket related fields
	types.TicketInfo `bson:",inline"`

	InPbehaviorInterval bool `bson:"in_pbh,omitempty" json:"in_pbh,omitempty"`
}
