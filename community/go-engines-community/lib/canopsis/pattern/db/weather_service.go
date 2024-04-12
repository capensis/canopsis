package db

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"go.mongodb.org/mongo-driver/bson"
)

func WeatherServicePatternToMongoQuery(p pattern.WeatherServicePattern, prefix string) (bson.M, error) {
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
