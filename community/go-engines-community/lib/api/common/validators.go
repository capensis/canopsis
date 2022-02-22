package common

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
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
	if field.Type() == reflect.TypeOf(types.CpsTime{}) {
		val := field.Interface().(types.CpsTime).Time
		if val.IsZero() {
			return nil
		}

		return val
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

func ValidateAlarmPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}
	p, ok := i.(pattern.Alarm)
	if !ok {
		return false
	}

	if len(p) == 0 {
		return true
	}

	for _, group := range p {
		if len(group) == 0 {
			return false
		}
	}

	_, err := p.Match(types.Alarm{})
	if err != nil {
		return false
	}

	return true
}

func ValidateEntityPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}
	p, ok := i.(pattern.Entity)
	if !ok {
		return false
	}

	if len(p) == 0 {
		return true
	}

	for _, group := range p {
		if len(group) == 0 {
			return false
		}
	}

	_, err := p.Match(types.Entity{})
	if err != nil {
		return false
	}

	return true
}

func ValidatePbehaviorPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}
	p, ok := i.(pattern.Pbehavior)
	if !ok {
		return false
	}

	if len(p) == 0 {
		return true
	}

	for _, group := range p {
		if len(group) == 0 {
			return false
		}
	}

	_, err := p.Match(pbehavior.PBehavior{})
	if err != nil {
		return false
	}

	return true
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
		panic(fmt.Sprintf("request does not have field %s", v.field))
	}
	if field.IsZero() {
		return
	}
	fieldType, ok := sl.Current().Type().FieldByName(v.field)
	if !ok {
		panic(fmt.Sprintf("request does not have field %s", v.field))
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
	} else if err != mongodriver.ErrNoDocuments {
		panic(err)
	}
}

func NewUniqueBulkFieldValidator(field string) FieldValidator {
	return &uniqueBulkFieldValidator{
		field: field,
	}
}

type uniqueBulkFieldValidator struct {
	field string
}

func (v *uniqueBulkFieldValidator) Validate(_ context.Context, sl validator.StructLevel) {
	vals := make(map[interface{}][]int)

	var arr *reflect.Value
	fieldName := ""
	for i := 0; i < sl.Current().NumField(); i++ {
		field := sl.Current().Field(i)
		fieldName = sl.Current().Type().Field(i).Name
		k := field.Kind()
		if k == reflect.Array || k == reflect.Slice {
			arr = &field
		}
	}

	if arr == nil || sl.Current().NumField() > 1 {
		panic("request is not array")
	}

	for i := 0; i < arr.Len(); i++ {
		item := arr.Index(i)
		field := item.FieldByName(v.field)
		if !field.IsValid() {
			panic(fmt.Sprintf("request does not have field %s", v.field))
		}
		if field.IsZero() {
			continue
		}
		val := field.Interface()

		vals[val] = append(vals[val], i)
	}

	for val, indexes := range vals {
		if len(indexes) > 1 {
			for i := 1; i < len(indexes); i++ {
				path := fmt.Sprintf("%s[%d].%s", fieldName, i, v.field)
				sl.ReportError(val, path, v.field, "unique", "")
			}
		}
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
		panic(fmt.Sprintf("request does not have field %s", v.field))
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
		if err == mongodriver.ErrNoDocuments {
			sl.ReportError(val, v.field, v.field, "not_exist", "")
		} else {
			panic(err)
		}
	}
}
