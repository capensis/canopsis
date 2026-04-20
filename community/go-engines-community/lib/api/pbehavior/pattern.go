package pbehavior

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type searchPattern [][]pattern.FieldCondition

func (p searchPattern) ToMongoQuery() (bson.M, error) {
	if len(p) == 0 {
		return nil, nil
	}

	groupQueries := make([]bson.M, len(p))
	var err error

	now := time.Now()

	for i, group := range p {
		condQueries := make([]bson.M, len(group))
		for j, cond := range group {
			switch cond.Field {
			case "name", "author", "rrule", "reason", "type":
				condQueries[j], err = cond.Condition.StringToMongoQuery(cond.Field, false)
			case "enabled":
				condQueries[j], err = cond.Condition.BoolToMongoQuery(cond.Field)
			case "tstart", "tstop", "rrule_end", "created", "updated", "last_alarm_date":
				condQueries[j], err = cond.Condition.TimeToMongoQuery(cond.Field, now)
			case "alarm_count":
				condQueries[j], err = cond.Condition.IntToMongoQuery(cond.Field, true)
			default:
				return nil, pattern.ErrUnsupportedField
			}

			if err != nil {
				return nil, err
			}
		}

		groupQueries[i] = bson.M{"$and": condQueries}
	}

	return bson.M{"$or": groupQueries}, nil
}
