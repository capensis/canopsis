package oldpattern

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// EntityRegexMatches is a type that contains the values of the sub-expressions
// of regular expressions for each of the fields of an Entity that contain
// strings.
type EntityRegexMatches struct {
	ID             RegexMatches
	Name           RegexMatches
	Component      RegexMatches
	Category       RegexMatches
	Infos          map[string]InfoRegexMatches
	ComponentInfos map[string]InfoRegexMatches
	Type           RegexMatches
}

// NewEntityRegexMatches creates an EntityRegexMatches, with the Infos field
// initialized to an empty map.
func NewEntityRegexMatches() EntityRegexMatches {
	return EntityRegexMatches{
		Infos:          make(map[string]InfoRegexMatches),
		ComponentInfos: make(map[string]InfoRegexMatches),
	}
}

// EntityFields is a type representing a pattern that can be applied to an
// entity
type EntityFields struct {
	ID             StringPattern          `bson:"_id"`
	Name           StringPattern          `bson:"name"`
	Enabled        BoolPattern            `bson:"enabled"`
	Infos          map[string]InfoPattern `bson:"infos"`
	ComponentInfos map[string]InfoPattern `bson:"component_infos"`
	Type           StringPattern          `bson:"type"`
	Component      StringPattern          `bson:"component"`
	Category       StringPattern          `bson:"category"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

// AsMongoDriverQuery returns a mongodb filter from the EntityFields for mongo-driver.
func (e EntityFields) AsMongoDriverQuery() bson.M {
	query := bson.M{}
	if !e.ID.Empty() {
		query["_id"] = e.ID.AsMongoDriverQuery()
	}
	if !e.Name.Empty() {
		query["name"] = e.Name.AsMongoDriverQuery()
	}
	if !e.Enabled.Empty() {
		query["enabled"] = e.Enabled.AsMongoDriverQuery()
	}
	if len(e.Infos) != 0 {
		for key, value := range e.Infos {
			if value.ShouldNotBeSet {
				query["infos."+key] = nil
			} else {
				for valueKey, valueValue := range value.AsMongoDriverQuery() {
					query["infos."+key+"."+valueKey] = valueValue
				}
			}
		}
	}
	if len(e.ComponentInfos) != 0 {
		for key, value := range e.ComponentInfos {
			if value.ShouldNotBeSet {
				query["component_infos."+key] = nil
			} else {
				for valueKey, valueValue := range value.AsMongoDriverQuery() {
					query["component_infos."+key+"."+valueKey] = valueValue
				}
			}
		}
	}
	if !e.Type.Empty() {
		query["type"] = e.Type.AsMongoDriverQuery()
	}
	if !e.Component.Empty() {
		query["component"] = e.Component.AsMongoDriverQuery()
	}
	if !e.Category.Empty() {
		query["category"] = e.Category.AsMongoDriverQuery()
	}
	return query
}

func (e EntityFields) AsSqlQuery(table ...string) (string, error) {
	prefix := ""
	if len(table) > 1 {
		panic(fmt.Errorf("too many arguments, expected one: %+v", table))
	}
	if len(table) > 0 {
		prefix = table[0] + "."
	}

	conds := make([]string, 0)
	if !e.ID.Empty() {
		// Int id is used in SQL database. String id is stored to custom_id field.
		conds = append(conds, fmt.Sprintf("%scustom_id %s", prefix, e.ID.AsSqlQuery()))
	}
	if !e.Name.Empty() {
		conds = append(conds, fmt.Sprintf("%sname %s", prefix, e.Name.AsSqlQuery()))
	}
	if !e.Type.Empty() {
		conds = append(conds, fmt.Sprintf("%stype %s", prefix, e.Type.AsSqlQuery()))
	}
	if !e.Component.Empty() {
		conds = append(conds, fmt.Sprintf("%scomponent %s", prefix, e.Component.AsSqlQuery()))
	}
	if !e.Category.Empty() {
		conds = append(conds, fmt.Sprintf("%scategory %s", prefix, e.Category.AsSqlQuery()))
	}
	if !e.Enabled.Empty() {
		conds = append(conds, fmt.Sprintf("%senabled %s", prefix, e.Enabled.AsSqlQuery()))
	}
	if len(e.Infos) != 0 {
		for key, value := range e.Infos {
			if value.Value.IsSet() {
				conds = append(conds, fmt.Sprintf("%sinfos->>'%s' %s", prefix, key, value.Value.AsSqlQuery()))
			}
			if value.Name.IsSet() {
				return "", fmt.Errorf("where clause for infos.name is not supported for SQL")
			}
			if value.Description.IsSet() {
				return "", fmt.Errorf("where clause for infos.description is not supported for SQL")
			}
		}
	}
	if len(e.ComponentInfos) != 0 {
		for key, value := range e.ComponentInfos {
			if value.Value.IsSet() {
				conds = append(conds, fmt.Sprintf("%scomponent_infos->>'%s' %s", prefix, key, value.Value.AsSqlQuery()))
			}
			if value.Name.IsSet() {
				return "", fmt.Errorf("where clause for infos.name is not supported for SQL")
			}
			if value.Description.IsSet() {
				return "", fmt.Errorf("where clause for infos.description is not supported for SQL")
			}
		}
	}
	return strings.Join(conds, " AND "), nil
}

// EntityPattern is a type representing a pattern that can be applied to an
// entity
type EntityPattern struct {
	// ShouldNotBeNil is a boolean indicating that the entity should not be
	// nil, and ShouldBeNil is a boolean indicating that the entity should be
	// nil.
	// The two booleans are needed to be able to make the difference between
	// the case where no entity pattern was defined (in which case the entity
	// can be nil or not), the case where a nil entity pattern was defined (in
	// which case the entity should be nil), and the case where a non-nil
	// entity pattern was defined (in which case the entity should not be nil).
	ShouldNotBeNil bool
	ShouldBeNil    bool

	EntityFields
}

func (e EntityPattern) IsSet() bool {
	return e.ShouldBeNil ||
		e.EntityFields.Type.IsSet() ||
		e.EntityFields.Enabled.IsSet() ||
		e.EntityFields.Name.IsSet() ||
		e.EntityFields.ID.IsSet() ||
		e.EntityFields.Component.IsSet() ||
		e.EntityFields.Category.IsSet() ||
		len(e.EntityFields.Infos) > 0 ||
		len(e.EntityFields.ComponentInfos) > 0
}

// AsMongoDriverQuery returns a mongodb filter from the EntityPattern for mongo-driver
func (e EntityPattern) AsMongoDriverQuery() bson.M {
	query := make(bson.M)

	if e.ShouldBeNil {
		return nil
	} else if e.ShouldNotBeNil {
		return e.EntityFields.AsMongoDriverQuery()
	}
	return query
}

func (e EntityPattern) AsSqlQuery(table ...string) (string, error) {
	if e.ShouldNotBeNil {
		return e.EntityFields.AsSqlQuery(table...)
	}

	return "", nil
}

// Matches returns true if an entity is matched by a pattern. If the pattern
// contains regular expressions with sub-expressions, the values of the
// sub-expressions are written in the matches argument.
func (e EntityPattern) Matches(entity *types.Entity, matches *EntityRegexMatches) bool {
	if entity == nil {
		return !e.ShouldNotBeNil
	}

	match := !e.ShouldBeNil &&
		e.Component.Matches(entity.Component, &matches.Component) &&
		e.Category.Matches(entity.Category, &matches.Category) &&
		e.ID.Matches(entity.ID, &matches.ID) &&
		e.Name.Matches(entity.Name, &matches.Name) &&
		e.Enabled.Matches(entity.Enabled) &&
		e.Type.Matches(entity.Type, &matches.Type)
	if !match {
		return false
	}

	for infoName, infoPattern := range e.Infos {
		var infoRegexMatches InfoRegexMatches
		info, isSet := entity.Infos[infoName]
		match = infoPattern.Matches(info, isSet, &infoRegexMatches)
		if match {
			matches.Infos[infoName] = infoRegexMatches
		} else {
			return false
		}
	}

	for componentInfoName, componentInfoPattern := range e.ComponentInfos {
		var infoRegexMatches InfoRegexMatches
		info, isSet := entity.ComponentInfos[componentInfoName]
		match = componentInfoPattern.Matches(info, isSet, &infoRegexMatches)
		if match {
			matches.ComponentInfos[componentInfoName] = infoRegexMatches
		} else {
			return false
		}
	}

	return true
}

func (e EntityPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if e.ShouldBeNil {
		return bsontype.Null, []byte{}, nil
	}

	resultBson := bson.M{}

	if e.ID.IsSet() {
		bsonFieldName, err := GetFieldBsonName(e, "ID", "id")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = e.ID
	}

	if e.Name.IsSet() {
		bsonFieldName, err := GetFieldBsonName(e, "Name", "name")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = e.Name
	}

	if e.Enabled.IsSet() {
		bsonFieldName, err := GetFieldBsonName(e, "Enabled", "enabled")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = e.Enabled
	}

	if len(e.Infos) > 0 {
		bsonFieldName, err := GetFieldBsonName(e, "Infos", "infos")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = e.Infos
	}

	if len(e.ComponentInfos) > 0 {
		bsonFieldName, err := GetFieldBsonName(e, "ComponentInfos", "componentinfos")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = e.ComponentInfos
	}

	if e.Component.IsSet() {
		bsonFieldName, err := GetFieldBsonName(e, "Component", "component")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = e.Component
	}

	if e.Category.IsSet() {
		bsonFieldName, err := GetFieldBsonName(e, "Category", "category")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = e.Category
	}

	if e.Type.IsSet() {
		bsonFieldName, err := GetFieldBsonName(e, "Type", "type")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = e.Type
	}

	if len(resultBson) > 0 {
		return bson.MarshalValue(resultBson)
	}

	return bsontype.Undefined, nil, nil
}

func (e *EntityPattern) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Null:
		// The BSON value is null. The field should not be set.
		e.ShouldBeNil = true
		e.ShouldNotBeNil = false
		return nil
	default:
		err := bson.Unmarshal(b, &e.EntityFields)
		if err != nil {
			return err
		}

		if len(e.UnexpectedFields) != 0 {
			unexpectedFieldNames := make([]string, 0, len(e.UnexpectedFields))
			for key := range e.UnexpectedFields {
				unexpectedFieldNames = append(unexpectedFieldNames, key)
			}

			return UnexpectedFieldsError{
				Err: fmt.Errorf("unexpected pattern fields: %s", strings.Join(unexpectedFieldNames, ", ")),
			}
		}
	}

	e.ShouldBeNil = false
	e.ShouldNotBeNil = true

	return nil
}

// EntityPatternList is a type representing a list of entity patterns.
// An entity is matched by an EntityPatternList if it is matched by one of its
// EntityPatterns.
// The zero value of an EntityPatternList (i.e. an EntityPatternList that has
// not been set) is considered valid, and matches all entities.
// Deprecated : community/go-engines-community/lib/canopsis/pattern/Entity
type EntityPatternList struct {
	Patterns []EntityPattern `swaggerignore:"true"`

	// Set is a boolean indicating whether the EntityPatternList has been set
	// explicitly or not.
	Set bool `swaggerignore:"true"`

	// Valid is a boolean indicating whether the event patterns or valid or
	// not.
	// Valid is also false if the EntityPatternList has not been set.
	Valid bool `swaggerignore:"true"`
}

func (l *EntityPatternList) UnmarshalJSON(b []byte) error {
	var jsonPatterns interface{}
	err := json.Unmarshal(b, &jsonPatterns)
	if err != nil {
		return err
	}

	marshalled, err := bson.Marshal(bson.M{
		"list": jsonPatterns,
	})
	if err != nil {
		return err
	}

	var patternWrapper struct {
		PatternList EntityPatternList `bson:"list"`
	}

	err = bson.Unmarshal(marshalled, &patternWrapper)
	if err != nil {
		return err
	}

	*l = patternWrapper.PatternList
	return nil
}

func (l EntityPatternList) MarshalJSON() ([]byte, error) {
	bsonType, b, err := bson.MarshalValue(l)
	if err != nil {
		return nil, err
	}

	if bsonType == bsontype.Null {
		res, err := json.Marshal(nil)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	var unmarshalledBson []map[string]interface{}
	raw := bson.RawValue{
		Type:  bsontype.Array,
		Value: b,
	}
	err = raw.Unmarshal(&unmarshalledBson)
	if err != nil {
		return nil, err
	}

	res, err := json.Marshal(unmarshalledBson)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (l EntityPatternList) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !l.Set {
		return bsontype.Null, nil, nil
	}

	return bson.MarshalValue(l.Patterns)
}

func (l EntityPatternList) AsMongoDriverQuery() bson.M {
	var patternFilters []bson.M
	if !l.Set {
		return bson.M{}
	}
	for _, entitiesPattern := range l.Patterns {
		patternFilters = append(patternFilters, entitiesPattern.AsMongoDriverQuery())
	}
	return bson.M{"$or": patternFilters}
}

func (l EntityPatternList) AsNegativeMongoDriverQuery() bson.M {
	var patternFilters []bson.M
	if !l.Set {
		return bson.M{}
	}
	for _, entitiesPattern := range l.Patterns {
		patternFilters = append(patternFilters, entitiesPattern.AsMongoDriverQuery())
	}
	return bson.M{"$nor": patternFilters}
}

func (l EntityPatternList) AsSqlQuery(table ...string) (string, error) {
	if !l.Set {
		return "", nil
	}

	conds := make([]string, len(l.Patterns))
	var err error
	for i, v := range l.Patterns {
		conds[i], err = v.AsSqlQuery(table...)
		if err != nil {
			return "", err
		}
	}
	return strings.Join(conds, " OR "), nil
}

// Matches returns true if the entity is matched by the EntityPatternList.
func (l EntityPatternList) Matches(entity *types.Entity) bool {
	if !l.Set {
		return true
	}

	regexMatches := NewEntityRegexMatches()

	for _, pattern := range l.Patterns {
		if pattern.Matches(entity, &regexMatches) {
			return true
		}
	}

	return false
}

// IsSet returns true if the EntityPatternList has been set explicitly.
func (l EntityPatternList) IsSet() bool {
	return l.Set
}

// IsValid returns true if the EntityPatternList is valid.
func (l EntityPatternList) IsValid() bool {
	return !l.Set || l.Valid
}

func (l EntityPatternList) IsZero() bool {
	return !l.Set
}

func (l *EntityPatternList) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Null:
		l.Set = false
		l.Valid = false
	case bsontype.Array:
		l.Set = true
		l.Valid = false

		var raw bson.Raw
		err := bson.Unmarshal(b, &raw)
		if err != nil {
			return err
		}

		array, err := raw.Values()
		if err != nil {
			return err
		}

		for _, v := range array {
			document, ok := v.DocumentOK()
			if !ok {
				return fmt.Errorf("unable to parse entity pattern list element")
			}

			var pattern EntityPattern

			err = bson.Unmarshal(document, &pattern)
			if err != nil {
				if errors.As(err, &UnexpectedFieldsError{}) {
					return nil
				}

				return err
			}

			l.Patterns = append(l.Patterns, pattern)
		}

		l.Valid = true
	default:
		return fmt.Errorf("entity pattern list should be an array or nil")
	}

	return nil
}

func (l *EntityPatternList) AsInterface() (interface{}, error) {
	if l == nil {
		return nil, nil
	}
	b, err := l.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var m interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
