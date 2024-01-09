package idlerule

import (
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
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
	v.validateEditRequest(sl, r.EditRequest)
}

func (v *Validator) ValidateUpdateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(UpdateRequest)
	v.validateEditRequest(sl, r.EditRequest)
}

func (v *Validator) ValidateBulkUpdateRequestItem(sl validator.StructLevel) {
	r := sl.Current().Interface().(BulkUpdateRequestItem)
	v.validateEditRequest(sl, r.EditRequest)
}

func (v *Validator) validateEditRequest(sl validator.StructLevel, r EditRequest) {
	v.validateType(sl, r.Type)
	v.validateAlarmRule(sl, r)
	v.validateEntityRule(sl, r)
	v.validateDisableDuringPeriods(sl, r.DisableDuringPeriods)
}

func (v *Validator) validateType(sl validator.StructLevel, t string) {
	if t == "" {
		return
	}
	validTypes := []string{
		idlerule.RuleTypeAlarm,
		idlerule.RuleTypeEntity,
	}
	param := strings.Join(validTypes, " ")
	found := false
	for _, v := range validTypes {
		if v == t {
			found = true
			break
		}
	}

	if !found {
		sl.ReportError(t, "Type", "Type", "oneof", param)
	}
}

func (v *Validator) validateAlarmRule(sl validator.StructLevel, r EditRequest) {
	if r.Type != idlerule.RuleTypeAlarm {
		return
	}

	v.validateAlarmRulePatterns(sl, r)

	if r.AlarmCondition == "" {
		sl.ReportError(r.AlarmCondition, "AlarmCondition", "AlarmCondition", "required", "")
	} else {
		validValues := []string{
			idlerule.RuleAlarmConditionLastEvent,
			idlerule.RuleAlarmConditionLastUpdate,
		}
		param := strings.Join(validValues, " ")
		found := false
		for _, v := range validValues {
			if v == r.AlarmCondition {
				found = true
				break
			}
		}

		if !found {
			sl.ReportError(r.AlarmCondition, "AlarmCondition", "AlarmCondition", "oneof", param)
		}
	}

	if r.Operation == nil {
		sl.ReportError(r.Operation, "Operation", "Operation", "required", "")
	} else {
		validOperationTypes := []string{
			types.ActionTypeAck,
			types.ActionTypeAckRemove,
			types.ActionTypeCancel,
			types.ActionTypeAssocTicket,
			types.ActionTypeChangeState,
			types.ActionTypeSnooze,
			types.ActionTypePbehavior,
		}
		param := strings.Join(validOperationTypes, " ")
		found := false
		for _, v := range validOperationTypes {
			if v == r.Operation.Type {
				found = true
				break
			}
		}

		if !found {
			sl.ReportError(r.Operation.Type, "Operation.Type", "Type", "oneof", param)
		}

		v.validateOperationParametersRequest(sl, r.Operation.Type, r.Operation.Parameters)
	}
}

func (v *Validator) validateEntityRule(sl validator.StructLevel, r EditRequest) {
	if r.Type != idlerule.RuleTypeEntity {
		return
	}

	v.validateEntityRulePatterns(sl, r)

	if r.Operation != nil {
		sl.ReportError(r.Operation, "Operation", "Operation", "must_be_empty", "")
	}

	if r.AlarmCondition != "" {
		sl.ReportError(r.AlarmCondition, "AlarmCondition", "AlarmCondition", "must_be_empty", "")
	}
}

func (v *Validator) validateDisableDuringPeriods(sl validator.StructLevel, disableDuringPeriods []string) {
	if len(disableDuringPeriods) == 0 {
		return
	}
	validPeriods := []string{
		pbehavior.TypeMaintenance,
		pbehavior.TypePause,
		pbehavior.TypeInactive,
	}
	param := strings.Join(validPeriods, " ")
	for _, period := range disableDuringPeriods {
		found := false
		for _, v := range validPeriods {
			if v == period {
				found = true
				break
			}
		}

		if !found {
			sl.ReportError(disableDuringPeriods, "DisableDuringPeriods", "disableDuringPeriods", "oneof", param)
		}
	}
}

