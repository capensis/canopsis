package dataexport

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Transform converts a value before it's written, see Options.Transform.
type Transform func(ctx context.Context, field string, v any) (any, error)

// IsDate reports whether the field holds a timestamp.
type IsDate func(ctx context.Context, field string) (bool, error)

// WalkObjects returns a Transform applying transform to every value of a document:
// a value is transformed, then, if it's still an object, its keys are transformed as fields themselves.
// "pbehavior_info.timestamp" is therefore transformed whether it's exported on its own or as part
// of an exported "pbehavior_info" object.
//
// Objects are transformed in place, arrays aren't walked.
func WalkObjects(transform Transform) Transform {
	var walk Transform

	walk = func(ctx context.Context, field string, v any) (any, error) {
		v, err := transform(ctx, field, v)
		if err != nil {
			return nil, err
		}

		m, ok := v.(bson.M)
		if !ok {
			return v, nil
		}

		for k, sub := range m {
			sub, err = walk(ctx, field+"."+k, sub)
			if err != nil {
				return nil, err
			}

			m[k] = sub
		}

		return m, nil
	}

	return walk
}

// NewDateTransform returns a Transform formatting the timestamp of every isDate field with layout, in location.
//
// Any other field is left unchanged, so the transform can be chained with the ones of the exported resource.
// isDate is called at most once per field.
func NewDateTransform(isDate IsDate, location *time.Location, layout string) Transform {
	t := dateTransform{
		isDate:     isDate,
		location:   location,
		layout:     layout,
		dateFields: map[string]bool{},
	}

	return t.transform
}

type dateTransform struct {
	isDate   IsDate
	location *time.Location
	layout   string
	// dateFields caches the isDate results by field.
	dateFields map[string]bool
}

func (t dateTransform) transform(ctx context.Context, field string, v any) (any, error) {
	isDate, err := t.fieldIsDate(ctx, field)
	if err != nil {
		return nil, err
	}

	if !isDate {
		return v, nil
	}

	i, ok := toInt64(v)
	if !ok {
		return v, nil
	}

	return datetime.NewCpsTime(i).In(t.location).Time.Format(t.layout), nil
}

func (t dateTransform) fieldIsDate(ctx context.Context, field string) (bool, error) {
	isDate, ok := t.dateFields[field]
	if ok {
		return isDate, nil
	}

	isDate, err := t.isDate(ctx, field)
	if err != nil {
		return false, err
	}

	t.dateFields[field] = isDate

	return isDate, nil
}

func toInt64(v any) (int64, bool) {
	switch i := v.(type) {
	case int64:
		return i, true
	case int32:
		return int64(i), true
	case int:
		return int64(i), true
	case float32:
		return int64(i), true
	case float64:
		return int64(i), true
	default:
		return 0, false
	}
}
