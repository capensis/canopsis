package resolverule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
)

func ValidateEditRequest(sl validator.StructLevel) {
	var r = sl.Current().Interface().(EditRequest)

	if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" &&
		len(r.AlarmPattern) == 0 && r.CorporateAlarmPattern == "" {
		sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required_or", "EntityPattern")
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_or", "AlarmPattern")
	}

	if r.CorporateEntityPattern == "" && len(r.EntityPattern) > 0 &&
		!r.EntityPattern.Validate(common.GetForbiddenFieldsInEntityPattern(mongo.ResolveRuleMongoCollection)) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.CorporateAlarmPattern == "" && len(r.AlarmPattern) > 0 &&
		!r.AlarmPattern.Validate(
			common.GetForbiddenFieldsInAlarmPattern(mongo.ResolveRuleMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.ResolveRuleMongoCollection),
		) {
		sl.ReportError(r.EntityPattern, "AlarmPattern", "AlarmPattern", "alarm_pattern", "")
	}
}
