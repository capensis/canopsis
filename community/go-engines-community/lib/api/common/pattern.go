package common

import (
	"context"
	"errors"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrNotExistCorporateAlarmPattern          = NewValidationError("corporate_alarm_pattern", "CorporateAlarmPattern doesn't exist.")
	ErrNotExistCorporateEntityPattern         = NewValidationError("corporate_entity_pattern", "CorporateEntityPattern doesn't exist.")
	ErrNotExistCorporateTotalEntityPattern    = NewValidationError("corporate_total_entity_pattern", "CorporateTotalEntityPattern doesn't exist.")
	ErrNotExistCorporatePbehaviorPattern      = NewValidationError("corporate_pbehavior_pattern", "CorporatePbehaviorPattern doesn't exist.")
	ErrNotExistCorporateWeatherServicePattern = NewValidationError("corporate_weather_service_pattern", "CorporateWeatherServicePattern doesn't exist.")
)

type AlarmPatternFieldsRequest struct {
	AlarmPattern          pattern.Alarm `json:"alarm_pattern" binding:"alarm_pattern"`
	CorporateAlarmPattern string        `json:"corporate_alarm_pattern"`

	CorporatePattern savedpattern.SavedPattern `json:"-"`
	IsPrivate        bool                      `json:"-"`
	User             string                    `json:"-"`
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
	Aliases          []string                  `json:"-"`
	IsPrivate        bool                      `json:"-"`
	User             string                    `json:"-"`
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
	IsPrivate        bool                      `json:"-"`
	User             string                    `json:"-"`
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

type WeatherServicePatternFieldsRequest struct {
	WeatherServicePattern          pattern.WeatherServicePattern `json:"weather_service_pattern" binding:"weather_service_pattern"`
	CorporateWeatherServicePattern string                        `json:"corporate_weather_service_pattern"`

	CorporatePattern savedpattern.SavedPattern `json:"-"`
	IsPrivate        bool                      `json:"-"`
	User             string                    `json:"-"`
}

func (r WeatherServicePatternFieldsRequest) ToModel() savedpattern.WeatherServicePatternFields {
	if r.CorporatePattern.ID == "" {
		return savedpattern.WeatherServicePatternFields{
			WeatherServicePattern: r.WeatherServicePattern,
		}
	}

	return savedpattern.WeatherServicePatternFields{
		WeatherServicePattern:               r.CorporatePattern.WeatherServicePattern,
		CorporateWeatherServicePattern:      r.CorporatePattern.ID,
		CorporateWeatherServicePatternTitle: r.CorporatePattern.Title,
	}
}

type PatternFieldsTransformer interface {
	TransformAlarmPatternFieldsRequest(ctx context.Context, r AlarmPatternFieldsRequest) (AlarmPatternFieldsRequest, error)
	TransformEntityPatternFieldsRequest(ctx context.Context, r EntityPatternFieldsRequest) (EntityPatternFieldsRequest, error)
	TransformPbehaviorPatternFieldsRequest(ctx context.Context, r PbehaviorPatternFieldsRequest) (PbehaviorPatternFieldsRequest, error)
	TransformWeatherServicePatternFieldsRequest(ctx context.Context, r WeatherServicePatternFieldsRequest) (WeatherServicePatternFieldsRequest, error)
	TransformTotalEntityPatternFieldsRequest(ctx context.Context, r TotalEntityPatternFieldsRequest) (TotalEntityPatternFieldsRequest, error)
}

func NewPatternFieldsTransformer(client mongo.DbClient) PatternFieldsTransformer {
	return &basePatternFieldsTransformer{
		patternCollection:             client.Collection(mongo.PatternMongoCollection),
		entityInfosPropertyCollection: client.Collection(mongo.EntityInfosPropertyCollection),
	}
}

type basePatternFieldsTransformer struct {
	patternCollection             mongo.DbCollection
	entityInfosPropertyCollection mongo.DbCollection
}

func (t *basePatternFieldsTransformer) getCommonFilter(isPrivate bool, user string) bson.M {
	if isPrivate {
		return bson.M{
			"$or": []bson.M{
				{"is_corporate": true},
				{"author": user},
			},
		}
	}

	return bson.M{
		"is_corporate": true,
	}
}

func (t *basePatternFieldsTransformer) TransformAlarmPatternFieldsRequest(ctx context.Context, r AlarmPatternFieldsRequest) (AlarmPatternFieldsRequest, error) {
	if r.CorporateAlarmPattern != "" {
		filter := t.getCommonFilter(r.IsPrivate, r.User)
		filter["_id"] = r.CorporateAlarmPattern
		filter["type"] = savedpattern.TypeAlarm

		err := t.patternCollection.FindOne(ctx, filter).Decode(&r.CorporatePattern)
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
		filter := t.getCommonFilter(r.IsPrivate, r.User)
		filter["_id"] = r.CorporateEntityPattern
		filter["type"] = savedpattern.TypeEntity

		err := t.patternCollection.FindOne(ctx, filter).Decode(&r.CorporatePattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return r, ErrNotExistCorporateEntityPattern
			}

			return r, err
		}
	}

	entityPattern := r.EntityPattern
	if r.CorporateEntityPattern != "" {
		entityPattern = r.CorporatePattern.EntityPattern
	}

	var err error

	r.EntityPattern, r.Aliases, err = t.transformAliases(ctx, entityPattern, "entity_pattern")
	return r, err
}

