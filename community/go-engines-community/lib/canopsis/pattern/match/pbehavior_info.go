package match

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func ValidatePbehaviorInfoPattern(p pattern.PbehaviorInfo) bool {
	emptyPbhInfo := types.PbehaviorInfo{}

	for _, group := range p {
		if len(group) == 0 {
			return false
		}

		for _, v := range group {
			f := v.Field
			cond := v.Condition
			var err error

			if str, ok := emptyPbhInfo.GetStringField(f); ok {
				_, err = cond.MatchString(str)
			} else {
				err = pattern.ErrUnsupportedField
			}

			if err != nil {
				return false
			}
		}
	}

	return true
}

func MatchPbehaviorInfoPattern(p pattern.PbehaviorInfo, pbhInfo *types.PbehaviorInfo) (bool, error) {
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

			if str, ok := pbhInfo.GetStringField(f); ok {
				matched, err = cond.MatchString(str)
			} else {
				err = pattern.ErrUnsupportedField
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
