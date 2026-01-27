package linkrule

import (
	apiexternaldata "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/externaldatatable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	liblink "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	templateExecutor template.Executor
}

func NewValidator(templateExecutor template.Executor) *Validator {
	return &Validator{
		templateExecutor: templateExecutor,
	}
}

func (v *Validator) ValidateEditRequest(sl validator.StructLevel) {
	var r = sl.Current().Interface().(EditRequest)

	if r.CorporateEntityPattern == "" && len(r.EntityPattern) > 0 &&
		!match.ValidateEntityPattern(r.EntityPattern, patternfields.GetForbiddenFieldsInEntityPattern(mongo.LinkRuleMongoCollection)) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.CorporateAlarmPattern == "" && len(r.AlarmPattern) > 0 &&
		!match.ValidateAlarmPattern(r.AlarmPattern,
			patternfields.GetForbiddenFieldsInAlarmPattern(mongo.LinkRuleMongoCollection),
			patternfields.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.LinkRuleMongoCollection),
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
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
		}

		if len(r.AlarmPattern) > 0 {
			sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "excluded", "")
		}

		if r.CorporateAlarmPattern != "" {
			sl.ReportError(r.CorporateAlarmPattern, "CorporateAlarmPattern", "CorporateAlarmPattern", "excluded", "")
		}
	}

	if len(r.Links) > 0 && r.SourceCode != "" {
		sl.ReportError(r.SourceCode, "SourceCode", "SourceCode", "required_not_both", "Links")
	}
	if len(r.Links) == 0 && r.SourceCode == "" {
		sl.ReportError(r.Links, "Links", "Links", "required_or", "SourceCode")
		sl.ReportError(r.SourceCode, "SourceCode", "SourceCode", "required_or", "Links")
	}

	apiexternaldata.ValidateRefParameters(sl, v.templateExecutor, r.ExternalData, []string{externaldata.RefTypeTable})
}

func (v *Validator) ValidateTemplateRequest(sl validator.StructLevel) {
	var r = sl.Current().Interface().(TemplateRequest)
	switch r.Rule.Type {
	case liblink.TypeAlarm:
		if r.TestData.Entity != "" {
			sl.ReportError(r.TestData.Entity, "Entity", "TestData.Entity", "must_be_empty", "")
		}
	case liblink.TypeEntity:
		if r.TestData.Alarm != "" {
			sl.ReportError(r.TestData.Alarm, "Alarm", "TestData.Alarm", "must_be_empty", "")
		}
	}

	if len(r.Rule.Links) > 0 && r.Rule.SourceCode != "" {
		sl.ReportError(r.Rule.SourceCode, "SourceCode", "Rule.SourceCode", "required_not_both", "Links")
	}
	if len(r.Rule.Links) == 0 && r.Rule.SourceCode == "" {
		sl.ReportError(r.Rule.Links, "Links", "Rule.Links", "required_or", "SourceCode")
		sl.ReportError(r.Rule.SourceCode, "SourceCode", "Rule.SourceCode", "required_or", "Links")
	}
}
