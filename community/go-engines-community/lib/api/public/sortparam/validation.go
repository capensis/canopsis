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

func Validate(sort string, allowedFields []string) (bson.M, error) {
	sortable := make(map[string]bool, len(allowedFields))
	for _, field := range allowedFields {
		sortable[field] = true
	}

	return ValidateFunc(sort,
		func(f string) bool {
			return sortable[f]
		},
		strings.Join(allowedFields, " "),
	)
}

// ValidateFunc works as Validate but resolves sortable fields with isSortable instead of a fixed list.
// errParam is used in validation errors to describe the expected values.
func ValidateFunc(sort string, isSortable func(string) bool, errParam string) (bson.M, error) {
	return ValidateResolveFunc(sort,
		func(f string) (string, bool) {
			return f, isSortable(f)
		},
		errParam,
	)
}

// ValidateResolveFunc works as ValidateFunc but resolveField also maps a requested field
// to the field which is used in the sort query. It returns false if the field is not sortable.
// errParam is used in validation errors to describe the expected values.
func ValidateResolveFunc(sort string, resolveField func(string) (string, bool), errParam string) (bson.M, error) {
	parts := strings.Split(sort, ",")
	q := make(bson.D, 0, len(parts))
	errs := validator.ValidationErrors{}
	hasID := false
	for i, p := range parts {
		f, d, ok := strings.Cut(p, ":")
		if !ok {
			d = Asc
		}

		resolved, ok := resolveField(f)
		if !ok {
			errs = append(errs, validation.NewFieldErrorWithParam("oneof", strconv.Itoa(i), "sort."+strconv.Itoa(i), errParam))
			continue
		}

		if resolved == idField {
			hasID = true
		}

		switch d {
		case Asc:
			q = append(q, bson.E{Key: resolved, Value: 1})
		case Desc:
			q = append(q, bson.E{Key: resolved, Value: -1})
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
