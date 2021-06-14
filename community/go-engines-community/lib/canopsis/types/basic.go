package types

import (
	"bytes"
	crand "crypto/rand"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	mdbson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
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

// GetBSON implements bson.Getter interface
func (t CpsNumber) GetBSON() (interface{}, error) {
	return int64(t), nil
}

// SetBSON implements bson.Setter interface
func (t *CpsNumber) SetBSON(raw bson.Raw) error {
	var num float64

	if err := raw.Unmarshal(&num); err != nil {
		return err
	}

	*t = CpsNumber(num)
	return nil
}

// Float64 ...
func (t CpsNumber) Float64() float64 {
	return float64(t)
}

// CpsTimestamp convert a number to a timestamp
func (t CpsNumber) CpsTimestamp() CpsTime {
	return NewCpsTime(int64(t))
}

// CpsTime allows conversion from time.Time to time.Time.Unix()
type CpsTime struct {
	time.Time
}

// NewCpsTime create a CpsTime from a timestamp
func NewCpsTime(timestamp int64) CpsTime {
	return CpsTime{time.Unix(timestamp, 0)}
}

// MarshalJSON converts from CpsTime to timestamp as bytes
func (t CpsTime) MarshalJSON() ([]byte, error) {
	ts := time.Time(t.Time).Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

// UnmarshalJSON converts from string to CpsTime
func (t *CpsTime) UnmarshalJSON(b []byte) error {
	if nl := bytes.TrimPrefix(b, []byte("{\"$numberLong\":\"")); len(nl) < len(b) {
		// json value can be recorded as "tstop":{"$numberLong":"4733481300"}
		nl = bytes.TrimSuffix(nl, []byte("\"}"))
		b = nl[:]
	}
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		tt, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			return err
		}
		ts = int(tt)
	}

	*t = NewCpsTime(int64(ts))
	return nil
}

// MarshalBSONValue converts from CpsTime to timestamp as bytes
func (t CpsTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return mdbson.MarshalValue(t.Time.Unix())
}

// UnmarshalBSONValue converts from timestamp as bytes to CpsTime
func (t *CpsTime) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Double:
		double, _, ok := bsoncore.ReadDouble(b)
		if !ok {
			return errors.New("invalid value, expected double")
		}

		*t = NewCpsTime(int64(double))
	case bsontype.Int64:
		i64, _, ok := bsoncore.ReadInt64(b)
		if !ok {
			return errors.New("invalid value, expected int64")
		}

		*t = NewCpsTime(i64)
	case bsontype.Int32:
		i32, _, ok := bsoncore.ReadInt32(b)
		if !ok {
			return errors.New("invalid value, expected int64")
		}

		*t = NewCpsTime(int64(i32))
	case bsontype.Null:
		//do nothing
	default:
		return fmt.Errorf("unexpected type %v", valueType)
	}

	return nil
}

// GetBSON implements the bson.Getter interface
func (t CpsTime) GetBSON() (interface{}, error) {
	ts := time.Time(t.Time).Unix()
	return ts, nil
}

// SetBSON implements the bson.Setter interface
func (t *CpsTime) SetBSON(raw bson.Raw) error {
	var ts int64

	if err := raw.Unmarshal(&ts); err != nil {
		return err
	}

	*t = NewCpsTime(int64(ts))
	return nil
}

// Format easilly format a CpsTime as a string
func (t CpsTime) Format() string {
	return t.Time.Format(time.RFC3339Nano)
}

// CpsDuration allow conversions from/to time.Duration to/from string
type CpsDuration time.Duration

// MarshalJSON converts a CpsDuration to string
func (t CpsDuration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Duration(t).String() + "\""), nil
}

// UnmarshalJSON converts a string to CpsDuration
func (t *CpsDuration) UnmarshalJSON(b []byte) error {
	sduration := string(b)

	if len(sduration) < 3 {
		return errors.New("Bad duration: string length below 3 chars")
	}

	if sduration[0] != '"' || sduration[len(sduration)-1] != '"' {
		return errors.New("Bad duration: not a string")
	}

	parsed, err := time.ParseDuration(sduration[1 : len(sduration)-1])
	if err != nil {
		return err
	}
	*t = CpsDuration(parsed)
	return nil
}

// GetBSON implements the bson.Getter interface
func (t CpsDuration) GetBSON() (interface{}, error) {
	return time.Duration(t).String(), nil
}

// SetBSON implements the bson.Setter interface
func (t *CpsDuration) SetBSON(raw bson.Raw) error {
	var d string

	if err := raw.Unmarshal(&d); err != nil {
		return err
	}

	parsed, err := time.ParseDuration(d)
	if err != nil {
		return err
	}
	*t = CpsDuration(parsed)
	return nil
}

// MarshalBSONValue converts from CpsDuration to bytes
func (t CpsDuration) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return mdbson.MarshalValue(t.Duration().String())
}

// UnmarshalBSONValue converts from bytes to CpsDuration
func (t *CpsDuration) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.String:
		str, _, ok := bsoncore.ReadString(b)
		if !ok {
			return errors.New("invalid value, expected string")
		}

		d, err := time.ParseDuration(str)
		if err != nil {
			return err
		}

		*t = CpsDuration(d)
	default:
		return fmt.Errorf("unexpected type %v", valueType)
	}

	return nil
}

// Duration return the CpsDuration casted to time.Duration
func (t CpsDuration) Duration() time.Duration {
	return time.Duration(t)
}

// CpsShortDuration allow conversions from/to time.Duration to/from string
type CpsShortDuration int64

// MarshalJSON converts a CpsDuration to string
func (t CpsShortDuration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.String() + "\""), nil
}

