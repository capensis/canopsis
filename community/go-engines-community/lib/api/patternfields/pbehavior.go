package patternfields

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"github.com/go-playground/validator/v10"
)

type PbehaviorRequest struct {
	PbehaviorPattern          pattern.PbehaviorInfo `json:"pbehavior_pattern" binding:"pbehavior_pattern"`
	CorporatePbehaviorPattern string                `json:"corporate_pbehavior_pattern"`

	CorporatePattern savedpattern.SavedPattern `json:"-"`
	IsPrivate        bool                      `json:"-"`
	User             string                    `json:"-"`
}

func (r PbehaviorRequest) ToModel() savedpattern.PbehaviorPatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.PbehaviorPatternFields{
			PbehaviorPattern: r.PbehaviorPattern,
		}
	}

	return savedpattern.PbehaviorPatternFields{
		PbehaviorPattern:               r.CorporatePattern.PbehaviorPattern,
		CorporatePbehaviorPattern:      r.CorporatePattern.ID,
		CorporatePbehaviorPatternTitle: r.CorporatePattern.Title,
	}
}

func ValidatePbehaviorPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}
	p, ok := i.(pattern.PbehaviorInfo)
	if !ok {
		return false
	}

	return match.ValidatePbehaviorInfoPattern(p)
}
