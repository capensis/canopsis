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

func (v *eventfilterValidator) Validate(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EventFilter)
	if r.Patterns != nil && !r.Patterns.IsValid() {
		sl.ReportError(r.Patterns, "patterns", "Patterns", "eventfilter_patterns_invalid", "")
	}
	if eventfilter.Type(r.Type) == eventfilter.RuleTypeEnrichment {
		if len(r.Actions) == 0 {
			sl.ReportError(r.Actions, "actions", "Actions", "required", "")
		}
		validOutcome := []string{
			string(eventfilter.Pass),
			string(eventfilter.Drop),
			string(eventfilter.Break),
		}
		if r.OnSuccess != "" {
			switch eventfilter.Outcome(r.OnSuccess) {
			case eventfilter.Pass, eventfilter.Drop, eventfilter.Break:
			default:
				sl.ReportError(r.OnSuccess, "OnSuccess", "OnSuccess", "oneof", strings.Join(validOutcome, " "))
			}
		}
		if r.OnFailure != "" {
			switch eventfilter.Outcome(r.OnFailure) {
			case eventfilter.Pass, eventfilter.Drop, eventfilter.Break:
			default:
				sl.ReportError(r.OnFailure, "OnFailure", "OnFailure", "oneof", strings.Join(validOutcome, " "))
			}
		}
	}

	if r.ID != "" {
		err := v.dbClient.Collection(mongo.EventFilterRulesMongoCollection).FindOne(ctx, bson.M{"_id": r.ID}).Err()
		if err == nil {
			sl.ReportError("_id", "ID", "ID", "unique", "")
		} else if err != mongodriver.ErrNoDocuments {
			panic(err)
		}
	}

}

func NewValidator(dbClient mongo.DbClient) *eventfilterValidator {
	return &eventfilterValidator{dbClient: dbClient}
}
