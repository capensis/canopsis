package pattern

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func Match(
	entity types.Entity,
	alarm types.Alarm,
	entityPattern Entity,
	alarmPattern Alarm,
	oldEntityPatterns oldpattern.EntityPatternList,
	oldAlarmPatterns oldpattern.AlarmPatternList,
) (bool, error) {
	if !oldEntityPatterns.IsSet() && len(entityPattern) == 0 &&
		!oldAlarmPatterns.IsSet() && len(alarmPattern) == 0 {
		return false, nil
	}

	if len(entityPattern) > 0 {
		ok, _, err := entityPattern.Match(entity)
		if err != nil {
			return false, fmt.Errorf("entity pattern is invalid: %w", err)
		}
		if !ok {
			return false, nil
		}
	} else if oldEntityPatterns.IsSet() {
		if !oldEntityPatterns.IsValid() {
			return false, ErrInvalidOldEntityPattern
		}
		if !oldEntityPatterns.Matches(&entity) {
			return false, nil
		}
	}

	if len(alarmPattern) > 0 {
		ok, err := alarmPattern.Match(alarm)
		if err != nil {
			return false, fmt.Errorf("alarm pattern is invalid: %w", err)
		}
		if !ok {
			return false, nil
		}
	} else if oldAlarmPatterns.IsSet() {
		if !oldAlarmPatterns.IsValid() {
			return false, ErrInvalidOldAlarmPattern
		}
		if !oldAlarmPatterns.Matches(&alarm) {
			return false, nil
		}
	}

	return true, nil
}

func EntityPatternToMongoQuery(
	prefix string,
	entityPattern Entity,
	oldEntityPatterns oldpattern.EntityPatternList,
) (bson.M, error) {
	if len(entityPattern) > 0 {
		return entityPattern.ToMongoQuery(prefix)
	}

	if oldEntityPatterns.IsSet() {
		if !oldEntityPatterns.IsValid() {
			return nil, ErrInvalidOldEntityPattern
		}

		return addPrefixToOldPatternQuery(prefix, oldEntityPatterns.AsMongoDriverQuery()), nil
	}

	return nil, nil
}

func AlarmPatternToMongoQuery(
	prefix string,
	alarmPattern Alarm,
	oldAlarmPatterns oldpattern.AlarmPatternList,
) (bson.M, error) {
	if len(alarmPattern) > 0 {
		return alarmPattern.ToMongoQuery(prefix)
	}

	if oldAlarmPatterns.IsSet() {
		if !oldAlarmPatterns.IsValid() {
			return nil, ErrInvalidOldEntityPattern
		}

		return addPrefixToOldPatternQuery(prefix, oldAlarmPatterns.AsMongoDriverQuery()), nil
	}

	return nil, nil
}

func addPrefixToOldPatternQuery(prefix string, patternBson bson.M) bson.M {
	if prefix == "" {
		return patternBson
	}

	newBson := make(bson.M)
	patternListInterface, ok := patternBson["$or"]
	if !ok {
		return patternBson
	}

	patternList := patternListInterface.([]bson.M)
	newPatternsList := make([]bson.M, len(patternList))
	for i, pattern := range patternList {
		newPattern := make(bson.M)
		for k, vv := range pattern {
			newPattern[prefix+"."+k] = vv
		}

		newPatternsList[i] = newPattern
	}

	// Just in case when an entity pattern's function AsMongoDriverQuery returns an empty bson array
	// that might happen if the pattern has bad format in mongo, but after unmarshalling it has isSet = true
	// since it's not possible for $or has empty array we just fill it with an empty bson
	if len(newPatternsList) == 0 {
		newPatternsList = append(newPatternsList, bson.M{})
	}

	newBson["$or"] = newPatternsList

	return newBson
}