func (v *Validator) validateOperationParametersRequest(sl validator.StructLevel, t string, params idlerule.Parameters) {
	switch t {
	case types.ActionTypeAssocTicket:
		if params.Ticket == "" {
			sl.ReportError(params.Ticket, "Operation.Parameters.Ticket", "Ticket", "required", "")
		}
	case types.ActionTypeChangeState:
		if params.State == nil {
			sl.ReportError(params.State, "Operation.Parameters.State", "State", "required", "")
		} else {
			validTypes := []types.CpsNumber{
				types.AlarmStateOK,
				types.AlarmStateMinor,
				types.AlarmStateMajor,
				types.AlarmStateCritical,
			}
			param := ""
			for i := range validTypes {
				param += strconv.Itoa(int(validTypes[i]))
				if i < len(validTypes)-1 {
					param += " "
				}
			}

			found := false
			for _, v := range validTypes {
				if v == *params.State {
					found = true
					break
				}
			}

			if !found {
				sl.ReportError(params.State, "Operation.Parameters.State", "State", "oneof", param)
			}
		}
	case types.ActionTypeSnooze:
		if params.Duration == nil {
			sl.ReportError(params.Duration, "Operation.Parameters.Duration", "Duration", "required", "")
		}
	case types.ActionTypePbehavior:
		if params.Name == "" {
			sl.ReportError(params.Name, "Operation.Parameters.Name", "Name", "required", "")
		}
		if params.Reason == "" {
			sl.ReportError(params.Reason, "Operation.Parameters.Reason", "Reason", "required", "")
		}
		if params.Type == "" {
			sl.ReportError(params.Type, "Operation.Parameters.Type", "Type", "required", "")
		}
		// Validate rrule
		if params.RRule != "" {
			_, err := rrule.StrToROption(params.RRule)
			if err != nil {
				sl.ReportError(params.RRule, "Operation.Parameters.RRule", "RRule", "rrule", "")
			}
		}

		// Validate time
		if params.Tstart == nil && params.Tstop != nil {
			sl.ReportError(params.Tstart, "Operation.Parameters.Tstart", "Tstart", "required_with", "Tstop")
		}
		if params.Tstart != nil && params.Tstop == nil {
			sl.ReportError(params.Tstop, "Operation.Parameters.Tstop", "Tstop", "required_with", "Tstart")
		}
		if params.Duration == nil && params.StartOnTrigger != nil && *params.StartOnTrigger {
			sl.ReportError(params.Duration, "Operation.Parameters.Duration", "Duration", "required_with", "StartOnTrigger")
		}
		if params.Duration != nil && (params.StartOnTrigger == nil || !*params.StartOnTrigger) {
			sl.ReportError(params.StartOnTrigger, "Operation.Parameters.StartOnTrigger", "StartOnTrigger", "required_with", "Duration")
		}
		if params.Tstart == nil && params.Tstop == nil && params.Duration == nil && (params.StartOnTrigger == nil || !*params.StartOnTrigger) {
			sl.ReportError(params.Tstart, "Operation.Parameters.Tstart", "Tstart", "required_or", "StartOnTrigger")
			sl.ReportError(params.StartOnTrigger, "Operation.Parameters.StartOnTrigger", "StartOnTrigger", "required_or", "Tstart")
		}
		if params.Tstart != nil && params.StartOnTrigger != nil && *params.StartOnTrigger {
			sl.ReportError(params.Tstart, "Operation.Parameters.Tstart", "Tstart", "required_or", "StartOnTrigger")
			sl.ReportError(params.StartOnTrigger, "Operation.Parameters.StartOnTrigger", "StartOnTrigger", "required_or", "Tstart")
		}
		if params.Tstart != nil && params.Tstop != nil && params.Tstop.Before(*params.Tstart) {
			sl.ReportError(params.Tstop, "Operation.Parameters.Tstop", "Tstop", "gtfield", "Tstart")
		}
	}
}

func (v *Validator) validateEntityRulePatterns(sl validator.StructLevel, r EditRequest) {
	if r.CorporateEntityPattern == "" && len(r.EntityPattern) > 0 &&
		!match.ValidateEntityPattern(r.EntityPattern, common.GetForbiddenFieldsInEntityPattern(mongo.IdleRuleMongoCollection)) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.CorporateAlarmPattern != "" || len(r.AlarmPattern) > 0 {
		sl.ReportError(r.Type, "AlarmPattern", "AlarmPattern", "must_be_empty", "")
	}

	if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
	}
}

func (v *Validator) validateAlarmRulePatterns(sl validator.StructLevel, r EditRequest) {
	if r.CorporateEntityPattern == "" && len(r.EntityPattern) > 0 &&
		!match.ValidateEntityPattern(r.EntityPattern, common.GetForbiddenFieldsInEntityPattern(mongo.IdleRuleMongoCollection)) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.CorporateAlarmPattern == "" && len(r.AlarmPattern) > 0 &&
		!match.ValidateAlarmPattern(r.AlarmPattern,
			common.GetForbiddenFieldsInAlarmPattern(mongo.IdleRuleMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.IdleRuleMongoCollection),
		) {
		sl.ReportError(r.EntityPattern, "AlarmPattern", "AlarmPattern", "alarm_pattern", "")
	}

	if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" &&
		len(r.AlarmPattern) == 0 && r.CorporateAlarmPattern == "" {
		sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required_or", "EntityPattern")
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_or", "AlarmPattern")
	}
}
