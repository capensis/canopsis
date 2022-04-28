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
		err := t.patternCollection.FindOne(ctx, bson.M{"_id": r.CorporateAlarmPattern, "type": savedpattern.TypeAlarm}).Decode(&r.CorporatePattern)
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
		err := t.patternCollection.FindOne(ctx, bson.M{"_id": r.CorporateEntityPattern, "type": savedpattern.TypeEntity}).Decode(&r.CorporatePattern)
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
		err := t.patternCollection.FindOne(ctx, bson.M{"_id": r.CorporatePbehaviorPattern, "type": savedpattern.TypePbehavior}).Decode(&r.CorporatePattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return r, ErrNotExistCorporatePbehaviorPattern
			}

			return r, err
		}
	}

	return r, nil
}

func ValidateAlarmPatternFieldsRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(AlarmPatternFieldsRequest)

	if r.CorporateAlarmPattern != "" && len(r.AlarmPattern) > 0 {
		sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required_not_both", "CorporateAlarmPattern")
	}
}

func ValidateEntityPatternFieldsRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EntityPatternFieldsRequest)

	if r.CorporateEntityPattern != "" && len(r.EntityPattern) > 0 {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_not_both", "CorporateEntityPattern")
	}
}

func ValidatePbehaviorPatternFieldsRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(PbehaviorPatternFieldsRequest)

	if r.CorporatePbehaviorPattern != "" && len(r.PbehaviorPattern) > 0 {
		sl.ReportError(r.PbehaviorPattern, "PbehaviorPattern", "PbehaviorPattern", "required_not_both", "CorporatePbehaviorPattern")
	}
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

	return p.Validate()
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

	return p.Validate()
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
