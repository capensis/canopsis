package entityservice

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
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
	GetDependencies(ctx context.Context, r ContextGraphRequest) (*ContextGraphAggregationResult, error)
	GetImpacts(ctx context.Context, r ContextGraphRequest) (*ContextGraphAggregationResult, error)
	Create(ctx context.Context, request CreateRequest) (*Response, error)
	Update(ctx context.Context, request UpdateRequest) (*Response, ServiceChanges, error)
	Delete(ctx context.Context, id string) (bool, error)
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

	linkGenerator link.Generator

	logger zerolog.Logger
}

func NewStore(db mongo.DbClient, linkGenerator link.Generator, logger zerolog.Logger) Store {
	return &store{
		dbClient:                  db,
		dbCollection:              db.Collection(mongo.EntityMongoCollection),
		alarmDbCollection:         db.Collection(mongo.AlarmMongoCollection),
		resolvedAlarmDbCollection: db.Collection(mongo.ResolvedAlarmMongoCollection),

		linkGenerator: linkGenerator,

		logger: logger,
	}
}

func (s *store) GetOneBy(ctx context.Context, id string) (*Response, error) {
	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id, "type": types.EntityTypeService, "soft_deleted": bson.M{"$exists": false}}},
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

func (s *store) GetDependencies(ctx context.Context, r ContextGraphRequest) (*ContextGraphAggregationResult, error) {
	service := types.Entity{}
	err := s.dbCollection.
		FindOne(ctx, bson.M{
			"_id":          r.ID,
			"type":         types.EntityTypeService,
			"soft_deleted": bson.M{"$exists": false},
		}).
		Decode(&service)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	now := types.NewCpsTime()
	match := bson.M{
		"services":     service.ID,
		"soft_deleted": bson.M{"$exists": false},
	}
	pipeline := s.getQueryBuilder().CreateTreeOfDepsAggregationPipeline(match, r.Query, r.SortRequest, r.Category, r.Search,
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

	err = s.fillLinks(ctx, result)
	if err != nil {
		s.logger.Err(err).Msg("cannot fetch links")
	}

	return result, nil
}

func (s *store) GetImpacts(ctx context.Context, r ContextGraphRequest) (*ContextGraphAggregationResult, error) {
	e := types.Entity{}
	err := s.dbCollection.
		FindOne(ctx, bson.M{"_id": r.ID, "soft_deleted": bson.M{"$exists": false}}, options.FindOne().SetProjection(bson.M{"services": 1})).
		Decode(&e)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	result := &ContextGraphAggregationResult{}
	if len(e.Services) == 0 {
		result.Data = make([]entity.Entity, 0)
		return result, nil
	}

	match := bson.M{
		"_id":          bson.M{"$in": e.Services},
		"type":         types.EntityTypeService,
		"soft_deleted": bson.M{"$exists": false},
	}
	now := types.NewCpsTime()
	pipeline := s.getQueryBuilder().CreateTreeOfDepsAggregationPipeline(match, r.Query, r.SortRequest, r.Category, r.Search,
		r.WithFlags, now)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(result)
		if err != nil {
			return nil, err
		}
	}

	err = s.fillLinks(ctx, result)
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
	service := entityservice.EntityService{
		Entity: types.Entity{
			ID:            utils.NewID(),
			Name:          request.Name,
			EnableHistory: []types.CpsTime{},
			Enabled:       enabled,
			Infos:         transformInfos(request.EditRequest),
			Type:          types.EntityTypeService,
			Services:      []string{},
			Category:      request.Category,
			ImpactLevel:   request.ImpactLevel,
			SliAvailState: sliAvailState,
			Created:       types.CpsTime{Time: time.Now()},
		},
		EntityPatternFields: request.EntityPatternFieldsRequest.ToModelWithoutFields(common.GetForbiddenFieldsInEntityPattern(mongo.EntityMongoCollection)),
		OutputTemplate:      request.OutputTemplate,
	}
	if request.Coordinates != nil {
		service.Coordinates = *request.Coordinates
	}

	if request.ID == "" {
		request.ID = utils.NewID()
	}

	service.ID = request.ID
	var response *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.dbCollection.InsertOne(ctx, service)
		if err != nil {
			return err
		}
		response, err = s.GetOneBy(ctx, service.ID)
		return err
	})

	return response, err
}

func (s *store) Update(ctx context.Context, request UpdateRequest) (*Response, ServiceChanges, error) {
	pattern := request.EntityPatternFieldsRequest.ToModelWithoutFields(common.GetForbiddenFieldsInEntityPattern(mongo.EntityMongoCollection))
	set := bson.M{
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
	}
	unset := bson.M{}

	if request.CorporateEntityPattern != "" || len(request.EntityPattern) > 0 {
		unset["old_entity_patterns"] = ""
	}
	if request.Coordinates == nil {
		unset["coordinates"] = ""
	} else {
		set["coordinates"] = request.Coordinates
	}

	update := bson.M{"$set": set}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	var service *Response
	serviceChanges := ServiceChanges{}
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		service = nil
		serviceChanges = ServiceChanges{}
		oldValues := &entityservice.EntityService{}
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

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	updateRes, err := s.dbCollection.UpdateOne(ctx, bson.M{
		"_id":          id,
		"type":         types.EntityTypeService,
		"soft_deleted": nil,
	}, bson.M{"$set": bson.M{
		"enabled":      false,
		"soft_deleted": types.NewCpsTime(),
	}})
	if err != nil || updateRes.MatchedCount == 0 {
		return false, err
	}

	return true, nil
}

func (s *store) fillLinks(ctx context.Context, response *ContextGraphAggregationResult) error {
	if response == nil || len(response.Data) == 0 {
		return nil
	}

	ids := make([]string, len(response.Data))
	for i, v := range response.Data {
		ids[i] = v.ID
	}

	linksByEntityId, err := s.linkGenerator.GenerateForEntities(ctx, ids)
	if err != nil || len(linksByEntityId) == 0 {
		return err
	}

	for i, v := range response.Data {
		response.Data[i].Links = linksByEntityId[v.ID]
		for _, links := range response.Data[i].Links {
			sort.Slice(links, func(i, j int) bool {
				return links[i].Label < links[j].Label
			})
		}
	}

	return nil
}

func (s *store) getQueryBuilder() *entity.MongoQueryBuilder {
	return entity.NewMongoQueryBuilder(s.dbClient)
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
