package utils

import "go.mongodb.org/mongo-driver/bson"

func IsStringSlice(arg any) ([]string, bool) {
	switch v := arg.(type) {
	case []string:
		return v, true
	case bson.A:
		return InterfaceSliceToStringSlice(v)
	case []any:
		return InterfaceSliceToStringSlice(v)
	default:
		return nil, false
	}
}

func InterfaceSliceToStringSlice(v []any) ([]string, bool) {
	res := make([]string, len(v))
	for i := range v {
		str, ok := v[i].(string)
		if !ok {
			return nil, false
		}

		res[i] = str
	}

	return res, true
}