// UnmarshalJSON converts a string to CpsShortDuration
func (t *CpsShortDuration) UnmarshalJSON(b []byte) error {
	str := string(b)

	if len(str) < 3 {
		return errors.New("bad duration: string length below 3 chars")
	}

	if str[0] != '"' || str[len(str)-1] != '"' {
		return errors.New("bad duration: not a string")
	}

	parsed, err := ParseCpsShortDuration(str[1 : len(str)-1])
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}

// MarshalBSONValue converts from CpsShortDuration to bytes
func (t CpsShortDuration) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return mdbson.MarshalValue(t.String())
}

// UnmarshalBSONValue converts from bytes to CpsShortDuration
func (t *CpsShortDuration) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.String:
		str, _, ok := bsoncore.ReadString(b)
		if !ok {
			return errors.New("invalid value, expected string")
		}

		d, err := ParseCpsShortDuration(str)
		if err != nil {
			return err
		}

		*t = d
	default:
		return fmt.Errorf("unexpected type %v", valueType)
	}

	return nil
}

// Duration return the CpsShortDuration casted to time.Duration
func (t CpsShortDuration) Duration() time.Duration {
	return time.Duration(t)
}

func ParseCpsShortDuration(str string) (CpsShortDuration, error) {
	if str == "" {
		return 0, nil
	}

	r := regexp.MustCompile(`(?P<val>\d+)(?P<t>[hms])`)
	res := r.FindStringSubmatch(str)
	if len(res) == 0 {
		return 0, nil
	}

	val, err := strconv.Atoi(res[1])
	if err != nil {
		return 0, err
	}

	switch res[2] {
	case "h":
		return CpsShortDuration(int64(val) * int64(time.Hour)), nil
	case "m":
		return CpsShortDuration(int64(val) * int64(time.Minute)), nil
	case "s":
		return CpsShortDuration(int64(val) * int64(time.Second)), nil
	}

	return 0, nil
}

func (t CpsShortDuration) String() string {
	d := t.Duration()
	if h := d / time.Hour; h > 0 {
		return fmt.Sprintf("%dh", h)
	}
	if m := d / time.Minute; m > 0 {
		return fmt.Sprintf("%dm", m)
	}
	if s := d / time.Second; s > 0 {
		return fmt.Sprintf("%ds", s)
	}

	return ""
}

// DurationWithUnit represent duration with user-preferred units
type DurationWithUnit struct {
	Seconds int64  `bson:"seconds" json:"seconds" binding:"required"`
	Unit    string `bson:"unit" json:"unit" binding:"required"`
}

func listOfInterfaceToString(v []interface{}) (string, error) {
	values := make([]string, len(v))
	for i, vv := range v {
		sval, err := InterfaceToString(vv)
		if err != nil {
			return "", fmt.Errorf("list of: %v", err)
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
			return []string{}, fmt.Errorf("list of: %v", err)
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
		return strconv.FormatInt(int64(vt), 10), nil
	case int64:
		return strconv.FormatInt(vt, 10), nil
	case uint:
		return strconv.FormatUint(uint64(vt), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(vt), 10), nil
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
	case int64:
		return typedValue, true
	case uint64:
		return int64(typedValue), true
	case int:
		return int64(typedValue), true
	case uint:
		return int64(typedValue), true
	case CpsNumber:
		return int64(typedValue), true
	case time.Time:
		return typedValue.Unix(), true
	case CpsTime:
		return typedValue.Unix(), true
	default:
		return 0, false
	}
}

const NamingCharacterSet = "abcdefghijklmnopqrstuvwxyz1234567890"
const DefaultNameSchemeTemplate = "{{ rand_string 2 }}-{{ rand_string 2 }}-{{ rand_string 2 }}"
const (
	RandomWithoutPrefixNamingScheme = iota
	RandomWithTimePrefixNamingScheme
	RandomWithResourcePrefix
	RandomWithConnectorPrefix
)

var NumberOfCharacter = int64(len(NamingCharacterSet))

var genDisplayNameFuns = template.FuncMap{
	"rand_string": func(n int) string {
		return RandString(n)
	},
}

var displayTpl *template.Template
var oldScheme = ""

func getDisplayNameTemplate(displayNameScheme string) *template.Template {
	if displayTpl == nil || oldScheme != displayNameScheme {
		if displayNameScheme != "" {
			tml, err := template.New("displayname_gen_scheme").Funcs(genDisplayNameFuns).
				Parse(displayNameScheme)
			if err != nil {
				log.Println("Can not parse display name scheme: ", err)
				displayTpl, _ = template.New("default_displayname_gen_scheme").Funcs(genDisplayNameFuns).
					Parse(DefaultNameSchemeTemplate)
			} else {
				oldScheme = displayNameScheme
				displayTpl = tml
			}
		} else {
			log.Println("Displayname scheme is not set")
			displayTpl, _ = template.New("default_displayname_gen_scheme").Funcs(genDisplayNameFuns).
				Parse(DefaultNameSchemeTemplate)
		}

	}
	return displayTpl
}

// RandString generate a random string
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		index, _ := crand.Int(crand.Reader, big.NewInt(NumberOfCharacter))
		b[i] = NamingCharacterSet[index.Int64()]
	}
	return string(b)
}

func GenDisplayName(displayNameScheme string) string {
	tml := getDisplayNameTemplate(displayNameScheme)
	var b bytes.Buffer
	err := tml.Execute(io.Writer(&b), nil)
	if err != nil {
		log.Println("Gen display name had error: ", err)
		return defaultDisplayNameFunc()
	}
	return strings.ToUpper(b.String())
}

// GenDisplayName generate a random uppercased display name
func defaultDisplayNameFunc() string {
	name := RandString(2) + "-" + RandString(2) + "-" + RandString(2)
	return strings.ToUpper(name)
}
