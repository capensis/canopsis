package eventfilter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type eventfilterValidator struct {
	dbClient mongo.DbClient
}

func (v *eventfilterValidator) ValidateEventFilter(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(eventfilter.Rule)
	v.validateFields(r, sl)

	if r.ID != "" {
		err := v.dbClient.Collection(mongo.EventFilterRulesMongoCollection).FindOne(ctx, bson.M{"_id": r.ID}).Err()
		if err == nil {
			sl.ReportError("_id", "ID", "ID", "unique", "")
		} else if err != mongodriver.ErrNoDocuments {
			panic(err)
		}
	}
}

func (v *eventfilterValidator) validateFields(p eventfilter.Rule, sl validator.StructLevel) {
	if !p.Patterns.IsValid() {
		sl.ReportError(p.Patterns, "patterns", "Patterns", "eventfilter_patterns_invalid", "")
	}

	if p.Type == eventfilter.RuleTypeChangeEntity &&
		p.Config.Component == "" &&
		p.Config.Resource == "" &&
		p.Config.Connector == "" &&
		p.Config.ConnectorName == "" {
		sl.ReportError(p.Config, "config", "Config", "required", "")
	}

	if p.Type == eventfilter.RuleTypeEnrichment {
		if len(p.Config.Actions) == 0 {
			sl.ReportError(p.Config.Actions, "actions", "Actions", "required", "")
		}

		validOutcome := []string{
			eventfilter.OutcomePass,
			eventfilter.OutcomeDrop,
			eventfilter.OutcomeBreak,
		}
		if p.Config.OnSuccess == "" {
			sl.ReportError(p.Config.OnSuccess, "on_success", "OnSuccess", "required_if", "Type enrichment")
		} else {
			switch p.Config.OnSuccess {
			case eventfilter.OutcomePass, eventfilter.OutcomeDrop, eventfilter.OutcomeBreak:
			default:
				sl.ReportError(p.Config.OnSuccess, "on_success", "OnSuccess", "oneof", strings.Join(validOutcome, " "))
			}
		}

		if p.Config.OnFailure == "" {
			sl.ReportError(p.Config.OnFailure, "on_failure", "OnFailure", "required_if", "Type enrichment")
		} else {
			switch p.Config.OnFailure {
			case eventfilter.OutcomePass, eventfilter.OutcomeDrop, eventfilter.OutcomeBreak:
			default:
				sl.ReportError(p.Config.OnFailure, "on_failure", "OnFailure", "oneof", strings.Join(validOutcome, " "))
			}
		}
	}
}

func NewValidator(dbClient mongo.DbClient) *eventfilterValidator {
	return &eventfilterValidator{dbClient: dbClient}
}
