package common

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNotExistCorporateAlarmPattern     = ValidationError{field: "corporate_alarm_pattern", err: errors.New("CorporateAlarmPattern doesn't exist.")}
	ErrNotExistCorporateEntityPattern    = ValidationError{field: "corporate_entity_pattern", err: errors.New("CorporateEntityPattern doesn't exist.")}
	ErrNotExistCorporatePbehaviorPattern = ValidationError{field: "corporate_pbehavior_pattern", err: errors.New("CorporatePbehaviorPattern doesn't exist.")}
)

type AlarmPatternFieldsRequest struct {
	AlarmPattern          pattern.Alarm `json:"alarm_pattern" binding:"alarm_pattern"`
	CorporateAlarmPattern string        `json:"corporate_alarm_pattern"`

	CorporatePattern savedpattern.SavedPattern `json:"-"`
}

func (r AlarmPatternFieldsRequest) ToModel() savedpattern.AlarmPatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.AlarmPatternFields{
			AlarmPattern: r.AlarmPattern,
		}
	}

	return savedpattern.AlarmPatternFields{
		AlarmPattern:               r.CorporatePattern.AlarmPattern,
		CorporateAlarmPattern:      r.CorporatePattern.ID,
		CorporateAlarmPatternTitle: r.CorporatePattern.Title,
	}
}

func (r AlarmPatternFieldsRequest) ToModelWithoutFields(forbiddenFields, onlyTimeAbsoluteFields []string) savedpattern.AlarmPatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.AlarmPatternFields{
			AlarmPattern: r.AlarmPattern,
		}
	}

	return savedpattern.AlarmPatternFields{
		AlarmPattern:               r.CorporatePattern.AlarmPattern.RemoveFields(forbiddenFields, onlyTimeAbsoluteFields),
		CorporateAlarmPattern:      r.CorporatePattern.ID,
		CorporateAlarmPatternTitle: r.CorporatePattern.Title,
	}
}

type EntityPatternFieldsRequest struct {
	EntityPattern          pattern.Entity `json:"entity_pattern" binding:"entity_pattern"`
	CorporateEntityPattern string         `json:"corporate_entity_pattern"`

	CorporatePattern savedpattern.SavedPattern `json:"-"`
}

func (r EntityPatternFieldsRequest) ToModel() savedpattern.EntityPatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.EntityPatternFields{
			EntityPattern: r.EntityPattern,
		}
	}

	return savedpattern.EntityPatternFields{
		EntityPattern:               r.CorporatePattern.EntityPattern,
		CorporateEntityPattern:      r.CorporatePattern.ID,
		CorporateEntityPatternTitle: r.CorporatePattern.Title,
	}
}

func (r EntityPatternFieldsRequest) ToModelWithoutFields(forbiddenFields []string) savedpattern.EntityPatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.EntityPatternFields{
			EntityPattern: r.EntityPattern,
		}
	}

	return savedpattern.EntityPatternFields{
		EntityPattern:               r.CorporatePattern.EntityPattern.RemoveFields(forbiddenFields),
		CorporateEntityPattern:      r.CorporatePattern.ID,
		CorporateEntityPatternTitle: r.CorporatePattern.Title,
	}
}

type PbehaviorPatternFieldsRequest struct {
	PbehaviorPattern          pattern.PbehaviorInfo `json:"pbehavior_pattern" binding:"pbehavior_pattern"`
	CorporatePbehaviorPattern string                `json:"corporate_pbehavior_pattern"`

	CorporatePattern savedpattern.SavedPattern `json:"-"`
}

func (r PbehaviorPatternFieldsRequest) ToModel() savedpattern.PbehaviorPatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.PbehaviorPatternFields{
			PbehaviorPattern: r.PbehaviorPattern,
		}
	}

	return savedpattern.PbehaviorPatternFields{
		PbehaviorPattern:               r.CorporatePattern.PbehaviorPattern,
		CorporatePbehaviorPattern:      r.CorporatePattern.ID,
		CorporatePbehaviorPatternTitle: r.CorporatePattern.Title,
	}
}

type PatternFieldsTransformer interface {
	TransformAlarmPatternFieldsRequest(ctx context.Context, r AlarmPatternFieldsRequest) (AlarmPatternFieldsRequest, error)
	TransformEntityPatternFieldsRequest(ctx context.Context, r EntityPatternFieldsRequest) (EntityPatternFieldsRequest, error)
	TransformPbehaviorPatternFieldsRequest(ctx context.Context, r PbehaviorPatternFieldsRequest) (PbehaviorPatternFieldsRequest, error)
}

func NewPatternFieldsTransformer(client mongo.DbClient) PatternFieldsTransformer {
	return &basePatternFieldsTransformer{
		patternCollection: client.Collection(mongo.PatternMongoCollection),
	}
}

