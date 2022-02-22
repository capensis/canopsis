package pattern

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"go.mongodb.org/mongo-driver/bson"
)

type Pbehavior [][]FieldCondition

func (p Pbehavior) Match(pbh pbehavior.PBehavior) (bool, error) {
	for _, group := range p {
		matched := false

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error
			matched = false

			if str, ok := getPbehaviorStringField(pbh, f); ok {
				matched, _, err = cond.MatchString(str)
			} else {
				err = ErrUnsupportedField
			}

			if err != nil {
				return false, err
			}

			if !matched {
				break
			}
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}

func (p Pbehavior) ToMongoQuery(prefix string) ([]bson.M, error) {
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
			f := prefix + cond.Field
			condQueries[j], err = cond.Condition.ToMongoQuery(f)
			if err != nil {
				return nil, err
			}
		}

		groupQueries[i] = bson.M{"$and": condQueries}
	}

	return []bson.M{{"$match": bson.M{"$or": groupQueries}}}, nil
}

func getPbehaviorStringField(pbh pbehavior.PBehavior, f string) (string, bool) {
	switch f {
	case "_id":
		return pbh.ID, true
	case "type":
		return pbh.Type, true
	case "reason":
		return pbh.Reason, true
	default:
		return "", false
	}
}
