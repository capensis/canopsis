package entityservice

import (
	"context"
	"errors"
	"reflect"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	GetOneBy(ctx context.Context, id string) (*Response, error)
	GetDependencies(ctx context.Context, apiKey string, r ContextGraphRequest) (*ContextGraphAggregationResult, error)
	GetImpacts(ctx context.Context, apiKey string, r ContextGraphRequest) (*ContextGraphAggregationResult, error)
	Create(ctx context.Context, request CreateRequest) (*Response, error)
	Update(ctx context.Context, request UpdateRequest) (*Response, ServiceChanges, error)
	Delete(ctx context.Context, id string) (bool, *types.Alarm, error)
}

type ServiceChanges struct {
	IsPatternChanged bool
	IsToggled        bool
}

type store struct {
	dbClient                  mongo.DbClient
	dbCollection              mongo.DbCollection
	alarmDbCollection         mongo.DbCollection
	resolvedAlarmDbCollection mongo.DbCollection

	queryBuilder *entity.MongoQueryBuilder

	linksFetcher common.LinksFetcher

	logger zerolog.Logger
}

func NewStore(db mongo.DbClient, linksFetcher common.LinksFetcher, logger zerolog.Logger) Store {
	return &store{
		dbClient:                  db,
		dbCollection:              db.Collection(mongo.EntityMongoCollection),
		alarmDbCollection:         db.Collection(mongo.AlarmMongoCollection),
		resolvedAlarmDbCollection: db.Collection(mongo.ResolvedAlarmMongoCollection),

		queryBuilder: entity.NewMongoQueryBuilder(db),

		linksFetcher: linksFetcher,

		logger: logger,
	}
}

func (s *store) GetOneBy(ctx context.Context, id string) (*Response, error) {
	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id, "type": types.EntityTypeService}},
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
	})
	if err != nil {
		return nil, err
	}

	if cursor.Next(ctx) {
		res := Response{}
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, nil
}

func (s *store) GetDependencies(ctx context.Context, apiKey string, r ContextGraphRequest) (*ContextGraphAggregationResult, error) {
	service := types.Entity{}
	err := s.dbCollection.
		FindOne(ctx, bson.M{"_id": r.ID, "type": types.EntityTypeService}).
		Decode(&service)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	now := types.NewCpsTime()
	match := bson.M{"_id": bson.M{"$in": service.Depends}}
	pipeline := s.queryBuilder.CreateTreeOfDepsAggregationPipeline(match, r.Query, r.SortRequest, r.Category, r.Search,
		r.WithFlags, now)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	result := &ContextGraphAggregationResult{}

	if cursor.Next(ctx) {
		err = cursor.Decode(result)
		if err != nil {
			return nil, err
		}
	}

	err = s.fillLinks(ctx, apiKey, result)
	if err != nil {
		s.logger.Err(err).Msg("cannot fetch links")
	}

	return result, nil
}

func (s *store) GetImpacts(ctx context.Context, apiKey string, r ContextGraphRequest) (*ContextGraphAggregationResult, error) {
	e := types.Entity{}
	err := s.dbCollection.
		FindOne(ctx, bson.M{"_id": r.ID}, options.FindOne().SetProjection(bson.M{"impact": 1})).
		Decode(&e)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	match := bson.M{
		"_id":  bson.M{"$in": e.Impacts},
		"type": types.EntityTypeService,
	}
	now := types.NewCpsTime()
	pipeline := s.queryBuilder.CreateTreeOfDepsAggregationPipeline(match, r.Query, r.SortRequest, r.Category, r.Search,
		r.WithFlags, now)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	result := &ContextGraphAggregationResult{}

	if cursor.Next(ctx) {
		err = cursor.Decode(result)
		if err != nil {
			return nil, err
		}
	}

	err = s.fillLinks(ctx, apiKey, result)
	if err != nil {
		s.logger.Err(err).Msg("cannot fetch links")
	}

	return result, nil
}

func (s *store) Create(ctx context.Context, request CreateRequest) (*Response, error) {
	var enabled bool
	if request.Enabled != nil {
		enabled = *request.Enabled
	}
	var sliAvailState int64
	if request.SliAvailState != nil {
		sliAvailState = *request.SliAvailState
	}
	entity := entityservice.EntityService{
		Entity: types.Entity{
			ID:            utils.NewID(),
			Name:          request.Name,
			Depends:       []string{},
			Impacts:       []string{},
			EnableHistory: []types.CpsTime{},
			Enabled:       enabled,
			Infos:         transformInfos(request.EditRequest),
			Type:          types.EntityTypeService,
			Category:      request.Category,
			ImpactLevel:   request.ImpactLevel,
			SliAvailState: sliAvailState,
			Created:       types.CpsTime{Time: time.Now()},
		},
		EntityPatternFields: request.EntityPatternFieldsRequest.ToModelWithoutFields(common.GetForbiddenFieldsInEntityPattern(mongo.EntityMongoCollection)),
		OutputTemplate:      request.OutputTemplate,
	}

	if request.ID == "" {
		request.ID = utils.NewID()
	}

	entity.ID = request.ID
	var response *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.dbCollection.InsertOne(ctx, entity)
		if err != nil {
			return err
		}
		response, err = s.GetOneBy(ctx, entity.ID)
		return err
	})

	return response, err
}

