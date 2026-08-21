package sortparam

import (
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	Asc  = "asc"
	Desc = "desc"
)

const idField = "_id"

func Validate(sort string, allowedFields []string, defaultSort string) (bson.M, error) {
	if sort == "" {
		return defaultQuery(defaultSort), nil
	}

	sortable := make(map[string]bool, len(allowedFields))
	for _, field := range allowedFields {
		sortable[field] = true
	}

	return ValidateFunc(sort,
		func(f string) bool {
			return sortable[f]
		},
		strings.Join(allowedFields, " "),
		defaultSort,
	)
}

// ValidateFunc works as Validate but resolves sortable fields with isSortable instead of a fixed list.
// errParam is used in validation errors to describe the expected values.
func ValidateFunc(sort string, isSortable func(string) bool, errParam, defaultSort string) (bson.M, error) {
	if sort == "" {
		return defaultQuery(defaultSort), nil
	}

	parts := strings.Split(sort, ",")
	q := make(bson.D, 0, len(parts))
	errs := validator.ValidationErrors{}
	hasID := false
	for i, p := range parts {
		f, d, ok := strings.Cut(p, ":")
		if !ok {
			d = Asc
		}

		if !isSortable(f) {
			errs = append(errs, validation.NewFieldErrorWithParam("oneof", strconv.Itoa(i), "sort."+strconv.Itoa(i), errParam))
			continue
		}

		if f == idField {
			hasID = true
		}

		switch d {
		case Asc:
			q = append(q, bson.E{Key: f, Value: 1})
		case Desc:
			q = append(q, bson.E{Key: f, Value: -1})
		default:
			errs = append(errs, validation.NewFieldErrorWithParam("oneof", strconv.Itoa(i), "sort."+strconv.Itoa(i), Asc+" "+Desc))
			continue
		}
	}

	if len(errs) > 0 {
		return nil, validation.NewError(errs, nil)
	}

	if !hasID {
		q = append(q, bson.E{Key: idField, Value: 1})
	}

	return bson.M{"$sort": q}, nil
}

func defaultQuery(defaultSort string) bson.M {
	q := bson.D{{Key: defaultSort, Value: 1}}
	if defaultSort != idField {
		q = append(q, bson.E{Key: idField, Value: 1})
	}

	return bson.M{"$sort": q}
}
