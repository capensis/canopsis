package export

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
)

// ToCsv fetches data and saves it in csv file.
func ToCsv(
	ctx context.Context,
	exportFields Fields,
	separator rune,
	dataCursor DataCursor,
) (resFileName string, resErr error) {
	defer func() {
		err := dataCursor.Close(ctx)
		if err != nil && resErr == nil {
			resErr = err
		}
	}()
	file, err := os.CreateTemp("", "export.*.csv")
	if err != nil {
		return "", err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			return
		}

		if resErr != nil {
			err := os.Remove(file.Name())
			if err != nil {
				return
			}
		}
	}()

	w := csv.NewWriter(file)
	if separator != 0 {
		w.Comma = separator
	}

	var fields []string
	if len(exportFields) > 0 {
		fieldsLabels := exportFields.Labels()
		fields = exportFields.Fields()

		err = w.WriteAll([][]string{fieldsLabels})
		if err != nil {
			return "", err
		}
	}

	var data []map[string]any
	for dataCursor.Next(ctx) {
		var item map[string]any
		err := dataCursor.Scan(&item)
		if err != nil {
			return "", err
		}

		if len(fields) == 0 {
			for field := range item {
				fields = append(fields, field)
			}

			sort.Strings(fields)

			err = w.WriteAll([][]string{fields})
			if err != nil {
				return "", err
			}
		}

		data = append(data, item)

		if len(data) >= canopsis.DefaultBulkSize {
			err = writeToCsv(w, data, fields)
			if err != nil {
				return "", err
			}
			data = nil
		}
	}

	if len(data) > 0 {
		err = writeToCsv(w, data, fields)
		if err != nil {
			return "", err
		}
	}

	return file.Name(), nil
}

// writeToCsv adds input array of structs to csv writer.
func writeToCsv(
	w *csv.Writer,
	data []map[string]any,
	fields []string,
) error {
	if len(fields) == 0 {
		return nil
	}

	records, err := convertToCsvData(data, fields)
	if err != nil {
		return err
	}

	if err = w.WriteAll(records); err != nil {
		return err
	}

	return nil
}

// convertToCsvData converts array of structs to format which is used by csv writer.
func convertToCsvData(
	data []map[string]any,
	fields []string,
) ([][]string, error) {
	records := make([][]string, len(data))
	var err error

	for i := range data {
		records[i] = make([]string, len(fields))

		for fi, f := range fields {
			records[i][fi], err = toString(data[i][f])
			if err != nil {
				return nil, err
			}
		}
	}

	return records, nil
}

func toString(v any) (string, error) {
	if v == nil {
		return "", nil
	}

	rv := reflect.ValueOf(v)
	rv = unwrapPointer(rv)
	if !rv.IsValid() {
		return "", nil
	}

	str := ""
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		str = string(b)
	case reflect.Float32, reflect.Float64:
		str = strconv.FormatFloat(rv.Float(), 'f', -1, 64)
	default:
		str = fmt.Sprintf("%v", v)
	}

	return str, nil
}

func unwrapPointer(v reflect.Value) reflect.Value {
	for {
		switch v.Kind() {
		case reflect.Interface, reflect.Ptr:
			v = v.Elem()
		default:
			return v
		}
	}
}
