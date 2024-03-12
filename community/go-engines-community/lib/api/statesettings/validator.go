package statesettings

import (
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateEditRequest(sl validator.StructLevel)
	ValidateJUnitThresholds(sl validator.StructLevel)
	ValidateStateThreshold(sl validator.StructLevel)
}

type baseValidator struct {
	invalidRulesPatternFields           []string
	invalidInheritedEntityPatternFields []string
}

func NewValidator() Validator {
	return &baseValidator{
		invalidRulesPatternFields:           common.GetForbiddenFieldsInEntityPattern(mongo.StateSettingsMongoCollection),
		invalidInheritedEntityPatternFields: append(common.GetForbiddenFieldsInEntityPattern(mongo.StateSettingsMongoCollection), "connector"),
	}
}

func (v *baseValidator) ValidateEditRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	if r.ID == statesetting.JUnitID {
		v.validateJUnitSettings(sl, r.StateSetting)
		return
	}

	v.validateStateSetting(sl, r.StateSetting)
}

func (v *baseValidator) validateStateSetting(sl validator.StructLevel, r StateSetting) {
	switch r.Method {
	case statesetting.MethodInherited:
		if r.InheritedEntityPattern == nil {
			sl.ReportError(r.InheritedEntityPattern, "InheritedEntityPattern", "InheritedEntityPattern", "required", "")
		}

		if r.InheritedEntityPattern != nil && !match.ValidateEntityPattern(*r.InheritedEntityPattern, v.invalidInheritedEntityPatternFields) {
			sl.ReportError(r.InheritedEntityPattern, "InheritedEntityPattern", "InheritedEntityPattern", "entity_pattern", "")
		}

		if r.StateThresholds != nil {
			sl.ReportError(r.StateThresholds, "StateThresholds", "StateThresholds", "must_be_empty", "")
		}
	case statesetting.MethodDependencies:
		if r.InheritedEntityPattern != nil {
			sl.ReportError(r.InheritedEntityPattern, "InheritedEntityPattern", "InheritedEntityPattern", "must_be_empty", "")
		}

		if r.StateThresholds == nil {
			sl.ReportError(r.StateThresholds, "StateThresholds", "StateThresholds", "required", "")
		}
	default:
		sl.ReportError(r.Method, "Method", "Method", "oneof", strings.Join([]string{
			statesetting.MethodInherited,
			statesetting.MethodDependencies,
		}, ","))
	}

	if r.Title == nil {
		sl.ReportError(r.Title, "Title", "Title", "required", "")
	}

	if r.Enabled == nil {
		sl.ReportError(r.Enabled, "Enabled", "Enabled", "required", "")
	}

	if r.EntityPattern == nil {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
	}

	if r.EntityPattern != nil && !match.ValidateEntityPattern(*r.EntityPattern, v.invalidRulesPatternFields) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.JUnitThresholds != nil {
		sl.ReportError(r.JUnitThresholds, "JUnitThresholds", "JUnitThresholds", "must_be_empty", "")
	}
}

func (v *baseValidator) validateJUnitSettings(sl validator.StructLevel, r StateSetting) {
	switch r.Method {
	case statesetting.MethodWorst:
		if r.JUnitThresholds != nil {
			sl.ReportError(r.JUnitThresholds, "JUnitThresholds", "JUnitThresholds", "must_be_empty", "")
		}
	case statesetting.MethodWorstOfShare:
		if r.JUnitThresholds == nil {
			sl.ReportError(r.JUnitThresholds, "JUnitThresholds", "JUnitThresholds", "notblank", "")
		}
	default:
		sl.ReportError(r.Method, "Method", "Method", "oneof", strings.Join([]string{
			statesetting.MethodWorst,
			statesetting.MethodWorstOfShare,
		}, ","))
	}

	if r.StateThresholds != nil {
		sl.ReportError(r.StateThresholds, "StateThresholds", "StateThresholds", "must_be_empty", "")
	}

	if r.EntityPattern != nil {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "must_be_empty", "")
	}

	if r.InheritedEntityPattern != nil {
		sl.ReportError(r.InheritedEntityPattern, "InheritedEntityPattern", "InheritedEntityPattern", "must_be_empty", "")
	}

	if r.Enabled != nil {
		sl.ReportError(r.Enabled, "Enabled", "Enabled", "must_be_empty", "")
	}

	if r.Priority != 0 {
		sl.ReportError(r.Priority, "Priority", "Priority", "must_be_empty", "")
	}

	if r.Title != nil {
		sl.ReportError(r.Title, "Title", "Title", "must_be_empty", "")
	}
}

func (v *baseValidator) ValidateStateThreshold(sl validator.StructLevel) {
	r := sl.Current().Interface().(StateThreshold)

	if r.Method == statesetting.CalculationMethodShare {
		if r.Value > 99 {
			sl.ReportError(r.Value, "Value", "Value", "ltfield", "100")
		}
	}

	if r.Value < 0 {
		sl.ReportError(r.Value, "Value", "Value", "gtfield", "-1")
	}
}

func (v *baseValidator) ValidateJUnitThresholds(sl validator.StructLevel) {
	r := sl.Current().Interface().(JUnitThreshold)

	if r.Type != nil {
		switch *r.Type {
		case statesetting.TypeNumber, statesetting.TypePercentage:
			/*do nothing*/
		default:
			sl.ReportError(r.Type, "Type", "Type", "oneof", strings.Join([]string{
				strconv.Itoa(statesetting.TypeNumber),
				strconv.Itoa(statesetting.TypePercentage),
			}, ","))
		}
	}
}