func (s *store) Update(ctx context.Context, request UpdateRequest) (*Response, ServiceChanges, error) {
	var service *Response
	serviceChanges := ServiceChanges{}
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		service = nil
		serviceChanges = ServiceChanges{}
		oldValues := &entityservice.EntityService{}
		pattern := request.EntityPatternFieldsRequest.ToModelWithoutFields(common.GetForbiddenFieldsInEntityPattern(mongo.EntityMongoCollection))
		update := bson.M{"$set": bson.M{
			"name":            request.Name,
			"output_template": request.OutputTemplate,
			"category":        request.Category,
			"impact_level":    request.ImpactLevel,
			"enabled":         request.Enabled,
			"infos":           transformInfos(request.EditRequest),
			"sli_avail_state": request.SliAvailState,

			"entity_pattern":                 pattern.EntityPattern,
			"corporate_entity_pattern":       pattern.CorporateEntityPattern,
			"corporate_entity_pattern_title": pattern.CorporateEntityPatternTitle,
		}}
		if request.CorporateEntityPattern != "" || len(request.EntityPattern) > 0 {
			update["$unset"] = bson.M{"old_entity_patterns": ""}
		}
		err := s.dbCollection.FindOneAndUpdate(
			ctx,
			bson.M{
				"_id":  request.ID,
				"type": types.EntityTypeService,
			},
			update,
			options.FindOneAndUpdate().
				SetProjection(bson.M{"enabled": 1, "entity_pattern": 1, "old_entity_patterns": 1}).
				SetReturnDocument(options.Before),
		).Decode(oldValues)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}
			return err
		}
		serviceChanges.IsToggled = request.Enabled != nil && oldValues.Enabled != *request.Enabled
		if oldValues.OldEntityPatterns.IsSet() {
			serviceChanges.IsPatternChanged = len(request.EntityPattern) > 0
		} else {
			serviceChanges.IsPatternChanged = !reflect.DeepEqual(oldValues.EntityPattern, request.EntityPattern)
		}
		service, err = s.GetOneBy(ctx, request.ID)
		return err
	})

	return service, serviceChanges, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, *types.Alarm, error) {
	res := false
	var alarm *types.Alarm
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		alarm = nil
		deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{
			"_id":  id,
			"type": types.EntityTypeService,
		})
		if err != nil || deleted == 0 {
			return err
		}

		// Delete open alarm.
		alarm = &types.Alarm{}
		err = s.alarmDbCollection.FindOneAndDelete(ctx, bson.M{"d": id, "v.resolved": nil}).Decode(alarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				alarm = nil
			} else {
				return err
			}
		}
		// Delete resolved alarms.
		_, err = s.alarmDbCollection.DeleteMany(ctx, bson.M{"d": id})
		if err != nil {
			return err
		}
		_, err = s.resolvedAlarmDbCollection.DeleteMany(ctx, bson.M{"d": id})
		if err != nil {
			return err
		}

		res = true
		return nil
	})

	return res, alarm, err
}

func (s *store) fillLinks(ctx context.Context, apiKey string, response *ContextGraphAggregationResult) error {
	if response == nil || len(response.Data) == 0 {
		return nil
	}

	reqEntities := make([]common.FetchLinksRequestItem, len(response.Data))
	for i, e := range response.Data {
		reqEntities[i] = common.FetchLinksRequestItem{
			EntityID: e.ID,
		}
	}
	linksRes, err := s.linksFetcher.Fetch(ctx, apiKey, common.FetchLinksRequest{Entities: reqEntities})
	if err != nil || linksRes == nil {
		return err
	}

	linksByEntityID := make(map[string]map[string]interface{}, len(reqEntities))
	for _, item := range linksRes.Data {
		if len(item.Links) > 0 {
			links := make(map[string]interface{}, len(item.Links))
			for category, link := range item.Links {
				links[category] = link
			}
			linksByEntityID[item.EntityID] = links
		}
	}

	for i, e := range response.Data {
		response.Data[i].Links = linksByEntityID[e.ID]
	}

	return nil
}

func transformInfos(request EditRequest) map[string]types.Info {
	infos := make(map[string]types.Info, len(request.Infos))
	for _, v := range request.Infos {
		infos[v.Name] = types.Info{
			Name:        v.Name,
			Description: v.Description,
			Value:       v.Value,
		}
	}

	return infos
}
