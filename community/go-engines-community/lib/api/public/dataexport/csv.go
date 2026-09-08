package dataexport

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/link"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var csvDelimiters = map[string]rune{"comma": ',', "semicolon": ';', "tab": '\t', "space": ' '}

type CsvOptions struct {
	Options
	// Delimiter is the csv delimiter, comma if empty.
	Delimiter string
}

// ToCsv streams every document of cursor into a csv file, see export.
func ToCsv(ctx context.Context, cursor libmongo.Cursor, opts CsvOptions) error {
	return export(ctx, cursor, opts.Options, "csv", func(f io.Writer) (recordWriter, error) {
		w := csv.NewWriter(f)
		if opts.Delimiter != "" {
			var ok bool
			w.Comma, ok = csvDelimiters[opts.Delimiter]
			if !ok {
				return nil, fmt.Errorf("invalid delimiter: %s", opts.Delimiter)
			}
		}

		err := writeCsvHeader(w, opts.Fields)
		if err != nil {
			return nil, err
		}

		return &csvWriter{
			w:      w,
			opts:   opts.Options,
			record: make([]string, len(opts.Fields)),
		}, nil
	})
}

func writeCsvHeader(w *csv.Writer, fields []Field) error {
	if len(fields) == 0 {
		return nil
	}

	header := make([]string, len(fields))
	for i, field := range fields {
		header[i] = field.header()
	}

	return w.Write(header)
}

type csvWriter struct {
	w    *csv.Writer
	opts Options
	// record is reused between documents: every field is always set.
	record []string
}

func (w *csvWriter) Write(ctx context.Context, doc bson.M) error {
	for i, field := range w.opts.Fields {
		v, ok, err := w.opts.value(ctx, doc, field)
		if err != nil {
			return err
		}

		if !ok {
			w.record[i] = ""

			continue
		}

		w.record[i], err = csvValue(v)
		if err != nil {
			return err
		}
	}

	return w.w.Write(w.record)
}

func (w *csvWriter) Close() error {
	w.w.Flush()

	return w.w.Error()
}

// csvValue formats a value as a csv cell.
func csvValue(v any) (string, error) {
	if v == nil {
		return "", nil
	}

	switch rv := v.(type) {
	case bson.M,
		bson.A,
		[]link.Link:
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}

		return string(b), nil
	case float32:
		return strconv.FormatFloat(float64(rv), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(rv, 'f', -1, 64), nil
	default:
		return fmt.Sprintf("%v", rv), nil
	}
}
