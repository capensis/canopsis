package pattern

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

type PbehaviorInfo [][]FieldCondition

func (p PbehaviorInfo) Match(pbhInfo types.PbehaviorInfo) (bool, error) {
	if len(p) == 0 {
		return true, nil
	}

	for _, group := range p {
		matched := false

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error
			matched = false

			if str, ok := getPbehaviorInfoStringField(pbhInfo, f); ok {
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

func (p PbehaviorInfo) Validate(forbiddenFields []string) bool {
	emptyPbhInfo := types.PbehaviorInfo{}
	forbiddenFieldsMap := make(map[string]bool, len(forbiddenFields))
	for _, field := range forbiddenFields {
		forbiddenFieldsMap[field] = true
	}

	for _, group := range p {
		if len(group) == 0 {
			return false
		}

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error

			if forbiddenFieldsMap[f] {
				return false
			}

			if str, ok := getPbehaviorInfoStringField(emptyPbhInfo, f); ok {
				_, _, err = cond.MatchString(str)
			} else {
				err = ErrUnsupportedField
			}

			if err != nil {
				return false
			}
		}
	}

	return true
}

func (p PbehaviorInfo) ToMongoQuery(prefix string) (bson.M, error) {
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

	return bson.M{"$or": groupQueries}, nil
}

func (p PbehaviorInfo) HasField(field string) bool {
	for _, group := range p {
		for _, condition := range group {
			if condition.Field == field {
				return true
			}
		}
	}

	return false
}

func getPbehaviorInfoStringField(pbhInfo types.PbehaviorInfo, f string) (string, bool) {
	switch f {
	case "pbehavior_info.id":
		return pbhInfo.ID, true
	case "pbehavior_info.type":
		return pbhInfo.TypeID, true
	case "pbehavior_info.reason":
		return pbhInfo.Reason, true
	default:
		return "", false
	}
}
