package pbehavior

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"github.com/teambition/rrule-go"
)

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateCreateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)
	if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
	}
}

func (v *Validator) ValidateUpdateRequest(sl validator.StructLevel) {
	corporateEntityPattern := ""
	var entityPattern pattern.Entity
	switch r := sl.Current().Interface().(type) {
	case UpdateRequest:
		entityPattern = r.EntityPattern
		corporateEntityPattern = r.CorporateEntityPattern
	case BulkUpdateRequestItem:
		entityPattern = r.EntityPattern
		corporateEntityPattern = r.CorporateEntityPattern
	}

	if len(entityPattern) == 0 && corporateEntityPattern == "" {
		sl.ReportError(entityPattern, "EntityPattern", "EntityPattern", "required", "")
	}
}

func (v *Validator) ValidateEditRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	if r.CorporateEntityPattern == "" && len(r.EntityPattern) > 0 &&
		!match.ValidateEntityPattern(r.EntityPattern, patternfields.GetForbiddenFieldsInEntityPattern(mongo.PbehaviorMongoCollection)) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.RRule != "" && !v.checkRrule(r.RRule) {
		sl.ReportError(r.RRule, "RRule", "RRule", "rrule", "")
	}

	if r.Stop != nil && r.Start != nil && r.Stop.Before(*r.Start) {
		sl.ReportError(r.Stop, "Stop", "Stop", "gtfield", "Start")
	}

	if r.Timezone != "" {
		if _, err := time.LoadLocation(r.Timezone); err != nil {
			sl.ReportError(r.Timezone, "Timezone", "Timezone", "timezone", "")
		}
	}
}

func (v *Validator) ValidatePatchRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(PatchRequest)

	if r.RRule != nil && *r.RRule != "" && !v.checkRrule(*r.RRule) {
		sl.ReportError(r.RRule, "RRule", "RRule", "rrule", "")
	}

	if r.CorporateEntityPattern != nil && *r.CorporateEntityPattern == "" {
		sl.ReportError(r.CorporateEntityPattern, "CorporateEntityPattern", "CorporateEntityPattern", "required", "")
	}

	if r.CorporateEntityPattern == nil && r.EntityPattern != nil {
		if len(r.EntityPattern) == 0 {
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
		} else if !match.ValidateEntityPattern(r.EntityPattern, patternfields.GetForbiddenFieldsInEntityPattern(mongo.PbehaviorMongoCollection)) {
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
		}
	}

	if r.Name != nil && *r.Name == "" {
		sl.ReportError(r.Name, "Name", "Name", "required", "")
	}

	if r.Type != nil && *r.Type == "" {
		sl.ReportError(r.Type, "Type", "Type", "required", "")
	}

	if r.Reason != nil && *r.Reason == "" {
		sl.ReportError(r.Reason, "Reason", "Reason", "required", "")
	}

	if r.Color != nil && *r.Color != "" {
		err := validator.New().Var(*r.Color, "iscolor")
		if err != nil {
			sl.ReportError(r.Color, "Color", "Color", "iscolor", "")
		}
	}

	if r.Start != nil && r.Stop.isSet && r.Stop.val != nil && *r.Start >= *r.Stop.val {
		sl.ReportError(r.Stop, "Stop", "Stop", "gtfield", "Start")
	}

	if r.Timezone != nil && *r.Timezone != "" {
		if _, err := time.LoadLocation(*r.Timezone); err != nil {
			sl.ReportError(r.Timezone, "Timezone", "Timezone", "timezone", "")
		}
	}
}

func (v *Validator) ValidateConnectorCreateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(BulkConnectorCreateRequestItem)

	if r.Stop != nil && r.Start != nil && r.Stop.Before(*r.Start) {
		sl.ReportError(r.Stop, "Stop", "Stop", "gtfield", "Start")
	}
}

func (v *Validator) ValidateConnectorEditRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(BulkConnectorEditRequestItem)

	if r.Action == BulkConnectorActionCreate {
		if r.Name == "" {
			sl.ReportError(r.Name, "Name", "Name", "required", "")
		}

		if r.Reason == "" {
			sl.ReportError(r.Reason, "Reason", "Reason", "required", "")
		}

		if r.Type == "" {
			sl.ReportError(r.Type, "Type", "Type", "required", "")
		}
	}

	if r.Stop != nil && r.Start != nil && r.Stop.Before(*r.Start) {
		sl.ReportError(r.Stop, "Stop", "Stop", "gtfield", "Start")
	}
}

func (v *Validator) checkRrule(r string) bool {
	_, err := rrule.StrToROption(r)
	return err == nil
}

func (v *Validator) ValidateCalendarRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CalendarByEntityIDRequest)
	if r.To.Unix() > 0 && r.From.Unix() > 0 && r.To.Before(r.From) {
		sl.ReportError(r.To, "To", "To", "gtfield", "From")
	}
}
