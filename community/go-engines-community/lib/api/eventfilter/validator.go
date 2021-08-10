package eventfilter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/neweventfilter"
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
	r := sl.Current().Interface().(EventFilter)
	v.validateFields(r.EventFilterPayload, sl)

	if r.ID != "" {
		err := v.dbClient.Collection(mongo.EventFilterRulesMongoCollection).FindOne(ctx, bson.M{"_id": r.ID}).Err()
		if err == nil {
			sl.ReportError("_id", "ID", "ID", "unique", "")
		} else if err != mongodriver.ErrNoDocuments {
			panic(err)
		}
	}
}

func (v *eventfilterValidator) ValidateEventFilterPayload(sl validator.StructLevel) {
	r := sl.Current().Interface().(EventFilterPayload)
	v.validateFields(r, sl)
}

func (v *eventfilterValidator) validateFields(p EventFilterPayload, sl validator.StructLevel) {
	if p.Patterns != nil && !p.Patterns.IsValid() {
		sl.ReportError(p.Patterns, "patterns", "Patterns", "eventfilter_patterns_invalid", "")
	}

	if p.Type == neweventfilter.RuleTypeChangeEntity &&
		p.Config.Component == "" &&
		p.Config.Resource == "" &&
		p.Config.Connector == "" &&
		p.Config.ConnectorName == "" {
		sl.ReportError(p.Actions, "config", "Config", "required", "")
	}

	if eventfilter.Type(p.Type) == eventfilter.RuleTypeEnrichment {
		validOutcome := []string{
			string(eventfilter.Pass),
			string(eventfilter.Drop),
			string(eventfilter.Break),
		}
		if p.OnSuccess != "" {
			switch eventfilter.Outcome(p.OnSuccess) {
			case eventfilter.Pass, eventfilter.Drop, eventfilter.Break:
			default:
				sl.ReportError(p.OnSuccess, "OnSuccess", "OnSuccess", "oneof", strings.Join(validOutcome, " "))
			}
		}
		if p.OnFailure != "" {
			switch eventfilter.Outcome(p.OnFailure) {
			case eventfilter.Pass, eventfilter.Drop, eventfilter.Break:
			default:
				sl.ReportError(p.OnFailure, "OnFailure", "OnFailure", "oneof", strings.Join(validOutcome, " "))
			}
		}
	}
}

func NewValidator(dbClient mongo.DbClient) *eventfilterValidator {
	return &eventfilterValidator{dbClient: dbClient}
}
