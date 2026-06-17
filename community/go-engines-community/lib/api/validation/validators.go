package validation

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

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
