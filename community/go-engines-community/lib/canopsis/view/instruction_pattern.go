package view

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

type InstructionPattern [][]pattern.FieldCondition

func (p InstructionPattern) Validate() bool {
	for _, group := range p {
		if len(group) == 0 {
			return false
		}

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error

			switch f {
			case "instruction_execution":
				_, err = cond.MatchRef(nil)
			case "instructions":
				switch cond.Type {
				case pattern.ConditionHasOneOf, pattern.ConditionHasNot:
					_, err = cond.MatchStringArray([]string{})
				default:
					err = pattern.ErrUnsupportedConditionType
				}
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

func (p InstructionPattern) HasField(field string) bool {
	for _, group := range p {
		for _, condition := range group {
			if condition.Field == field {
				return true
			}
		}
	}

	return false
}
