package common

import (
	"errors"
	"math"
	"reflect"
	"regexp"
	"strings"

	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	libvalidator "git.canopsis.net/canopsis/go-engines/lib/validator"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		libvalidator.RegisterTranslations(v)
		err := v.RegisterValidation("notblank", validators.NotBlank)
		if err != nil {
			panic(err)
		}
		err = v.RegisterValidation("oneoforempty", ValidateOneOfOrEmpty)
		if err != nil {
			panic(err)
		}
		v.RegisterCustomTypeFunc(ValidateCpsTimeType, types.CpsTime{})
	}
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

func NewPaginatedResponse(q pagination.Query, d PaginatedData) (interface{}, error) {
	if q.Paginate == false {
		q.Limit = d.GetTotal()
	}

	pageCount := int64(math.Ceil(float64(d.GetTotal()) / float64(q.Limit)))
	if pageCount == 0 {
		pageCount = 1
	}

	if q.Page > pageCount {
		return nil, errors.New("page is out of range")
	}

	return &PaginatedListResponse{
		Data: d.GetData(),
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
var UnauthorizedResponse = ErrorResponse{Error: "Unauthorized"}
var ForbiddenResponse = ErrorResponse{Error: "Forbidden"}

// ValidationErrorResponse is response for failed validation.
type ValidationErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

// NewValidationErrorResponse creates response by validation errors.
func NewValidationErrorResponse(err error, request interface{}) interface{} {
	if errs, ok := err.(validator.ValidationErrors); ok {
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
	for i := range path {
		k := val.Kind()
		switch k {
		case reflect.Interface, reflect.Ptr:
			val = val.Elem()
			fallthrough
		case reflect.Struct:
			if f, ok := val.Type().FieldByName(path[i]); ok {
				val = val.FieldByName(path[i])
				tag := strings.Split(f.Tag.Get("json"), ",")
				if len(tag) > 0 {
					path[i] = tag[0]
				}
			}
		case reflect.Slice, reflect.Array:
			val = val.Index(0)
		case reflect.Map:
			val = val.MapIndex(val.MapKeys()[0])
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
