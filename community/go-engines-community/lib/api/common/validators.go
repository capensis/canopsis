package common

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/validator"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const (
	SortAsc     = "asc"
	SortDesc    = "desc"
	MaxIDLength = 255
)

var timeFormats = map[string]string{
	"YYYY-MM-DDThh:mm:ss":        "2006-01-02T15:04:05",
	"YYYY-MM-DDThh:mm:ssZ":       "2006-01-02T15:04:05-0700",
	"DD MMM YYYY hh:mm:ss":       "02 Jan 2006 15:04",
	"DD MMM YYYY hh:mm:ss ZZ":    "02 Jan 2006 15:04 MST",
	"W, DD MMM YYYY hh:mm:ss ZZ": "Mon, 02 Jan 2006 15:04:05 MST",
}

// ValidateCpsTimeType implements CustomTypeFunc and returns value to validate.
func ValidateCpsTimeType(field reflect.Value) interface{} {
	if field.Type() == reflect.TypeOf(datetime.CpsTime{}) {
		if t, ok := field.Interface().(datetime.CpsTime); ok {
			val := t.Time
			if val.IsZero() {
				return nil
			}

			return val
		}
	}

	return nil
}

func ValidateOneOfOrEmpty(fl validator.FieldLevel) bool {
	vals := strings.Split(fl.Param(), " ")
	field := fl.Field()

	var v string
	switch field.Kind() {
	case reflect.String:
		v = field.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = strconv.FormatUint(field.Uint(), 10)
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}

	if v == "" {
		return true
	}

	for i := 0; i < len(vals); i++ {
		prefix := strings.TrimSuffix(vals[i], "*")
		if prefix != "" && prefix != vals[i] {
			if strings.HasPrefix(v, prefix) {
				return true
			}
		} else if vals[i] == v {
			return true
		}
	}

	return false
}

func ValidateColorOrEmpty(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v == "" {
		return true
	}
	return validator.New().Var(v, "iscolor") == nil
}

func ValidateID(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v == "" {
		return true
	}

	return !strings.ContainsAny(v, libvalidator.InvalidIDChars) && len(v) <= MaxIDLength
}

func ValidateTimeFormat(fl validator.FieldLevel) bool {
	v := fl.Field().String()

	return v == "" || timeFormats[v] != ""
}

func GetRealFormatTime(f string) string {
	if f == "" {
		return ""
	}
	return timeFormats[f]
}

func ValidateFilteredQuery(sl validator.StructLevel) {
	r := sl.Current().Interface().(pagination.FilteredQuery)
	// Validate sort
	if r.Sort != "" {
		sorts := []string{
			SortAsc,
			SortDesc,
		}

		found := false
		for _, sort := range sorts {
			if sort == r.Sort {
				found = true
			}
		}

		if !found {
			param := strings.Join(sorts, " ")
			sl.ReportError(r.Sort, "Sort", "sort", "oneof", param)
		}
	}
}

type FieldValidator interface {
	Validate(ctx context.Context, sl validator.StructLevel)
}

func NewUniqueFieldValidator(
	db mongo.DbClient,
	collection string,
	field string,
) FieldValidator {
	return &uniqueFieldValidator{
		dbClient:     db,
		dbCollection: db.Collection(collection),
		field:        field,
	}
}

type uniqueFieldValidator struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
	field        string
}

func (v *uniqueFieldValidator) Validate(ctx context.Context, sl validator.StructLevel) {
	idField := sl.Current().FieldByNameFunc(func(name string) bool {
		return strings.ToLower(name) == "id"
	})
	id := ""
	if idField.IsValid() {
		var ok bool
		id, ok = idField.Interface().(string)
		if !ok {
			panic("request field id is not string")
		}
	}
	field := sl.Current().FieldByName(v.field)
	if !field.IsValid() {
		panic("request does not have field s" + v.field)
	}
	if field.IsZero() {
		return
	}
	fieldType, ok := sl.Current().Type().FieldByName(v.field)
	if !ok {
		panic("request does not have field %" + v.field)
	}
	val := field.Interface()
	var found struct {
		ID string `bson:"_id"`
	}
	err := v.dbCollection.FindOne(ctx, bson.M{fieldType.Tag.Get("json"): val}).Decode(&found)
	if err == nil {
		if found.ID != id || strings.ToLower(v.field) == "id" {
			sl.ReportError(val, v.field, v.field, "unique", "")
		}
	} else if !errors.Is(err, mongodriver.ErrNoDocuments) {
		panic(err)
	}
}

func NewExistFieldValidator(
	db mongo.DbClient,
	collection string,
	field string,
) FieldValidator {
	return &existFieldValidator{
		dbClient:     db,
		dbCollection: db.Collection(collection),
		field:        field,
	}
}

type existFieldValidator struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
	field        string
}

func (v *existFieldValidator) Validate(ctx context.Context, sl validator.StructLevel) {
	field := sl.Current().FieldByName(v.field)
	if !field.IsValid() {
		panic("request does not have field " + v.field)
	}
	val, ok := field.Interface().(string)
	if !ok {
		panic(fmt.Sprintf("request field %s is not string", v.field))
	}

	if val == "" {
		return
	}

	var found struct {
		ID string `bson:"_id"`
	}
	err := v.dbCollection.FindOne(ctx, bson.M{"_id": val}).Decode(&found)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			sl.ReportError(val, v.field, v.field, "not_exist", "")
		} else {
			panic(err)
		}
	}
}

func ValidateInfoValue(fl validator.FieldLevel) bool {
	return types.IsInfoValueValid(fl.Field().Interface())
}
