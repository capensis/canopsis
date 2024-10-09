package db

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func EventPatternToMongoQuery(p pattern.Event, prefix string) (bson.M, error) {
	groupQueries, err := getEventPatternGroupMongoQueries(p, prefix)
	if err != nil || len(groupQueries) == 0 {
		return nil, err
	}
	return bson.M{"$or": groupQueries}, nil
}

func getEventPatternGroupMongoQueries(p pattern.Event, prefix string) ([]bson.M, error) {
	if len(p) == 0 {
		return nil, nil
	}

	if prefix != "" {
		prefix += "."
	}

	emptyEvent := &types.Event{}
	groupQueries := make([]bson.M, len(p))
	var err error

	for i, group := range p {
		condQueries := make([]bson.M, len(group))
		for j, cond := range group {
			mongoField := prefix + cond.Field
			if _, ok := emptyEvent.GetStringField(cond.Field); ok {
				condQueries[j], err = cond.Condition.StringToMongoQuery(mongoField, false)
			} else if _, ok := emptyEvent.GetIntField(cond.Field); ok {
				condQueries[j], err = cond.Condition.IntToMongoQuery(mongoField, false)
			} else if extraName := pattern.GetEventExtraInfoName(cond.Field); extraName != "" {
				mongoField = prefix + "extra_infos." + extraName

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
					err = fmt.Errorf("invalid field type for %q field: %s", cond.Field, cond.FieldType)
				}

			} else {
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
