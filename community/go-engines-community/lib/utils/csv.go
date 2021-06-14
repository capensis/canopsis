package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ExportCsv adds input array of structs to csv writer.
// Only defined fields of struct are added to result.
func ExportCsv(w *csv.Writer, data interface{}, fields []string) error {
	if len(fields) == 0 {
		return nil
	}

	records, err := convertToCsvData(data, fields)
	if err != nil {
		return err
	}

	w.WriteAll(records)

	if err = w.Error(); err != nil {
		return err
	}

	return nil
}

// convertToCsvData converts array of structs to format which is used by csv writer.
func convertToCsvData(data interface{}, fields []string) ([][]string, error) {
	m, err := convertToArrayOfMaps(data)
	if err != nil {
		return nil, err
	}

	records := make([][]string, len(m))

	for i := range m {
		records[i] = make([]string, len(fields))

		for fi, f := range fields {
			var str string
			v := getNestedMapVal(m[i], strings.Split(f, "."))

			if v != nil {
				k := reflect.TypeOf(v).Kind()
				switch k {
				case reflect.Array, reflect.Slice, reflect.Interface, reflect.Map, reflect.Struct, reflect.Ptr:
					b, err := json.Marshal(v)
					if err != nil {
						panic(err)
					}
					str = string(b)
				case reflect.Float32:
					str = fmt.Sprintf("%s", strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64))
				case reflect.Float64:
					str = fmt.Sprintf("%s", strconv.FormatFloat(v.(float64), 'f', -1, 64))
				default:
					str = fmt.Sprintf("%v", v)
				}
			}

			records[i][fi] = str
		}
	}

	return records, nil
}

// convertToArrayOfMaps converts array of structs to array of maps.
func convertToArrayOfMaps(data interface{}) ([]map[string]interface{}, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var m []map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// getNestedMapVal returns item value of nested map by keys.
func getNestedMapVal(m map[string]interface{}, keys []string) interface{} {
	if len(keys) == 0 {
		return nil
	}

	k := keys[0]
	if v, ok := m[k]; ok {
		if len(keys) > 1 {
			if mv, ok := v.(map[string]interface{}); ok {
				return getNestedMapVal(mv, keys[1:])
			}

			return nil
		}

		return v
	}

	return nil
}
