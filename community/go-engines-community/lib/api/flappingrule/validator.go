package flappingrule

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Validator struct {
	dbClient mongo.DbClient
}

func NewValidator(client mongo.DbClient) *Validator {
	return &Validator{dbClient: client}
}

func (v *Validator) ValidateCreateRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)
	v.validatePatterns(ctx, sl, r.ID, r.EditRequest)
}

func (v *Validator) ValidateUpdateRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(UpdateRequest)
	v.validatePatterns(ctx, sl, r.ID, r.EditRequest)
}

func (v *Validator) validatePatterns(ctx context.Context, sl validator.StructLevel, id string, r EditRequest) {
	if r.CorporateEntityPattern == "" && len(r.EntityPattern) > 0 &&
		!r.EntityPattern.Validate(common.GetForbiddenFieldsInEntityPattern(mongo.FlappingRuleMongoCollection)) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.CorporateAlarmPattern == "" && len(r.AlarmPattern) > 0 &&
		!r.AlarmPattern.Validate(
			common.GetForbiddenFieldsInAlarmPattern(mongo.FlappingRuleMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.FlappingRuleMongoCollection),
		) {
		sl.ReportError(r.EntityPattern, "AlarmPattern", "AlarmPattern", "alarm_pattern", "")
	}

	if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" &&
		len(r.AlarmPattern) == 0 && r.CorporateAlarmPattern == "" {

		if id != "" {
			err := v.dbClient.Collection(mongo.FlappingRuleMongoCollection).FindOne(
				ctx,
				bson.M{
					"_id": id,
					"$or": bson.A{
						bson.M{"old_entity_patterns": bson.M{"$ne": nil}},
						bson.M{"old_alarm_patterns": bson.M{"$ne": nil}},
					},
				},
			).Err()

			if err == nil {
				return
			} else if !errors.Is(err, mongodriver.ErrNoDocuments) {
				panic(err)
			}
		}

		sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required_or", "EntityPattern")
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_or", "AlarmPattern")
	}
}
