package validation

import (
	"reflect"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"github.com/go-playground/validator/v10"
)

const DefaultTimeFormat = "2006-01-02T15:04:05"

var timeFormats = map[string]string{
	"YYYY-MM-DDThh:mm:ss":        DefaultTimeFormat,
	"YYYY-MM-DDThh:mm:ssZ":       "2006-01-02T15:04:05-0700",
	"DD MMM YYYY hh:mm:ss":       "02 Jan 2006 15:04",
	"DD MMM YYYY hh:mm:ss ZZ":    "02 Jan 2006 15:04 MST",
	"W, DD MMM YYYY hh:mm:ss ZZ": "Mon, 02 Jan 2006 15:04:05 MST",
}

func ValidateOneOfOrEmpty(fl validator.FieldLevel) bool {
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

	for val := range strings.SplitSeq(fl.Param(), " ") {
		prefix := strings.TrimSuffix(val, "*")
		if prefix != "" && prefix != val {
			if strings.HasPrefix(v, prefix) {
				return true
			}
		} else if val == v {
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
