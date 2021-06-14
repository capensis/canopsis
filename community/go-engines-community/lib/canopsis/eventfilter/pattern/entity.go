package pattern

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson/bsontype"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	mgobson "github.com/globalsign/mgo/bson"
	mongobson "go.mongodb.org/mongo-driver/bson"
)

// EntityRegexMatches is a type that contains the values of the sub-expressions
// of regular expressions for each of the fields of an Entity that contain
// strings.
type EntityRegexMatches struct {
	ID             RegexMatches
	Name           RegexMatches
	Component      RegexMatches
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

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

// AsMongoQuery returns a mongodb filter from the EntityFields for mgo-driver
func (e EntityFields) AsMongoQuery() mgobson.M {
	query := mgobson.M{}
	if !e.ID.Empty() {
		query["_id"] = e.ID.AsMongoQuery()
	}
	if !e.Name.Empty() {
		query["name"] = e.Name.AsMongoQuery()
	}
	if !e.Enabled.Empty() {
		query["enabled"] = e.Enabled.AsMongoQuery()
	}
	if len(e.Infos) != 0 {
		for key, value := range e.Infos {
			for valueKey, valueValue := range value.AsMongoQuery() {
				query["infos."+key+"."+valueKey] = valueValue
			}
		}
	}
	if len(e.ComponentInfos) != 0 {
		for key, value := range e.ComponentInfos {
			for valueKey, valueValue := range value.AsMongoQuery() {
				query["component_infos."+key+"."+valueKey] = valueValue
			}
		}
	}
	if !e.Type.Empty() {
		query["type"] = e.Type.AsMongoQuery()
	}
	if !e.Component.Empty() {
		query["component"] = e.Component.AsMongoQuery()
	}
	return query
}

// AsMongoQuery returns a mongodb filter from the EntityFields for mongo-driver
func (e EntityFields) AsMongoDriverQuery() mongobson.M {
	query := mongobson.M{}
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
			for valueKey, valueValue := range value.AsMongoDriverQuery() {
				query["infos."+key+"."+valueKey] = valueValue
			}
		}
	}
	if len(e.ComponentInfos) != 0 {
		for key, value := range e.ComponentInfos {
			for valueKey, valueValue := range value.AsMongoDriverQuery() {
				query["component_infos."+key+"."+valueKey] = valueValue
			}
		}
	}
	if !e.Type.Empty() {
		query["type"] = e.Type.AsMongoDriverQuery()
	}
	if !e.Component.Empty() {
		query["component"] = e.Component.AsMongoDriverQuery()
	}
	return query
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
		len(e.EntityFields.Infos) > 0 ||
		len(e.EntityFields.ComponentInfos) > 0
}

// AsMongoQuery returns a mongodb filter from the EntityPattern for mgo-driver
func (e EntityPattern) AsMongoQuery() mgobson.M {
	var query mgobson.M
	query = make(mgobson.M)

	if e.ShouldBeNil {
		return nil
	} else if e.ShouldNotBeNil {
		return e.EntityFields.AsMongoQuery()
	}
	return query
}

