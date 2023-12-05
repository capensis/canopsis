package datetime

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

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
	if ts <= 0 {
		return []byte("null"), nil
	}

	stamp := strconv.FormatInt(ts, 10)

	return []byte(stamp), nil
}

// UnmarshalJSON converts from string to CpsTime
func (t *CpsTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
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
	case bson.TypeDouble:
		double, _, ok := bsoncore.ReadDouble(b)
		if !ok {
			return errors.New("invalid value, expected double")
		}

		*t = NewCpsTime(int64(double))
	case bson.TypeInt64:
		i64, _, ok := bsoncore.ReadInt64(b)
		if !ok {
			return errors.New("invalid value, expected int64")
		}

		*t = NewCpsTime(i64)
	case bson.TypeInt32:
		i32, _, ok := bsoncore.ReadInt32(b)
		if !ok {
			return errors.New("invalid value, expected int64")
		}

		*t = NewCpsTime(int64(i32))
	case bson.TypeNull:
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

func (t MicroTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	unixMicro := t.Time.UnixMicro()
	if unixMicro <= 0 {
		return bson.TypeNull, nil, nil
	}

	return bson.MarshalValue(unixMicro)
}

func (t *MicroTime) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bson.TypeDouble:
		double, _, ok := bsoncore.ReadDouble(b)
		if !ok {
			return errors.New("invalid value, expected double")
		}

		*t = MicroTime{Time: time.UnixMicro(int64(double))}
	case bson.TypeInt64:
		i64, _, ok := bsoncore.ReadInt64(b)
		if !ok {
			return errors.New("invalid value, expected int64")
		}

		*t = MicroTime{Time: time.UnixMicro(i64)}
	case bson.TypeInt32:
		i32, _, ok := bsoncore.ReadInt32(b)
		if !ok {
			return errors.New("invalid value, expected int32")
		}

		*t = MicroTime{Time: time.UnixMicro(int64(i32))}
	case bson.TypeNull:
		//do nothing
	default:
		return fmt.Errorf("unexpected type %v", valueType)
	}

	return nil
}
