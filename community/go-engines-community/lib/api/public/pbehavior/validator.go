package pbehavior

import (
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"github.com/teambition/rrule-go"
)

func ValidateEditRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	if len(r.EntityPattern) > 0 && !patternfields.ValidateEntityPattern(r.EntityPattern, mongo.PbehaviorMongoCollection) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.RRule != "" && !validateRRULE(r.RRule) {
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

func ValidatePatchRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(PatchRequest)

	if r.RRule != nil && *r.RRule != "" && !validateRRULE(*r.RRule) {
		sl.ReportError(r.RRule, "RRule", "RRule", "rrule", "")
	}

	if r.EntityPattern != nil {
		if len(r.EntityPattern) == 0 {
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
		} else if !patternfields.ValidateEntityPattern(r.EntityPattern, mongo.PbehaviorMongoCollection) {
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

func ValidateConnectorCreateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(BulkConnectorCreateRequestItem)

	if r.Stop != nil && r.Start != nil && r.Stop.Before(*r.Start) {
		sl.ReportError(r.Stop, "Stop", "Stop", "gtfield", "Start")
	}
}

func ValidateConnectorEditRequest(sl validator.StructLevel) {
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

func validateRRULE(r string) bool {
	_, err := rrule.StrToROption(r)

	return err == nil
}
