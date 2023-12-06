package db

import (
	"fmt"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func PbehaviorInfoPatternToMongoQuery(p pattern.PbehaviorInfo, prefix string) (bson.M, error) {
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

			var ok bool
			if fieldCond.Field == "pbehavior_info.canonical_type" {
				condQueries[j], ok = canonicalTypeToMongoQuery(cond, mongoField)
				if ok {
					continue
				}
			}

			if _, ok = emptyPbhInfo.GetStringField(fieldCond.Field); ok {
				condQueries[j], err = cond.StringToMongoQuery(mongoField, false)
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

func canonicalTypeToMongoQuery(c pattern.Condition, f string) (bson.M, bool) {
	switch c.Type {
	case pattern.ConditionEqual:
		valueStr := c.GetValueStr()

		if valueStr != nil && *valueStr == types.PbhCanonicalTypeActive {
			return bson.M{f: bson.M{"$in": bson.A{nil, *valueStr}}}, true
		}
	case pattern.ConditionNotEqual:
		valueStr := c.GetValueStr()

		if valueStr != nil && *valueStr == types.PbhCanonicalTypeActive {
			return bson.M{f: bson.M{"$nin": bson.A{nil, *valueStr}}}, true
		}
	case pattern.ConditionIsOneOf:
		valueStrArray := c.GetValueStrArray()

		if slices.Contains(valueStrArray, types.PbhCanonicalTypeActive) {
			values := make([]interface{}, len(valueStrArray)+1)
			for k, s := range valueStrArray {
				values[k] = s
			}
			values[len(values)-1] = nil

			return bson.M{f: bson.M{"$in": values}}, true
		}
	case pattern.ConditionIsNotOneOf:
		valueStrArray := c.GetValueStrArray()

		if slices.Contains(valueStrArray, types.PbhCanonicalTypeActive) {
			values := make([]interface{}, len(valueStrArray)+1)
			for k, s := range valueStrArray {
				values[k] = s
			}
			values[len(values)-1] = nil

			return bson.M{f: bson.M{"$nin": values}}, true
		}
	}

	return nil, false
}
