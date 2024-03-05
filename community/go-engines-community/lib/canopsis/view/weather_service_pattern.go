package view

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"go.mongodb.org/mongo-driver/bson"
)

type WeatherServicePattern [][]pattern.FieldCondition

func (p WeatherServicePattern) Validate() bool {
	for _, group := range p {
		if len(group) == 0 {
			return false
		}

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error

			switch f {
			case "is_grey":
				_, err = cond.MatchBool(true)
			case "icon", "secondary_icon":
				_, err = cond.MatchString("")
			case "state.val":
				_, err = cond.MatchInt(0)
			default:
				err = pattern.ErrUnsupportedField
			}

			if err != nil {
				return false
			}
		}
	}

	return true
}

func (p WeatherServicePattern) ToMongoQuery(prefix string) (bson.M, error) {
	if len(p) == 0 {
		return nil, nil
	}

	if prefix != "" {
		prefix += "."
	}

	groupQueries := make([]bson.M, len(p))
	var err error

	for i, group := range p {
		condQueries := make([]bson.M, len(group))
		for j, cond := range group {
			mongoField := prefix + cond.Field
			switch cond.Field {
			case "is_grey":
				condQueries[j], err = cond.Condition.BoolToMongoQuery(mongoField)
			case "icon", "secondary_icon":
				condQueries[j], err = cond.Condition.StringToMongoQuery(mongoField, false)
			case "state.val":
				condQueries[j], err = cond.Condition.IntToMongoQuery(mongoField, false)
			default:
				err = pattern.ErrUnsupportedField
			}

			if err != nil {
				return nil, fmt.Errorf("invalid condition for %q field: %w", cond.Field, err)
			}
		}

		groupQueries[i] = bson.M{"$and": condQueries}
	}

	return bson.M{"$or": groupQueries}, nil
}

func (p WeatherServicePattern) HasField(field string) bool {
	for _, group := range p {
		for _, condition := range group {
			if condition.Field == field {
				return true
			}
		}
	}

	return false
}
