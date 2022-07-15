package eventfilter

import (
	"context"
	"errors"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type eventfilterValidator struct {
	dbClient mongo.DbClient
}

func NewValidator(client mongo.DbClient) *eventfilterValidator {
	return &eventfilterValidator{dbClient: client}
}

func (v *eventfilterValidator) ValidateCreateRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)
	v.validateEventFilter(ctx, sl, r.ID, r.EditRequest)
}

func (v *eventfilterValidator) ValidateUpdateRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(UpdateRequest)
	v.validateEventFilter(ctx, sl, r.ID, r.EditRequest)
}

func (v *eventfilterValidator) ValidateBulkUpdateRequestItem(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(BulkUpdateRequestItem)
	v.validateEventFilter(ctx, sl, r.ID, r.EditRequest)
}

func (v *eventfilterValidator) validateEventFilter(ctx context.Context, sl validator.StructLevel, id string, r EditRequest) {
	if r.Type == eventfilter.RuleTypeChangeEntity &&
		r.Config.Component == "" &&
		r.Config.Resource == "" &&
		r.Config.Connector == "" &&
		r.Config.ConnectorName == "" {
		sl.ReportError(r.Config, "config", "Config", "required", "")
	}

	if r.Type == eventfilter.RuleTypeEnrichment {
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

	if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" && len(r.EventPattern) == 0 {
		if id != "" {
			err := v.dbClient.Collection(mongo.EventFilterRulesMongoCollection).FindOne(
				ctx,
				bson.M{
					"_id":          id,
					"old_patterns": bson.M{"$ne": nil},
				},
			).Err()

			if err == nil {
				return
			} else if !errors.Is(err, mongodriver.ErrNoDocuments) {
				panic(err)
			}
		}

		sl.ReportError(r.EventPattern, "EventPattern", "EventPattern", "required_or", "EntityPattern")
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_or", "EventPattern")
	}
}