func (t *basePatternFieldsTransformer) TransformTotalEntityPatternFieldsRequest(ctx context.Context, r TotalEntityPatternFieldsRequest) (TotalEntityPatternFieldsRequest, error) {
	if r.CorporateTotalEntityPattern != "" {
		err := t.patternCollection.FindOne(ctx, bson.M{"_id": r.CorporateTotalEntityPattern, "type": savedpattern.TypeEntity}).Decode(&r.CorporatePattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return r, ErrNotExistCorporateTotalEntityPattern
			}

			return r, err
		}
	}

	entityPattern := r.TotalEntityPattern
	if r.CorporateTotalEntityPattern != "" {
		entityPattern = r.CorporatePattern.EntityPattern
	}

	var err error

	r.TotalEntityPattern, r.Aliases, err = t.transformAliases(ctx, entityPattern, "total_entity_pattern")
	return r, err
}

func (t *basePatternFieldsTransformer) transformAliases(ctx context.Context, origPattern pattern.Entity, errPrefix string) (pattern.Entity, []string, error) {
	var uniqueAliasesMap = make(map[string]string)
	var uniqueAliases []string

	mutPattern := make([][]pattern.FieldCondition, len(origPattern))
	for i := range origPattern {
		mutPattern[i] = slices.Clone(origPattern[i])
	}

	for _, p := range mutPattern {
		for idx, subP := range p {
			if subP.Alias == "" {
				continue
			}

			var doc struct {
				ID   string `bson:"_id"`
				Name string `bson:"name"`
			}

			if k, ok := uniqueAliasesMap[subP.Alias]; ok {
				p[idx].Field = "infos." + k
			} else {
				err := t.entityInfosPropertyCollection.FindOne(
					ctx,
					bson.M{"alias": subP.Alias},
					options.FindOne().SetProjection(bson.M{"name": 1}),
				).Decode(&doc)
				if err != nil {
					if errors.Is(err, mongodriver.ErrNoDocuments) {
						return nil, nil, NewValidationError(errPrefix+".alias", "Alias doesn't exist.")
					}

					return nil, nil, err
				}

				p[idx].Field = "infos." + doc.Name
				uniqueAliases = append(uniqueAliases, doc.ID)
				uniqueAliasesMap[subP.Alias] = doc.Name
			}
		}
	}

	return mutPattern, uniqueAliases, nil
}

func (t *basePatternFieldsTransformer) TransformPbehaviorPatternFieldsRequest(ctx context.Context, r PbehaviorPatternFieldsRequest) (PbehaviorPatternFieldsRequest, error) {
	if r.CorporatePbehaviorPattern != "" {
		filter := t.getCommonFilter(r.IsPrivate, r.User)
		filter["_id"] = r.CorporatePbehaviorPattern
		filter["type"] = savedpattern.TypePbehavior

		err := t.patternCollection.FindOne(ctx, filter).Decode(&r.CorporatePattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return r, ErrNotExistCorporatePbehaviorPattern
			}

			return r, err
		}
	}

	return r, nil
}

func (t *basePatternFieldsTransformer) TransformWeatherServicePatternFieldsRequest(ctx context.Context, r WeatherServicePatternFieldsRequest) (WeatherServicePatternFieldsRequest, error) {
	if r.CorporateWeatherServicePattern != "" {
		filter := t.getCommonFilter(r.IsPrivate, r.User)
		filter["_id"] = r.CorporateWeatherServicePattern
		filter["type"] = savedpattern.TypeWeatherService

		err := t.patternCollection.FindOne(ctx, filter).Decode(&r.CorporatePattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return r, ErrNotExistCorporateWeatherServicePattern
			}

			return r, err
		}
	}

	return r, nil
}

