package dataexport

import (
	"context"
	"io"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ToJSON streams every document of cursor into a json file, see export.
func ToJSON(ctx context.Context, cursor libmongo.Cursor, opts Options) error {
	return export(ctx, cursor, opts, "json", func(w io.Writer) (recordWriter, error) {
		_, err := io.WriteString(w, "[\n")
		if err != nil {
			return nil, err
		}

		return &jsonWriter{
			w:       w,
			opts:    opts,
			encoder: json.NewEncoder(),
			record:  make(map[string]any, len(opts.Fields)),
		}, nil
	})
}

type jsonWriter struct {
	w       io.Writer
	opts    Options
	encoder encoding.Encoder
	// record is reused between documents: every field is always set.
	record  map[string]any
	written bool
}

func (w *jsonWriter) Write(ctx context.Context, doc bson.M) error {
	for _, field := range w.opts.Fields {
		v, ok, err := w.opts.value(ctx, doc, field)
		if err != nil {
			return err
		}

		if !ok {
			v = ""
		}

		w.record[field.header()] = v
	}

	if w.written {
		_, err := io.WriteString(w.w, ",\n")
		if err != nil {
			return err
		}
	}

	w.written = true

	b, err := w.encoder.Encode(w.record)
	if err != nil {
		return err
	}

	_, err = w.w.Write(b)

	return err
}

func (w *jsonWriter) Close() error {
	_, err := io.WriteString(w.w, "\n]")

	return err
}
