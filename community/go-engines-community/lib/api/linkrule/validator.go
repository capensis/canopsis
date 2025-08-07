package linkrule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apiexternaldata "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/externaldatatable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
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

	apiexternaldata.ValidateRefParameters(sl, r.ExternalData, []string{externaldata.RefTypeTable})
}

func ValidateTemplateRequest(sl validator.StructLevel) {
	var r = sl.Current().Interface().(TemplateRequest)
	switch r.Rule.Type {
	case liblink.TypeAlarm:
		if r.TestData.Entity != "" {
			sl.ReportError(r.TestData.Entity, "TestData.Entity", "Entity", "must_be_empty", "")
		}
	case liblink.TypeEntity:
		if r.TestData.Alarm != "" {
			sl.ReportError(r.TestData.Alarm, "TestData.Alarm", "Alarm", "must_be_empty", "")
		}
	}
}
