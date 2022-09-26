package oldpattern

import (
	"fmt"
	"reflect"
	"strings"
)

func GetFieldBsonName(str interface{}, fieldName string, defaultBsonName string) (string, error) {
	v := reflect.ValueOf(str)
	i := reflect.Indirect(v)

	field, ok := i.Type().FieldByName(fieldName)
	if !ok {
		return "", fmt.Errorf("couldn't find %s field", fieldName)
	}

	bsonTags := strings.Split(field.Tag.Get("bson"), ",")
	if len(bsonTags) > 0 {
		return bsonTags[0], nil
	} else {
		return defaultBsonName, nil
	}
}
