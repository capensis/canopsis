package flappingrule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateCreateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)
	v.validatePatterns(sl, r.EditRequest)
}

func (v *Validator) ValidateUpdateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(UpdateRequest)
	v.validatePatterns(sl, r.EditRequest)
}

func (v *Validator) validatePatterns(sl validator.StructLevel, r EditRequest) {
	if r.CorporateEntityPattern == "" && len(r.EntityPattern) > 0 &&
		!match.ValidateEntityPattern(r.EntityPattern, common.GetForbiddenFieldsInEntityPattern(mongo.FlappingRuleMongoCollection)) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.CorporateAlarmPattern == "" && len(r.AlarmPattern) > 0 &&
		!match.ValidateAlarmPattern(r.AlarmPattern,
			common.GetForbiddenFieldsInAlarmPattern(mongo.FlappingRuleMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.FlappingRuleMongoCollection),
		) {
		sl.ReportError(r.EntityPattern, "AlarmPattern", "AlarmPattern", "alarm_pattern", "")
	}

	if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" &&
		len(r.AlarmPattern) == 0 && r.CorporateAlarmPattern == "" {
		sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required_or", "EntityPattern")
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_or", "AlarmPattern")
	}
}
