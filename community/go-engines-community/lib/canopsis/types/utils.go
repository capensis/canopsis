package types

import (
	"errors"
	"math"

	"go.mongodb.org/mongo-driver/bson"
)

var ErrInvalidInfoType = errors.New("info value should be int, string, bool or array of strings")

func IsInfoValueValid(value interface{}) bool {
	switch val := value.(type) {
	case int, int32, int64, string, bool, []string:
	// must support this special case for floats, because api saves json numbers as double in mongo for interface{} fields
	case float32:
		return float64(val) == math.Trunc(float64(val))
	case float64:
		return val == math.Trunc(val)
	// bson unmarshalling
	case bson.A:
		for _, v := range val {
			_, ok := v.(string)
			if !ok {
				return false
			}
		}
	case []interface{}:
		for _, v := range val {
			_, ok := v.(string)
			if !ok {
				return false
			}
		}
	default:
		return false
	}

	return true
}
