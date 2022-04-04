package common

import (
	"errors"
	"github.com/valyala/fastjson"
	"math"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	libvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/validator"
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"crecord_name" json:"name"`
}

type Role struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"crecord_name" json:"name"`
}

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
	if !q.Paginate {
		q.Limit = d.GetTotal()
	}

	pageCount := int64(math.Ceil(float64(d.GetTotal()) / float64(q.Limit)))
	if pageCount == 0 {
		pageCount = 1
	}

	if q.Page > pageCount {
		return PaginatedListResponse{}, errors.New("page is out of range")
	}

	data := d.GetData()
	if data == nil {
		data = []interface{}{}
	}

	return PaginatedListResponse{
		Data: data,
		Meta: PaginatedMeta{
			Page:       q.Page,
			PerPage:    q.Limit,
			PageCount:  pageCount,
			TotalCount: d.GetTotal(),
		},
	}, nil
}

// ErrorResponse is base failed response.
type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Error: err.Error()}
}

var NotFoundResponse = ErrorResponse{Error: "Not found"}
var MethodNotAllowedResponse = ErrorResponse{Error: http.StatusText(http.StatusMethodNotAllowed)}
var UnauthorizedResponse = ErrorResponse{Error: http.StatusText(http.StatusUnauthorized)}
var InternalServerErrorResponse = ErrorResponse{Error: "Internal server error"}
var ForbiddenResponse = ErrorResponse{Error: http.StatusText(http.StatusForbidden)}
var ErrTimeoutResponse = ErrorResponse{Error: "Request timeout reached"}

// ValidationErrorResponse is response for failed validation.
type ValidationErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

// NewValidationErrorResponse creates response by validation errors.
func NewValidationErrorResponse(err error, request interface{}) interface{} {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		var res ValidationErrorResponse
		res.Errors = make(map[string]string)
		for _, fe := range errs {
			field := transformNamespace(fe.Namespace(), request)
			res.Errors[field] = libvalidator.TranslateError(fe)
		}

		return res
	}

	return ErrorResponse{Error: "request has invalid structure"}
}

func NewValidationErrorFastJsonValue(ar *fastjson.Arena, err error, request interface{}) *fastjson.Value {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		value := ar.NewObject()
		for _, fe := range errs {
			field := transformNamespace(fe.Namespace(), request)
			value.Set(field, ar.NewString(libvalidator.TranslateError(fe)))
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