// AsMongoDriverQuery returns a mongodb filter from the EntityPattern for mongo-driver
func (e EntityPattern) AsMongoDriverQuery() mongobson.M {
	var query mongobson.M
	query = make(mongobson.M)

	if e.ShouldBeNil {
		return nil
	} else if e.ShouldNotBeNil {
		return e.EntityFields.AsMongoDriverQuery()
	}
	return query
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

// SetBSON unmarshals a BSON value into an EntityPattern.
func (e *EntityPattern) SetBSON(raw mgobson.Raw) error {
	switch raw.Kind {
	case mgobson.ElementNil:
		// The BSON value is null. The field should not be set.
		e.ShouldBeNil = true
		e.ShouldNotBeNil = false
		return nil

	default:
		err := raw.Unmarshal(&e.EntityFields)
		if err != nil {
			return err
		}

		if len(e.UnexpectedFields) != 0 {
			unexpectedFieldNames := make([]string, 0, len(e.UnexpectedFields))
			for key := range e.UnexpectedFields {
				unexpectedFieldNames = append(unexpectedFieldNames, key)
			}
			return fmt.Errorf("Unexpected pattern fields: %s", strings.Join(unexpectedFieldNames, ", "))
		}

		e.ShouldBeNil = false
		e.ShouldNotBeNil = true
		return nil
	}
}

func (e EntityPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if e.ShouldBeNil {
		return bsontype.Null, []byte{}, nil
	}

	resultBson := mongobson.M{}

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

	if e.Type.IsSet() {
		bsonFieldName, err := GetFieldBsonName(e, "Type", "type")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = e.Type
	}

	if len(resultBson) > 0 {
		return mongobson.MarshalValue(resultBson)
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
		err := mongobson.Unmarshal(b, &e.EntityFields)
		if err != nil {
			return err
		}

		if len(e.UnexpectedFields) != 0 {
			unexpectedFieldNames := make([]string, 0, len(e.UnexpectedFields))
			for key := range e.UnexpectedFields {
				unexpectedFieldNames = append(unexpectedFieldNames, key)
			}

			return UnexpectedFieldsError{
				Err: fmt.Errorf("Unexpected pattern fields: %s", strings.Join(unexpectedFieldNames, ", ")),
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
type EntityPatternList struct {
	Patterns []EntityPattern

	// Set is a boolean indicating whether the EntityPatternList has been set
	// explicitly or not.
	Set bool

	// Valid is a boolean indicating whether the event patterns or valid or
	// not.
	// Valid is also false if the EntityPatternList has not been set.
	Valid bool
}

// @todo: temporary dirty solution, each pattern should have its own marshal/unmarshal json/bson functions
func (l *EntityPatternList) UnmarshalJSON(b []byte) error {
	var jsonPatterns interface{}
	err := json.Unmarshal(b, &jsonPatterns)

	marshalled, err := mongobson.Marshal(mongobson.M{
		"list": jsonPatterns,
	})
	if err != nil {
		return err
	}

	var patternWrapper struct {
		PatternList EntityPatternList `bson:"list"`
	}

	err = mongobson.Unmarshal(marshalled, &patternWrapper)
	if err != nil {
		return err
	}

	*l = patternWrapper.PatternList
	return nil
}

func (l EntityPatternList) MarshalJSON() ([]byte, error) {
	bsonType, bson, err := mongobson.MarshalValue(l)
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
	raw := mongobson.RawValue{
		Type:  bsontype.Array,
		Value: bson,
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

	return mongobson.MarshalValue(l.Patterns)
}

// AsMongoQuery returns a mongodb filter from the EntityPatternList
func (l EntityPatternList) AsMongoQuery() mgobson.M {
	var patternFilters []mgobson.M
	if !l.Set {
		return mgobson.M{}
	}
	for _, entitiesPattern := range l.Patterns {
		patternFilters = append(patternFilters, entitiesPattern.AsMongoQuery())
	}
	return mgobson.M{"$or": patternFilters}
}

func (l EntityPatternList) AsMongoDriverQuery() mongobson.M {
	var patternFilters []mongobson.M
	if !l.Set {
		return mongobson.M{}
	}
	for _, entitiesPattern := range l.Patterns {
		patternFilters = append(patternFilters, entitiesPattern.AsMongoDriverQuery())
	}
	return mongobson.M{"$or": patternFilters}
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

// SetBSON unmarshals a BSON value into a EntityPatternList.
// If it cannot be unmarshalled, it is marked as invalid.
func (l *EntityPatternList) SetBSON(raw mgobson.Raw) error {
	switch raw.Kind {
	case mgobson.ElementNil:
		l.Set = false
		l.Valid = false
	default:
		l.Set = true

		err := raw.Unmarshal(&l.Patterns)
		if err != nil {
			log.Printf("unable to parse entity pattern list: %v", err)
			var ipatterns interface{}
			unmarshalErr := raw.Unmarshal(&ipatterns)
			log.Printf("unable to parse entity pattern list: %+v:%v", ipatterns, err)
			l.Valid = false
			return unmarshalErr
		}

		l.Valid = true
	}
	return nil
}

func (l *EntityPatternList) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Null:
		l.Set = false
		l.Valid = false
	case bsontype.Array:
		l.Set = true
		l.Valid = false

		var raw mongobson.Raw
		err := mongobson.Unmarshal(b, &raw)
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

			err = mongobson.Unmarshal(document, &pattern)
			if err != nil {
				if _, ok = err.(UnexpectedFieldsError); ok {
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
