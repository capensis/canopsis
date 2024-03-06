package pattern

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

const pbhCanonicalTypeActive = "active"

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
				matched, err = cond.MatchString(str)
			} else {
				err = ErrUnsupportedField
			}

			if err != nil {
				return false, fmt.Errorf("invalid condition for %q field: %w", f, err)
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

func (p PbehaviorInfo) Validate() bool {
	emptyPbhInfo := types.PbehaviorInfo{}

	for _, group := range p {
		if len(group) == 0 {
			return false
		}

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error

			if str, ok := getPbehaviorInfoStringField(emptyPbhInfo, f); ok {
				_, err = cond.MatchString(str)
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

	emptyPbhInfo := types.PbehaviorInfo{}
	groupQueries := make([]bson.M, len(p))
	var err error

	for i, group := range p {
		condQueries := make([]bson.M, len(group))
		for j, fieldCond := range group {
			mongoField := prefix + fieldCond.Field
			cond := fieldCond.Condition

			if fieldCond.Field == "pbehavior_info.canonical_type" {
				switch cond.Type {
				case ConditionEqual:
					if cond.valueStr != nil && *cond.valueStr == pbhCanonicalTypeActive {
						condQueries[j] = bson.M{mongoField: bson.M{"$in": bson.A{nil, *cond.valueStr}}}
						continue
					}
				case ConditionNotEqual:
					if cond.valueStr != nil && *cond.valueStr == pbhCanonicalTypeActive {
						condQueries[j] = bson.M{mongoField: bson.M{"$nin": bson.A{nil, *cond.valueStr}}}
						continue
					}
				case ConditionIsOneOf:
					found := false
					for _, item := range cond.valueStrArray {
						if item == pbhCanonicalTypeActive {
							found = true
							break
						}
					}

					if found {
						values := make([]interface{}, len(cond.valueStrArray)+1)
						for k, s := range cond.valueStrArray {
							values[k] = s
						}
						values[len(values)-1] = nil
						condQueries[j] = bson.M{mongoField: bson.M{"$in": values}}
						continue
					}
				case ConditionIsNotOneOf:
					found := false
					for _, item := range cond.valueStrArray {
						if item == pbhCanonicalTypeActive {
							found = true
							break
						}
					}

					if found {
						values := make([]interface{}, len(cond.valueStrArray)+1)
						for k, s := range cond.valueStrArray {
							values[k] = s
						}
						values[len(values)-1] = nil
						condQueries[j] = bson.M{mongoField: bson.M{"$nin": values}}
						continue
					}
				}
			}

			if _, ok := getPbehaviorInfoStringField(emptyPbhInfo, fieldCond.Field); ok {
				condQueries[j], err = cond.StringToMongoQuery(mongoField, false)
			} else {
				err = ErrUnsupportedField
			}
			if err != nil {
				return nil, fmt.Errorf("invalid condition for %q field: %w", fieldCond.Field, err)
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
	case "pbehavior_info.canonical_type":
		if pbhInfo.CanonicalType == "" {
			return pbhCanonicalTypeActive, true
		}
		return pbhInfo.CanonicalType, true
	case "pbehavior_info.reason":
		return pbhInfo.ReasonID, true
	default:
		return "", false
	}
}
