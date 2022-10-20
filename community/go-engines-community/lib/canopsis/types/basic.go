package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const (
	DurationUnitSecond = "s"
	DurationUnitMinute = "m"
	DurationUnitHour   = "h"
	DurationUnitDay    = "d"
	DurationUnitWeek   = "w"
	DurationUnitMonth  = "M"
	DurationUnitYear   = "y"
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
func (t CpsNumber) CpsTimestamp() CpsTime {
	return NewCpsTime(int64(t))
}

// CpsTime allows conversion from time.Time to time.Time.Unix()
type CpsTime struct {
	time.Time
}

// NewCpsTime create a CpsTime from a timestamp
func NewCpsTime(timestamp ...int64) CpsTime {
	if len(timestamp) == 0 {
		return CpsTime{Time: time.Now()}
	}

	if len(timestamp) > 1 {
		panic(fmt.Errorf("too much arguments, expected one: %+v", timestamp))
	}

	return CpsTime{time.Unix(timestamp[0], 0)}
}

// MarshalJSON converts from CpsTime to timestamp as bytes
func (t CpsTime) MarshalJSON() ([]byte, error) {
	ts := t.Time.Unix()
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
	return bson.MarshalValue(t.Time.Unix())
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

// Format easilly format a CpsTime as a string
func (t CpsTime) Format() string {
	return t.Time.Format(time.RFC3339Nano)
}

func (t CpsTime) Before(u CpsTime) bool {
	return t.Time.Before(u.Time)
}

func (t CpsTime) After(u CpsTime) bool {
	return t.Time.After(u.Time)
}

func (t CpsTime) In(loc *time.Location) CpsTime {
	return CpsTime{Time: t.Time.In(loc)}
}

func (t CpsTime) EqualDay(u CpsTime) bool {
	dateFormat := "2006-01-02"
	return t.Time.In(time.UTC).Format(dateFormat) == u.Time.In(time.UTC).Format(dateFormat)
}

type MicroTime struct {
	time.Time
}

func NewMicroTime() MicroTime {
	return MicroTime{Time: time.Now()}
}

func (t MicroTime) MarshalJSON() ([]byte, error) {
	unixMicro := t.Time.UnixMicro()
	if unixMicro <= 0 {
		return json.Marshal(nil)
	}

	return json.Marshal(unixMicro)
}

func (t *MicroTime) UnmarshalJSON(b []byte) error {
	var unixMicro int64
	err := json.Unmarshal(b, &unixMicro)
	if err != nil {
		return err
	}
	if unixMicro <= 0 {
		*t = MicroTime{}
		return nil
	}

	*t = MicroTime{Time: time.UnixMicro(unixMicro)}
	return nil
}

// DurationWithUnit represent duration with user-preferred units
type DurationWithUnit struct {
	Value int64  `bson:"value" json:"value" binding:"required,min=1"`
	Unit  string `bson:"unit" json:"unit" binding:"required,oneof=s m h d w M y"`
}

func NewDurationWithUnit(value int64, unit string) DurationWithUnit {
	return DurationWithUnit{
		Value: value,
		Unit:  unit,
	}
}

func (d DurationWithUnit) AddTo(t CpsTime) CpsTime {
	var r time.Time

	switch d.Unit {
	case DurationUnitSecond:
		r = t.Add(time.Duration(d.Value) * time.Second)
	case DurationUnitMinute:
		r = t.Add(time.Duration(d.Value) * time.Minute)
	case DurationUnitHour:
		r = t.Add(time.Duration(d.Value) * time.Hour)
	case DurationUnitDay:
		r = t.AddDate(0, 0, int(d.Value))
	case DurationUnitWeek:
		r = t.AddDate(0, 0, 7*int(d.Value))
	case DurationUnitMonth:
		r = t.AddDate(0, int(d.Value), 0)
	case DurationUnitYear:
		r = t.AddDate(int(d.Value), 0, 0)
	default:
		r = t.Add(time.Duration(d.Value) * time.Second)
	}

	return CpsTime{Time: r}
}

func (d DurationWithUnit) SubFrom(t CpsTime) CpsTime {
	var r time.Time

	switch d.Unit {
	case DurationUnitSecond:
		r = t.Add(-time.Duration(d.Value) * time.Second)
	case DurationUnitMinute:
		r = t.Add(-time.Duration(d.Value) * time.Minute)
	case DurationUnitHour:
		r = t.Add(-time.Duration(d.Value) * time.Hour)
	case DurationUnitDay:
		r = t.AddDate(0, 0, -int(d.Value))
	case DurationUnitWeek:
		r = t.AddDate(0, 0, -7*int(d.Value))
	case DurationUnitMonth:
		r = t.AddDate(0, -int(d.Value), 0)
	case DurationUnitYear:
		r = t.AddDate(-int(d.Value), 0, 0)
	default:
		r = t.Add(-time.Duration(d.Value) * time.Second)
	}

	return CpsTime{Time: r}
}

func (d DurationWithUnit) To(unit string) (DurationWithUnit, error) {
	newDuration := DurationWithUnit{
		Value: d.Value,
		Unit:  unit,
	}

	if d.Unit == unit || d.Value == 0 {
		return newDuration, nil
	}

	in := int64(0)

	switch d.Unit {
	case DurationUnitMinute:
		if unit == DurationUnitSecond {
			in = 60
		}
	case DurationUnitHour:
		switch unit {
		case DurationUnitMinute:
			in = 60
		case DurationUnitSecond:
			in = 60 * 60
		}
	case DurationUnitDay:
		switch unit {
		case DurationUnitHour:
			in = 24
		case DurationUnitMinute:
			in = 24 * 60
		case DurationUnitSecond:
			in = 24 * 60 * 60
		}
	case DurationUnitWeek:
		switch unit {
		case DurationUnitDay:
			in = 7
		case DurationUnitHour:
			in = 7 * 24
		case DurationUnitMinute:
			in = 7 * 24 * 60
		case DurationUnitSecond:
			in = 7 * 24 * 60 * 60
		}
	}

	if in > 0 {
		newDuration.Value *= in
		return newDuration, nil
	}

	return DurationWithUnit{}, fmt.Errorf("cannot transform unit %q to %q", d.Unit, unit)
}

func (d DurationWithUnit) String() string {
	return fmt.Sprintf("%d%s", d.Value, d.Unit)
}

func (d DurationWithUnit) IsZero() bool {
	return d == DurationWithUnit{}
}

func ParseDurationWithUnit(str string) (DurationWithUnit, error) {
	d := DurationWithUnit{}
	if str == "" {
		return d, fmt.Errorf("invalid duration %q", str)
	}

	r := regexp.MustCompile(`^(-?)(?P<val>\d+)(?P<t>[smhdwMy])$`)
	res := r.FindStringSubmatch(str)
	if len(res) == 0 {
		return d, fmt.Errorf("invalid duration %q", str)
	}

	val, err := strconv.Atoi(res[2])
	if err != nil {
		return d, fmt.Errorf("invalid duration %q: %w", str, err)
	}

	d.Value = int64(val)
	d.Unit = res[3]

	if res[1] == "-" {
		d.Value = -d.Value
	}

	return d, nil
}

type DurationWithEnabled struct {
	DurationWithUnit `bson:",inline"`
	Enabled          *bool `bson:"enabled" json:"enabled" binding:"required"`
}

func NewDurationWithEnabled(value int64, unit string, enabled *bool) DurationWithEnabled {
	return DurationWithEnabled{
		DurationWithUnit: DurationWithUnit{
			Value: value,
			Unit:  unit,
		},
		Enabled: enabled,
	}
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
	case float64:
		return int64(math.Round(value.(float64))), true
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
