package template

import (
	"encoding/json"

	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	enforcer security.Enforcer
}

func NewValidator(enforcer security.Enforcer) *Validator {
	return &Validator{enforcer: enforcer}
}

func (v *Validator) ValidateEditDataRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditDataRequest)
	if r.Type != nil && *r.Type == TypeTestDataEvent && r.Body != nil {
		b, err := json.Marshal(r.Body)
		if err == nil {
			event := types.Event{}
			err = json.Unmarshal(b, &event)
			if err == nil {
				err = event.InjectExtraInfos(b)
			}
		}

		if err != nil {
			sl.ReportError(r.Body, "Body", "Body", "invalid", "")
		}
	}
}

func (v *Validator) ValidateEditTestRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditTestRequest)

	if r.Type != nil {
		switch *r.Type {
		case TypeTestEventFilterRule,
			TypeTestActionScenario:
			if r.Data.Event == "" {
				sl.ReportError(r.Data.Event, "Event", "Data.Event", "required", "")
			}
		case TypeTestLinkRule:
			if r.Data.Alarm == "" && r.Data.Entity == "" {
				sl.ReportError(r.Data.Alarm, "Alarm", "Data.Alarm", "required_or", "Entity")
				sl.ReportError(r.Data.Entity, "Entity", "Data.Entity", "required_or", "Alarm")
			}
		case TypeTestDynamicInfosRule:
			if r.Data.Alarm == "" && r.Data.Event == "" {
				sl.ReportError(r.Data.Alarm, "Alarm", "Data.Alarm", "required_or", "Event")
				sl.ReportError(r.Data.Event, "Event", "Data.Event", "required_or", "Alarm")
			}
		case TypeTestWidget,
			TypeTestInstruction,
			TypeTestJob,
			TypeTestMetaAlarmRule:
			if r.Data.Alarm == "" {
				sl.ReportError(r.Data.Alarm, "Alarm", "Data.Alarm", "required", "")
			}
		case TypeTestDeclareTicketRule:
			if r.Data.Alarm == "" {
				sl.ReportError(r.Data.Alarm, "Alarm", "Data.Alarm", "required", "")
			}

			if len(r.Data.Responses) == 0 {
				sl.ReportError(r.Data.Responses, "Responses", "Data.Responses", "required", "")
			}
		case TypeTestWebhookTokenRule:
			if r.Data.Response == "" {
				sl.ReportError(r.Data.Response, "Response", "Data.Response", "required", "")
			}
		default:
			sl.ReportError(r.Type, "Type", "Type", "invalid", "")
		}
	}

	if r.Data.User != "" {
		ok, err := v.enforcer.Enforce(r.Author, apisecurity.PermAcl, securitymodel.PermissionRead)
		if err != nil {
			panic(err)
		}

		if !ok {
			sl.ReportError(r.Data.User, "User", "Data.User", "unauth", "")
		}
	}
}
