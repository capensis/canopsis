package types

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// CpsNumber is here for compatibility with old python engines.
// It will force an int64 from a float64.
type CpsNumber int64

// MarshalJSON implements json.Encoder interface
func (t CpsNumber) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint(t)), nil
}

// UnmarshalJSON implements json.Decoder interface
func (t *CpsNumber) UnmarshalJSON(b []byte) error {
	f, err := strconv.ParseFloat(string(b), 64)
	*t = CpsNumber(f)
	return err
}

// Float64 ...
func (t CpsNumber) Float64() float64 {
	return float64(t)
}

// CpsTimestamp convert a number to a timestamp
func (t CpsNumber) CpsTimestamp() datetime.CpsTime {
	return datetime.NewCpsTime(int64(t))
}

func listOfInterfaceToString(v []interface{}) (string, error) {
	values := make([]string, len(v))
	for i, vv := range v {
		sval, err := InterfaceToString(vv)
		if err != nil {
			return "", fmt.Errorf("list of: %w", err)
		}
		values[i] = sval
	}

	return strings.Join(values, ","), nil
}

func listOfInterfaceToStringSlice(v []interface{}) ([]string, error) {
	values := make([]string, len(v))
	for i, vv := range v {
		sval, err := InterfaceToString(vv)
		if err != nil {
			return []string{}, fmt.Errorf("list of: %w", err)
		}
		values[i] = sval
	}

	return values, nil
}

// InterfaceToString tries to convert v to it's string value.
// Supported types:
// * float64
// * string
// * bool
// * int|int64|uint|uint64
// * []interface{} : join elements with ","
// * nil: will return empty string
//
// Any other type will return empty string and an error, like lists
// or maps...
func InterfaceToString(v interface{}) (string, error) {
	switch vt := v.(type) {
	case bool:
		return strconv.FormatBool(vt), nil
	case float64:
		return strconv.FormatFloat(vt, 'f', -1, 64), nil
	case string:
		return vt, nil
	case int:
		return strconv.Itoa(vt), nil
	case int64:
		return strconv.FormatInt(vt, 10), nil
	case uint:
		return strconv.FormatUint(uint64(vt), 10), nil
	case uint64:
		return strconv.FormatUint(vt, 10), nil
	case []interface{}:
		return listOfInterfaceToString(vt)
	default:
		return "", fmt.Errorf("unsupported type: %T", v)
	}
}

func InterfaceToStringSlice(v interface{}) ([]string, error) {
	switch vt := v.(type) {
	case []interface{}:
		return listOfInterfaceToStringSlice(vt)
	case []string:
		return vt, nil
	default:
		return []string{}, fmt.Errorf("unsupported type: %T", v)
	}
}

// AsInteger tries to convert an interface{} into an int64, and returns its
// value and an integer indicating whether it succeeded or not.
//
// It works with int, uint, int64, uint64, CpsNumber and CpsTime (in this case,
// a unix timestamp is returned).
func AsInteger(value interface{}) (int64, bool) {
	switch typedValue := value.(type) {
	case float32:
		return int64(math.Round(float64(typedValue))), true
	case float64:
		return int64(math.Round(typedValue)), true
	case int:
		return int64(typedValue), true
	case int8:
		return int64(typedValue), true
	case int16:
		return int64(typedValue), true
	case int32:
		return int64(typedValue), true
	case int64:
		return typedValue, true
	case uint:
		return int64(typedValue), true
	case uint8:
		return int64(typedValue), true
	case uint16:
		return int64(typedValue), true
	case uint32:
		return int64(typedValue), true
	case uint64:
		return int64(typedValue), true
	case CpsNumber:
		return int64(typedValue), true
	case *CpsNumber:
		if typedValue == nil {
			return 0, false
		}
		return int64(*typedValue), true
	case time.Time:
		return typedValue.Unix(), true
	case *time.Time:
		if typedValue == nil {
			return 0, false
		}
		return typedValue.Unix(), true
	case datetime.CpsTime:
		return typedValue.Unix(), true
	case *datetime.CpsTime:
		if typedValue == nil {
			return 0, false
		}
		return typedValue.Unix(), true
	default:
		return 0, false
	}
}

func GenDisplayName(tpl *template.Template) string {
	var b bytes.Buffer
	err := tpl.Execute(io.Writer(&b), nil)
	if err != nil {
		log.Println("Gen display name had error: ", err)
		return defaultDisplayNameFunc()
	}
	return strings.ToUpper(b.String())
}

// GenDisplayName generate a random uppercased display name
func defaultDisplayNameFunc() string {
	name := utils.RandString(2) + "-" + utils.RandString(2) + "-" + utils.RandString(2)
	return strings.ToUpper(name)
}
