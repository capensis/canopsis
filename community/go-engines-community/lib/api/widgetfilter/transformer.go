package widgetfilter

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type PatternFieldsTransformer interface {
	TransformAlarmPatternFieldsRequest(ctx context.Context, r common.AlarmPatternFieldsRequest, isPrivate bool, user string) (common.AlarmPatternFieldsRequest, error)
	TransformEntityPatternFieldsRequest(ctx context.Context, r common.EntityPatternFieldsRequest, isPrivate bool, user string) (common.EntityPatternFieldsRequest, error)
	TransformPbehaviorPatternFieldsRequest(ctx context.Context, r common.PbehaviorPatternFieldsRequest, isPrivate bool, user string) (common.PbehaviorPatternFieldsRequest, error)
	TransformWeatherServicePatternFieldsRequest(ctx context.Context, r common.WeatherServicePatternFieldsRequest, isPrivate bool, user string) (common.WeatherServicePatternFieldsRequest, error)
}

type patternTransformer struct {
	patternCollection mongo.DbCollection
}

func NewPatternFieldsTransformer(client mongo.DbClient) PatternFieldsTransformer {
	return &patternTransformer{
		patternCollection: client.Collection(mongo.PatternMongoCollection),
	}
}

func (t *patternTransformer) TransformAlarmPatternFieldsRequest(ctx context.Context, r common.AlarmPatternFieldsRequest, isPrivate bool, user string) (common.AlarmPatternFieldsRequest, error) {
	if r.CorporateAlarmPattern != "" {
		filter := t.getCommonFilter(isPrivate, user)
		filter["_id"] = r.CorporateAlarmPattern
		filter["type"] = savedpattern.TypeAlarm
		err := t.patternCollection.FindOne(ctx, filter).Decode(&r.CorporatePattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return r, common.ErrNotExistCorporateAlarmPattern
			}

			return r, err
		}
	}

	return r, nil
}

func (t *patternTransformer) TransformEntityPatternFieldsRequest(ctx context.Context, r common.EntityPatternFieldsRequest, isPrivate bool, user string) (common.EntityPatternFieldsRequest, error) {
	if r.CorporateEntityPattern != "" {
		filter := t.getCommonFilter(isPrivate, user)
		filter["_id"] = r.CorporateEntityPattern
		filter["type"] = savedpattern.TypeEntity
		err := t.patternCollection.FindOne(ctx, filter).Decode(&r.CorporatePattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return r, common.ErrNotExistCorporateEntityPattern
			}

			return r, err
		}
	}

	return r, nil
}

func (t *patternTransformer) TransformPbehaviorPatternFieldsRequest(ctx context.Context, r common.PbehaviorPatternFieldsRequest, isPrivate bool, user string) (common.PbehaviorPatternFieldsRequest, error) {
	if r.CorporatePbehaviorPattern != "" {
		filter := t.getCommonFilter(isPrivate, user)
		filter["_id"] = r.CorporatePbehaviorPattern
		filter["type"] = savedpattern.TypePbehavior
		err := t.patternCollection.FindOne(ctx, filter).Decode(&r.CorporatePattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return r, common.ErrNotExistCorporatePbehaviorPattern
			}

			return r, err
		}
	}

	return r, nil
}

func (t *patternTransformer) TransformWeatherServicePatternFieldsRequest(ctx context.Context, r common.WeatherServicePatternFieldsRequest, isPrivate bool, user string) (common.WeatherServicePatternFieldsRequest, error) {
	if r.CorporateWeatherServicePattern != "" {
		filter := t.getCommonFilter(isPrivate, user)
		filter["_id"] = r.CorporateWeatherServicePattern
		filter["type"] = savedpattern.TypeWeatherService
		err := t.patternCollection.FindOne(ctx, filter).Decode(&r.CorporatePattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return r, common.ErrNotExistCorporateWeatherServicePattern
			}

			return r, err
		}
	}

	return r, nil
}

func (t *patternTransformer) getCommonFilter(isPrivate bool, user string) bson.M {
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
