package db

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func PBehaviorInfoPatternToMongoQuery(p pattern.PBehaviorInfo, prefix string) (bson.M, error) {
	if len(p) == 0 {
		return nil, nil
	}

	if prefix != "" {
		prefix += "."
	}

	emptyPbhInfo := types.PbehaviorInfo{}
	groupQueries := make([]bson.M, len(p))
	var err error

	for i, group := range p {
		condQueries := make([]bson.M, len(group))
		for j, fieldCond := range group {
			mongoField := prefix + fieldCond.Field
			cond := fieldCond.Condition

			if fieldCond.Field == "pbehavior_info.canonical_type" {
				var ok bool

				condQueries[j], ok = cond.CanonicalTypeToMongoQuery(mongoField)
				if ok {
					continue
				}
			}

			if _, ok := emptyPbhInfo.GetStringField(fieldCond.Field); ok {
				condQueries[j], err = cond.StringToMongoQuery(mongoField)
			} else {
				err = pattern.ErrUnsupportedField
			}
			if err != nil {
				return nil, fmt.Errorf("invalid condition for %q field: %w", fieldCond.Field, err)
			}
		}

		groupQueries[i] = bson.M{"$and": condQueries}
	}

	return bson.M{"$or": groupQueries}, nil
}
