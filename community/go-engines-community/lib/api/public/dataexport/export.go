package dataexport

import (
	"cmp"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// recordWriter writes the documents of an export in one format.
type recordWriter interface {
	// Write writes one document.
	Write(ctx context.Context, doc bson.M) error
	// Close writes what follows the documents.
	Close() error
}

func Filepath(dir, jobID, format string) string {
	return filepath.Join(dir, jobID+"."+format)
}

// export streams every document of cursor into the Filepath of the job,
// in the format of the writer returned by newWriter.
//
// The cursor is closed before returning, and the file is removed if anything fails,
// so a written file is always a complete export.
// Documents are written one by one: the result set is never held in memory.
func export(
	ctx context.Context,
	cursor libmongo.Cursor,
	opts Options,
	format string,
	newWriter func(io.Writer) (recordWriter, error),
) (resErr error) {
	defer func() {
		err := cursor.Close(ctx)
		if err != nil && resErr == nil {
			resErr = err
		}
	}()

	err := os.MkdirAll(opts.Dir, os.ModeDir|asyncjob.DirPerm)
	if err != nil {
		return err
	}

	file, err := os.Create(Filepath(opts.Dir, opts.Name, format))
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil && resErr == nil {
			resErr = err
		}

		if resErr != nil {
			_ = os.Remove(file.Name())
		}
	}()

	w, err := newWriter(file)
	if err != nil {
		return err
	}

	for cursor.Next(ctx) {
		var doc bson.M
		err = cursor.Decode(&doc)
		if err != nil {
			return err
		}

		if opts.Enrich != nil {
			err = opts.Enrich(ctx, doc)
			if err != nil {
				return err
			}
		}

		err = w.Write(ctx, doc)
		if err != nil {
			return err
		}
	}

	err = cursor.Err()
	if err != nil {
		return err
	}

	return w.Close()
}

// header returns the name of the field in the export, its path if it has no label.
func (f Field) header() string {
	return cmp.Or(f.Label, f.Name)
}

// value returns the value of the field in doc, transformed by opts.Transform.
// It returns false if doc doesn't have the field.
func (opts Options) value(ctx context.Context, doc bson.M, field Field) (any, bool, error) {
	v, ok := nestedValue(doc, strings.Split(field.Name, "."))
	if !ok || opts.Transform == nil {
		return v, ok, nil
	}

	v, err := opts.Transform(ctx, field.Name, v)
	if err != nil {
		return nil, false, err
	}

	return v, true, nil
}

func nestedValue(doc bson.M, keys []string) (any, bool) {
	if len(keys) == 0 {
		return nil, false
	}

	v, ok := doc[keys[0]]
	if !ok {
		return nil, false
	}

	if len(keys) == 1 {
		return v, true
	}

	switch sub := v.(type) {
	case bson.M:
		return nestedValue(sub, keys[1:])
	default:
		return nil, false
	}
}
