package linkrule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	liblink "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
)

func ValidateEditRequest(sl validator.StructLevel) {
	var r = sl.Current().Interface().(EditRequest)

	if r.CorporateEntityPattern == "" && len(r.EntityPattern) > 0 &&
		!match.ValidateEntityPattern(r.EntityPattern, common.GetForbiddenFieldsInEntityPattern(mongo.LinkRuleMongoCollection)) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.CorporateAlarmPattern == "" && len(r.AlarmPattern) > 0 &&
		!match.ValidateAlarmPattern(r.AlarmPattern,
			common.GetForbiddenFieldsInAlarmPattern(mongo.LinkRuleMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.LinkRuleMongoCollection),
		) {
		sl.ReportError(r.EntityPattern, "AlarmPattern", "AlarmPattern", "alarm_pattern", "")
	}

	switch r.Type {
	case liblink.TypeAlarm:
		if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" &&
			len(r.AlarmPattern) == 0 && r.CorporateAlarmPattern == "" {
			sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required_or", "EntityPattern")
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_or", "AlarmPattern")
		}
	case liblink.TypeEntity:
		if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" {
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "AlarmPattern")
		}
	}

	if len(r.Links) > 0 && r.SourceCode != "" {
		sl.ReportError(r.SourceCode, "SourceCode", "SourceCode", "required_not_both", "Links")
	}
	if len(r.Links) == 0 && r.SourceCode == "" {
		sl.ReportError(r.Links, "Links", "Links", "required_or", "SourceCode")
		sl.ReportError(r.SourceCode, "SourceCode", "SourceCode", "required_or", "Links")
	}

	for ref, params := range r.ExternalData {
		if len(params.Select) == 0 && len(params.Regexp) == 0 {
			sl.ReportError(params.Select, "ExternalData."+ref+".Select", "Select", "required_or", "Regexp")
			sl.ReportError(params.Regexp, "ExternalData."+ref+".Regexp", "Regexp", "required_or", "Select")
		}
	}
}