type TotalEntityPatternFieldsRequest struct {
	TotalEntityPattern          pattern.Entity `json:"total_entity_pattern" binding:"entity_pattern"`
	CorporateTotalEntityPattern string         `json:"corporate_total_entity_pattern"`

	CorporatePattern savedpattern.SavedPattern `json:"-"`
	Aliases          []string                  `json:"-"`
}

func (r TotalEntityPatternFieldsRequest) ToModel() correlation.TotalEntityPatternFields {
	if r.CorporatePattern.ID == "" {
		return correlation.TotalEntityPatternFields{
			TotalEntityPattern: r.TotalEntityPattern,
		}
	}

	return correlation.TotalEntityPatternFields{
		TotalEntityPattern:               r.CorporatePattern.EntityPattern,
		CorporateTotalEntityPattern:      r.CorporatePattern.ID,
		CorporateTotalEntityPatternTitle: r.CorporatePattern.Title,
	}
}

func (r TotalEntityPatternFieldsRequest) ToModelWithoutFields(forbiddenFields []string) correlation.TotalEntityPatternFields {
	if r.CorporatePattern.ID == "" {
		return correlation.TotalEntityPatternFields{
			TotalEntityPattern: r.TotalEntityPattern,
		}
	}

	return correlation.TotalEntityPatternFields{
		TotalEntityPattern:               r.CorporatePattern.EntityPattern.RemoveFields(forbiddenFields),
		CorporateTotalEntityPattern:      r.CorporatePattern.ID,
		CorporateTotalEntityPatternTitle: r.CorporatePattern.Title,
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

	return match.ValidateAlarmPattern(p, nil, nil)
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

	return match.ValidateEventPattern(p)
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

	return match.ValidateEntityPattern(p, nil)
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

	return match.ValidatePbehaviorInfoPattern(p)
}

func ValidateWeatherServicePattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}
	p, ok := i.(pattern.WeatherServicePattern)
	if !ok {
		return false
	}

	return match.ValidateWeatherServicePattern(p)
}

func GetForbiddenFieldsInEntityPattern(collection string) []string {
	switch collection {
	case mongo.StateSettingsMongoCollection:
		return []string{"last_event_date", "component", "component_infos"}
	case mongo.EntityMongoCollection:
		return []string{"last_event_date", "connector", "component_infos"}
	case mongo.PbehaviorMongoCollection,
		mongo.IdleRuleMongoCollection,
		mongo.DynamicInfosRulesMongoCollection,
		mongo.MetaAlarmRulesMongoCollection,
		mongo.FlappingRuleMongoCollection,
		mongo.ResolveRuleMongoCollection,
		mongo.ScenarioCollection,
		mongo.InstructionMongoCollection,
		mongo.KpiFilterMongoCollection,
		mongo.DeclareTicketRuleCollection,
		mongo.LinkRuleMongoCollection,
		mongo.AlarmTagCollection:
		return []string{"last_event_date"}
	default:
		return nil
	}
}

func GetForbiddenFieldsInAlarmPattern(collection string) []string {
	switch collection {
	case mongo.IdleRuleMongoCollection,
		mongo.MetaAlarmRulesMongoCollection,
		mongo.FlappingRuleMongoCollection,
		mongo.ResolveRuleMongoCollection,
		mongo.ScenarioCollection,
		mongo.InstructionMongoCollection,
		mongo.DeclareTicketRuleCollection,
		mongo.LinkRuleMongoCollection:
		return []string{"v.last_event_date", "v.last_update_date", "v.resolved"}
	case mongo.DynamicInfosRulesMongoCollection:
		return []string{"v.last_event_date", "v.last_update_date", "v.resolved", "v.duration", "v.infos"}
	case mongo.AlarmTagCollection:
		return []string{"v.last_event_date", "v.last_update_date", "v.resolved", "v.duration", "tags"}
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
		mongo.ScenarioCollection,
		mongo.InstructionMongoCollection,
		mongo.DeclareTicketRuleCollection,
		mongo.LinkRuleMongoCollection,
		mongo.AlarmTagCollection:
		return []string{"v.creation_date", "v.ack.t", "v.activation_date"}
	default:
		return nil
	}
}
