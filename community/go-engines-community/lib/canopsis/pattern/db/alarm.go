package db

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func AlarmPatternToMongoQuery(p pattern.Alarm, prefix string) (bson.M, error) {
	groupQueries, err := getAlarmPatternGroupMongoQueries(p, prefix)
	if err != nil || len(groupQueries) == 0 {
		return nil, err
	}
	return bson.M{"$or": groupQueries}, nil
}

func AlarmPatternToNegativeMongoQuery(p pattern.Alarm, prefix string) (bson.M, error) {
	groupQueries, err := getAlarmPatternGroupMongoQueries(p, prefix)
	if err != nil || len(groupQueries) == 0 {
		return nil, err
	}
	return bson.M{"$nor": groupQueries}, nil
}

func getAlarmPatternGroupMongoQueries(p pattern.Alarm, prefix string) ([]bson.M, error) {
	if len(p) == 0 {
		return nil, nil
	}

	if prefix != "" {
		prefix += "."
	}

	emptyAlarm := types.Alarm{}
	groupQueries := make([]bson.M, len(p))
	var err error

	for i, group := range p {
		condQueries := make([]bson.M, len(group))
		for j, cond := range group {
			if infoName := pattern.GetAlarmInfoName(cond.Field); infoName != "" {
				mongoField := prefix + "v.infos_array.v." + infoName

				switch cond.FieldType {
				case pattern.FieldTypeString:
					condQueries[j], err = cond.Condition.StringToMongoQuery(mongoField, true)
				case pattern.FieldTypeInt:
					condQueries[j], err = cond.Condition.IntToMongoQuery(mongoField, true)
				case pattern.FieldTypeBool:
					condQueries[j], err = cond.Condition.BoolToMongoQuery(mongoField)
				case pattern.FieldTypeStringArray:
					condQueries[j], err = cond.Condition.StringArrayToMongoQuery(mongoField, true)
				case "":
					condQueries[j], err = cond.Condition.RefToMongoQuery(mongoField)
				default:
					return nil, fmt.Errorf("invalid field type for %q field: %s", cond.Field, cond.FieldType)
				}
				if err != nil {
					return nil, fmt.Errorf("invalid condition for %q field: %w", cond.Field, err)
				}

				continue
			}

			mongoField := prefix + cond.Field
			foundField := false
			if _, ok := emptyAlarm.GetStringField(cond.Field); ok {
				foundField = true
				condQueries[j], err = cond.Condition.StringToMongoQuery(mongoField, false)
			}
			if !foundField || err != nil {
				if _, ok := emptyAlarm.GetIntField(cond.Field); ok {
					foundField = true
					condQueries[j], err = cond.Condition.IntToMongoQuery(mongoField, false)
				}
			}
			if !foundField || err != nil {
				if _, ok := emptyAlarm.GetRefField(cond.Field); ok {
					foundField = true
					condQueries[j], err = cond.Condition.RefToMongoQuery(mongoField)
				}
			}
			if !foundField || err != nil {
				if _, ok := emptyAlarm.GetTimeField(cond.Field); ok {
					foundField = true
					condQueries[j], err = cond.Condition.TimeToMongoQuery(mongoField)
				}
			}
			if !foundField || err != nil {
				if _, ok := emptyAlarm.GetDurationField(cond.Field); ok {
					foundField = true
					condQueries[j], err = cond.Condition.DurationToMongoQuery(mongoField)
				}
			}
			if !foundField || err != nil {
				if _, ok := emptyAlarm.GetStringArrayField(cond.Field); ok {
					foundField = true
					condQueries[j], err = cond.Condition.StringArrayToMongoQuery(mongoField, false)
				}
			}

			if !foundField {
				err = pattern.ErrUnsupportedField
			}
			if err != nil {
				return nil, fmt.Errorf("invalid condition for %q field: %w", cond.Field, err)
			}
		}

		groupQueries[i] = bson.M{"$and": condQueries}
	}

	return groupQueries, nil
}