type basePatternFieldsTransformer struct {
	patternCollection mongo.DbCollection
}

func (t *basePatternFieldsTransformer) TransformAlarmPatternFieldsRequest(ctx context.Context, r AlarmPatternFieldsRequest) (AlarmPatternFieldsRequest, error) {
	if r.CorporateAlarmPattern != "" {
		err := t.patternCollection.FindOne(ctx, bson.M{
			"_id":          r.CorporateAlarmPattern,
			"type":         savedpattern.TypeAlarm,
			"is_corporate": true,
		}).Decode(&r.CorporatePattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return r, ErrNotExistCorporateAlarmPattern
			}

			return r, err
		}
	}

	return r, nil
}

func (t *basePatternFieldsTransformer) TransformEntityPatternFieldsRequest(ctx context.Context, r EntityPatternFieldsRequest) (EntityPatternFieldsRequest, error) {
	if r.CorporateEntityPattern != "" {
		err := t.patternCollection.FindOne(ctx, bson.M{
			"_id":          r.CorporateEntityPattern,
			"type":         savedpattern.TypeEntity,
			"is_corporate": true,
		}).Decode(&r.CorporatePattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return r, ErrNotExistCorporateEntityPattern
			}

			return r, err
		}
	}

	return r, nil
}

func (t *basePatternFieldsTransformer) TransformPbehaviorPatternFieldsRequest(ctx context.Context, r PbehaviorPatternFieldsRequest) (PbehaviorPatternFieldsRequest, error) {
	if r.CorporatePbehaviorPattern != "" {
		err := t.patternCollection.FindOne(ctx, bson.M{
			"_id":          r.CorporatePbehaviorPattern,
			"type":         savedpattern.TypePbehavior,
			"is_corporate": true,
		}).Decode(&r.CorporatePattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return r, ErrNotExistCorporatePbehaviorPattern
			}

			return r, err
		}
	}

	return r, nil
}

func ValidateAlarmPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}
	p, ok := i.(pattern.Alarm)
	if !ok {
		return false
	}

	return p.Validate(nil, nil)
}

func ValidateEventPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}
	p, ok := i.(pattern.Event)
	if !ok {
		return false
	}

	return p.Validate()
}

func ValidateEntityPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}
	p, ok := i.(pattern.Entity)
	if !ok {
		return false
	}

	return p.Validate(nil)
}

func ValidatePbehaviorPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}
	p, ok := i.(pattern.PbehaviorInfo)
	if !ok {
		return false
	}

	return p.Validate()
}

func GetForbiddenFieldsInEntityPattern(collection string) []string {
	switch collection {
	case mongo.EntityMongoCollection:
		return []string{"last_event_date", "connector", "component_infos"}
	case mongo.PbehaviorMongoCollection,
		mongo.IdleRuleMongoCollection,
		mongo.DynamicInfosRulesMongoCollection,
		mongo.MetaAlarmRulesMongoCollection,
		mongo.FlappingRuleMongoCollection,
		mongo.ResolveRuleMongoCollection,
		mongo.ScenarioMongoCollection,
		mongo.InstructionMongoCollection,
		mongo.KpiFilterMongoCollection,
		mongo.DeclareTicketRuleMongoCollection,
		mongo.LinkRuleMongoCollection:
		return []string{"last_event_date"}
	default:
		return nil
	}
}

func GetForbiddenFieldsInAlarmPattern(collection string) []string {
	switch collection {
	case mongo.IdleRuleMongoCollection,
		mongo.DynamicInfosRulesMongoCollection,
		mongo.MetaAlarmRulesMongoCollection,
		mongo.FlappingRuleMongoCollection,
		mongo.ResolveRuleMongoCollection,
		mongo.ScenarioMongoCollection,
		mongo.InstructionMongoCollection,
		mongo.DeclareTicketRuleMongoCollection,
		mongo.LinkRuleMongoCollection:
		return []string{"v.last_event_date", "v.last_update_date", "v.resolved"}
	default:
		return nil
	}
}

func GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(collection string) []string {
	switch collection {
	case mongo.IdleRuleMongoCollection,
		mongo.DynamicInfosRulesMongoCollection,
		mongo.MetaAlarmRulesMongoCollection,
		mongo.FlappingRuleMongoCollection,
		mongo.ResolveRuleMongoCollection,
		mongo.ScenarioMongoCollection,
		mongo.InstructionMongoCollection,
		mongo.DeclareTicketRuleMongoCollection,
		mongo.LinkRuleMongoCollection:
		return []string{"v.creation_date", "v.ack.t", "v.activation_date"}
	default:
		return nil
	}
}
