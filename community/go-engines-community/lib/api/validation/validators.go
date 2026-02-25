package validation

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/go-playground/validator/v10"
)

const (
	MaxIDLength    = 255
	InvalidIDChars = "/?.$"

	tableNameRegexString = `^[a-zA-Z_]\w+$`
	tableNameMaxLen      = 63
)

var timeFormats = map[string]string{
	"YYYY-MM-DDThh:mm:ss":        "2006-01-02T15:04:05",
	"YYYY-MM-DDThh:mm:ssZ":       "2006-01-02T15:04:05-0700",
	"DD MMM YYYY hh:mm:ss":       "02 Jan 2006 15:04",
	"DD MMM YYYY hh:mm:ss ZZ":    "02 Jan 2006 15:04 MST",
	"W, DD MMM YYYY hh:mm:ss ZZ": "Mon, 02 Jan 2006 15:04:05 MST",
}

var tableNameRegex = regexp.MustCompile(tableNameRegexString)

// ValidateCpsTimeType implements CustomTypeFunc and returns value to validate.
func ValidateCpsTimeType(field reflect.Value) any {
	if field.Type() == reflect.TypeFor[datetime.CpsTime]() {
		if t, ok := field.Interface().(datetime.CpsTime); ok {
			val := t.Time
			if val.IsZero() {
				return nil
			}

			return val
		}
	}

	return nil
}

func ValidateOneOfOrEmpty(fl validator.FieldLevel) bool {
	vals := strings.Split(fl.Param(), " ")
	field := fl.Field()

	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	default:
		return false
	}

	if v == "" {
		return true
	}

	for i := 0; i < len(vals); i++ {
		prefix := strings.TrimSuffix(vals[i], "*")
		if prefix != "" && prefix != vals[i] {
			if strings.HasPrefix(v, prefix) {
				return true
			}
		} else if vals[i] == v {
			return true
		}
	}

	return false
}

func ValidateColorOrEmpty(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v == "" {
		return true
	}
	return validator.New().Var(v, "iscolor") == nil
}

func ValidateID(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v == "" {
		return true
	}

	return !strings.ContainsAny(v, InvalidIDChars) && len(v) <= MaxIDLength
}

func ValidateTimeFormat(fl validator.FieldLevel) bool {
	v := fl.Field().String()

	return v == "" || timeFormats[v] != ""
}

func GetRealFormatTime(f string) string {
	if f == "" {
		return ""
	}
	return timeFormats[f]
}

func ValidateFilteredQuery(sl validator.StructLevel) {
	r := sl.Current().Interface().(pagination.FilteredQuery)
	// Validate sort
	if r.Sort != "" {
		sorts := []string{
			pagination.SortAsc,
			pagination.SortDesc,
		}

		found := false
		for _, sort := range sorts {
			if sort == r.Sort {
				found = true
			}
		}

		if !found {
			param := strings.Join(sorts, " ")
			sl.ReportError(r.Sort, "Sort", "sort", "oneof", param)
		}
	}
}

func ValidateTableName(fl validator.FieldLevel) bool {
	return IsTableName(fl.Field().String())
}

func ValidateTemplate(templateExecutor template.Executor) validator.Func {
	return func(fl validator.FieldLevel) bool {
		v := fl.Field().String()
		if v == "" {
			return true
		}

		return templateExecutor.Parse(v).Err == nil
	}
}

func IsTableName(s string) bool {
	return tableNameRegex.MatchString(s) && len(s) <= tableNameMaxLen
}

func ValidateInfoValue(fl validator.FieldLevel) bool {
	return types.IsInfoValueValid(fl.Field().Interface())
}
