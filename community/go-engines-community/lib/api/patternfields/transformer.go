package patternfields

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrNotExistCorporateAlarmPattern          = common.NewValidationError("corporate_alarm_pattern", "CorporateAlarmPattern doesn't exist.")
	ErrNotExistCorporateEntityPattern         = common.NewValidationError("corporate_entity_pattern", "CorporateEntityPattern doesn't exist.")
	ErrNotExistCorporateTotalEntityPattern    = common.NewValidationError("corporate_total_entity_pattern", "CorporateTotalEntityPattern doesn't exist.")
	ErrNotExistCorporatePbehaviorPattern      = common.NewValidationError("corporate_pbehavior_pattern", "CorporatePbehaviorPattern doesn't exist.")
	ErrNotExistCorporateWeatherServicePattern = common.NewValidationError("corporate_weather_service_pattern", "CorporateWeatherServicePattern doesn't exist.")
)

type Transformer interface {
	TransformAlarmRequest(ctx context.Context, r AlarmRequest) (AlarmRequest, error)
	TransformEntityRequest(ctx context.Context, r EntityRequest) (EntityRequest, error)
	TransformPbehaviorRequest(ctx context.Context, r PbehaviorRequest) (PbehaviorRequest, error)
	TransformWeatherServiceRequest(ctx context.Context, r WeatherServiceRequest) (WeatherServiceRequest, error)
	TransformTotalEntityRequest(ctx context.Context, r TotalEntityRequest) (TotalEntityRequest, error)
}

func NewTransformer(client mongo.DbClient) Transformer {
	return &baseTransformer{
		patternCollection:             client.Collection(mongo.PatternMongoCollection),
		entityInfosPropertyCollection: client.Collection(mongo.EntityInfosPropertyCollection),
	}
}

type baseTransformer struct {
	patternCollection             mongo.DbCollection
	entityInfosPropertyCollection mongo.DbCollection
}

func (t *baseTransformer) TransformAlarmRequest(ctx context.Context, r AlarmRequest) (AlarmRequest, error) {
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

func (t *baseTransformer) TransformEntityRequest(ctx context.Context, r EntityRequest) (EntityRequest, error) {
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

func (t *baseTransformer) TransformTotalEntityRequest(ctx context.Context, r TotalEntityRequest) (TotalEntityRequest, error) {
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

func (t *baseTransformer) TransformPbehaviorRequest(ctx context.Context, r PbehaviorRequest) (PbehaviorRequest, error) {
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

func (t *baseTransformer) TransformWeatherServiceRequest(ctx context.Context, r WeatherServiceRequest) (WeatherServiceRequest, error) {
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

func (t *baseTransformer) getCommonFilter(isPrivate bool, user string) bson.M {
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

func (t *baseTransformer) transformAliases(ctx context.Context, origPattern pattern.Entity, errPrefix string) (pattern.Entity, []string, error) {
	mutPattern := make([][]pattern.FieldCondition, len(origPattern))
	aliases := make([]string, 0)
	for i := range origPattern {
		mutPattern[i] = slices.Clone(origPattern[i])
		for _, g := range mutPattern[i] {
			if g.Alias != "" {
				aliases = append(aliases, g.Alias)
			}
		}
	}

	if len(aliases) == 0 {
		return origPattern, nil, nil
	}

	cursor, err := t.entityInfosPropertyCollection.Find(
		ctx,
		bson.M{"alias": bson.M{"$in": aliases}},
		options.Find().SetProjection(bson.M{"alias": 1, "name": 1}),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot find entity aliases: %w", err)
	}

	defer cursor.Close(ctx)
	infoNamesByAlias := make(map[string]string, len(aliases))
	aliasPropIDs := make([]string, 0, len(aliases))
	for cursor.Next(ctx) {
		prop := struct {
			ID    string `bson:"_id"`
			Alias string `bson:"alias"`
			Name  string `bson:"name"`
		}{}
		err = cursor.Decode(&prop)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot decode entity alias: %w", err)
		}

		infoNamesByAlias[prop.Alias] = prop.Name
		aliasPropIDs = append(aliasPropIDs, prop.ID)
	}

	if err = cursor.Err(); err != nil {
		return nil, nil, fmt.Errorf("cannot fetch entity aliases: %w", err)
	}

	for _, g := range mutPattern {
		for idx, p := range g {
			if p.Alias == "" {
				continue
			}

			n, ok := infoNamesByAlias[p.Alias]
			if !ok {
				return nil, nil, common.NewValidationError(errPrefix+".alias", "Alias doesn't exist.")
			}

			g[idx].Field = "infos." + n
		}
	}

	return mutPattern, aliasPropIDs, nil
}
