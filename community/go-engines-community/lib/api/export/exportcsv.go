package export

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const pageLimit = 5000

// ExportCsv fetches data by page and saves it in csv file.
func ExportCsv(
	ctx context.Context,
	exportFields Fields,
	separator rune,
	dataFetcher DataFetcher,
) (resFileName string, resErr error) {
	if len(exportFields) == 0 {
		return "", fmt.Errorf("exportFields is empty")
	}

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

	fieldsLabels := exportFields.Labels()
	fields := exportFields.Fields()

	err = w.WriteAll([][]string{fieldsLabels})
	if err != nil {
		return "", err
	}

	var limit int64 = pageLimit
	data, totalCount, err := dataFetcher(ctx, 1, limit)
	if err != nil {
		return "", err
	}

	if totalCount > 0 {
		err = writeToCsv(w, data, fields)
		if err != nil {
			return "", err
		}

		pageCount := int64(math.Ceil(float64(totalCount) / float64(limit)))
		if pageCount > 1 {
			pageCh := make(chan int64)
			dataCh := runExportWorkers(ctx, pageCh, dataFetcher, limit)

			go func() {
				defer close(pageCh)
				for p := int64(2); p <= pageCount; p++ {
					pageCh <- p
				}
			}()

			var err error
			dataByPage := make(map[int64][]map[string]string)
			page := int64(2)
			for res := range dataCh {
				if res.Err != nil {
					err = res.Err
					continue
				}

				if res.Page == page {
					err = writeToCsv(w, res.Data, fields)
					if err != nil {
						return "", err
					}

					for p := page + 1; p <= pageCount; p++ {
						if d, ok := dataByPage[p]; ok {
							err = writeToCsv(w, d, fields)
							if err != nil {
								return "", err
							}

							delete(dataByPage, p)
						} else {
							page = p
							break
						}
					}
				} else {
					dataByPage[res.Page] = res.Data
				}
			}

			if err != nil {
				return "", err
			}
		}
	}

	return file.Name(), nil
}

func ExportCsvByCursor(
	ctx context.Context,
	exportFields Fields,
	separator rune,
	dataCursor DataCursor,
) (resFileName string, resErr error) {
	defer dataCursor.Close()
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

	var data []map[string]string
	for dataCursor.Next(ctx) {
		var item map[string]string
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

type workerResult struct {
	Page int64
	Data []map[string]string
	Err  error
}

// runExportWorkers starts workers. Each worker fetches specific page and sends csv data to output channel.
func runExportWorkers(
	ctx context.Context,
	in <-chan int64,
	dataFetcher DataFetcher,
	limit int64,
) <-chan workerResult {
	out := make(chan workerResult)

	go func() {
		defer close(out)

		wg := sync.WaitGroup{}

		for i := 0; i < 10; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for {
					select {
					case <-ctx.Done():
						return
					case page, ok := <-in:
						if !ok {
							return
						}

						data, _, err := dataFetcher(ctx, page, limit)
						if err != nil {
							out <- workerResult{
								Page: page,
								Err:  err,
							}
							return
						}

						out <- workerResult{
							Page: page,
							Data: data,
						}
					}
				}
			}()
		}

		wg.Wait()
	}()

	return out
}

// ConvertToMap converts array of structs to array of maps with flatten keys.
func ConvertToMap(
	data interface{},
	fields []string,
	timeFormat string,
	location *time.Location,
) ([]map[string]string, error) {
	val := reflect.ValueOf(data)
	kind := val.Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		return nil, fmt.Errorf("unexpected argument: %v, expected array or slice", data)
	}

	splittedFields := make(map[string][]string, len(fields))
	for _, f := range fields {
		splittedFields[f] = strings.Split(f, ".")
	}

	result := make([]map[string]string, val.Len())
	var err error

	for i := 0; i < val.Len(); i++ {
		result[i] = make(map[string]string)

		for f, splittedField := range splittedFields {
			result[i][f], err = getNestedMapValAsString(val.Index(i), splittedField, timeFormat, location)
			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

// getNestedMapValAsString returns item value of nested map by keys.
func getNestedMapValAsString(
	val reflect.Value,
	keys []string,
	timeFormat string,
	location *time.Location,
) (string, error) {
	if len(keys) == 0 {
		return "", errors.New("invalid keys argument")
	}

	if !val.IsValid() || val.IsZero() {
		return "", nil
	}

	k := keys[0]

	switch val.Kind() {
	case reflect.Interface, reflect.Ptr:
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		val = val.MapIndex(reflect.ValueOf(k))
		if !val.IsValid() || val.IsZero() {
			return "", nil
		}

		if len(keys) > 1 {
			return getNestedMapValAsString(val, keys[1:], timeFormat, location)
		}

		return interfaceToString(val.Interface(), timeFormat, location)
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			f := val.Type().Field(i)
			tag := strings.Split(f.Tag.Get("json"), ",")
			if len(tag) > 0 {
				if tag[0] == k {
					val = val.Field(i)
					if len(keys) > 1 {
						return getNestedMapValAsString(val, keys[1:], timeFormat, location)
					}

					return interfaceToString(val.Interface(), timeFormat, location)
				}
			}
		}

		return "", nil
	}

	return interfaceToString(val.Interface(), timeFormat, location)
}

// writeToCsv adds input array of structs to csv writer.
func writeToCsv(
	w *csv.Writer,
	data []map[string]string,
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
	data []map[string]string,
	fields []string,
) ([][]string, error) {
	records := make([][]string, len(data))

	for i := range data {
		records[i] = make([]string, len(fields))

		for fi, f := range fields {
			records[i][fi] = data[i][f]
		}
	}

	return records, nil
}

func interfaceToString(
	v interface{},
	timeFormat string,
	location *time.Location,
) (string, error) {
	if v == nil {
		return "", nil
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.Interface, reflect.Ptr:
		rVal := reflect.ValueOf(v)
		if !rVal.IsValid() || rVal.IsZero() {
			v = nil
		} else {
			v = rVal.Elem().Interface()
		}
	}

	if v == nil {
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
	case reflect.Struct:
		if timeFormat != "" {
			if t, ok := v.(time.Time); ok {
				if location != nil {
					t = t.In(location)
				}
				str = t.Format(timeFormat)
			} else if ct, ok := v.(types.CpsTime); ok {
				t := ct.Time
				if location != nil {
					t = t.In(location)
				}
				str = t.Format(timeFormat)
			}
		}

		if str == "" {
			b, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			str = string(b)
		}
	case reflect.Float32:
		str = strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64)
	case reflect.Float64:
		str = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	default:
		str = fmt.Sprintf("%v", v)
	}

	return str, nil
}
