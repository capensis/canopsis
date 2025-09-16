package eventfilter

import (
	"context"
	"strings"

	apiexternaldata "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/externaldatatable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"github.com/teambition/rrule-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Validator struct {
	dbClient mongo.DbClient
}

func NewValidator(client mongo.DbClient) *Validator {
	return &Validator{dbClient: client}
}

func (v *Validator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)
	switch r.Type {
	case eventfilter.RuleTypeChangeEntity:
		if r.Config.Component == "" &&
			r.Config.Resource == "" &&
			r.Config.Connector == "" &&
			r.Config.ConnectorName == "" {
			sl.ReportError(r.Config, "config", "Config", "required", "")
		}
	case eventfilter.RuleTypeEnrichment:
		if len(r.Config.Actions) == 0 {
			sl.ReportError(r.Config.Actions, "actions", "Actions", "required", "")
		}

		validOutcome := []string{
			eventfilter.OutcomePass,
			eventfilter.OutcomeDrop,
			eventfilter.OutcomeBreak,
		}
		if r.Config.OnSuccess == "" {
			sl.ReportError(r.Config.OnSuccess, "on_success", "OnSuccess", "required_if", "Type enrichment")
		} else {
			switch r.Config.OnSuccess {
			case eventfilter.OutcomePass, eventfilter.OutcomeDrop, eventfilter.OutcomeBreak:
			default:
				sl.ReportError(r.Config.OnSuccess, "on_success", "OnSuccess", "oneof", strings.Join(validOutcome, " "))
			}
		}

		if r.Config.OnFailure == "" {
			sl.ReportError(r.Config.OnFailure, "on_failure", "OnFailure", "required_if", "Type enrichment")
		} else {
			switch r.Config.OnFailure {
			case eventfilter.OutcomePass, eventfilter.OutcomeDrop, eventfilter.OutcomeBreak:
			default:
				sl.ReportError(r.Config.OnFailure, "on_failure", "OnFailure", "oneof", strings.Join(validOutcome, " "))
			}
		}
	}

	switch r.Type {
	case eventfilter.RuleTypeChangeEntity:
		if len(r.EventPattern) == 0 {
			sl.ReportError(r.EventPattern, "EventPattern", "EventPattern", "required", "")
		}

		if len(r.EntityPattern) > 0 {
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "must_be_empty", "")
		}
	default:
		if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" && len(r.EventPattern) == 0 {
			sl.ReportError(r.EventPattern, "EventPattern", "EventPattern", "required_or", "EntityPattern")
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_or", "EventPattern")
		}
	}

	if r.Start == nil && r.Stop != nil {
		sl.ReportError(r.Start, "Start", "Start", "required_with", "Stop")
	}

	if r.Stop == nil && r.Start != nil {
		sl.ReportError(r.Stop, "Stop", "Stop", "required_with", "Start")
	}

	if r.Stop != nil && r.Start != nil && r.Start.Unix() >= r.Stop.Unix() {
		sl.ReportError(r.Stop, "Stop", "Stop", "gtfield", "Start")
	}

	if r.RRule != "" && r.Stop == nil && r.Start == nil {
		sl.ReportError(r.Start, "Start", "Start", "required_with", "RRule")
		sl.ReportError(r.Stop, "Stop", "Stop", "required_with", "RRule")
	}

	if r.RRule != "" && !v.checkRrule(r.RRule) {
		sl.ReportError(r.RRule, "RRule", "RRule", "rrule", "")
	}

	ok, err := v.checkExceptions(ctx, r.Exceptions)
	if err != nil {
		panic(err)
	}
	if !ok {
		sl.ReportError(r.Exceptions, "Exceptions", "Exceptions", "not_exist", "")
	}

	apiexternaldata.ValidateRefParameters(sl, r.ExternalData, []string{externaldata.RefTypeAPI, externaldata.RefTypeTable})
}

func (v *Validator) checkRrule(r string) bool {
	_, err := rrule.StrToROption(r)
	return err == nil
}

func (v *Validator) checkExceptions(ctx context.Context, exceptions []string) (bool, error) {
	if len(exceptions) == 0 {
		return true, nil
	}
	count, err := v.dbClient.
		Collection(mongo.PbehaviorExceptionMongoCollection).
		CountDocuments(ctx, bson.M{"_id": bson.M{"$in": exceptions}})
	if err != nil {
		return false, err
	}

	return count == int64(len(exceptions)), nil
}
