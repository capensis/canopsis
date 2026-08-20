package entityservice

import (
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

	for i, group := range p {
		condQueries := make([]bson.M, len(group))
		for j, cond := range group {
			switch cond.Field {
			case "name", "type":
				condQueries[j], err = cond.Condition.StringToMongoQuery(cond.Field, false)
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
