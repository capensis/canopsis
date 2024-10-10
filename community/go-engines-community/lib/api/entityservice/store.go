package entityservice

import (
	"cmp"
	"context"
	"errors"
	"reflect"
	"sort"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
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
	GetDependencies(ctx context.Context, r ContextGraphRequest, userID string) (*ContextGraphAggregationResult, error)
	GetImpacts(ctx context.Context, r ContextGraphRequest, userID string) (*ContextGraphAggregationResult, error)
	Create(ctx context.Context, request CreateRequest) (*Response, error)
	Update(ctx context.Context, request UpdateRequest) (*Response, ServiceChanges, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
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
	entityCounters            mongo.DbCollection
	userDbCollection          mongo.DbCollection
	stateSettingDbCollection  mongo.DbCollection
	linkGenerator             link.Generator
	enableSameServiceNames    bool
	authorProvider            author.Provider
	logger                    zerolog.Logger
}

func NewStore(
	db mongo.DbClient,
	linkGenerator link.Generator,
	enableSameServiceNames bool,
	authorProvider author.Provider,
	logger zerolog.Logger,
) Store {
	return &store{
		dbClient:                  db,
		dbCollection:              db.Collection(mongo.EntityMongoCollection),
		alarmDbCollection:         db.Collection(mongo.AlarmMongoCollection),
		resolvedAlarmDbCollection: db.Collection(mongo.ResolvedAlarmMongoCollection),
		entityCounters:            db.Collection(mongo.EntityCountersCollection),
		userDbCollection:          db.Collection(mongo.UserCollection),
		stateSettingDbCollection:  db.Collection(mongo.StateSettingsMongoCollection),
		linkGenerator:             linkGenerator,
		enableSameServiceNames:    enableSameServiceNames,
		authorProvider:            authorProvider,
		logger:                    logger,
	}
}

func (s *store) GetOneBy(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{
		{"$match": bson.M{
			"_id":          id,
			"type":         types.EntityTypeService,
			"soft_deleted": bson.M{"$exists": false},
			"healthcheck":  bson.M{"$in": bson.A{nil, false}},
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
	}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
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

func (s *store) GetDependencies(ctx context.Context, r ContextGraphRequest, userID string) (*ContextGraphAggregationResult, error) {
	service := types.Entity{}
	err := s.dbCollection.
		FindOne(ctx, bson.M{
			"_id":          r.ID,
			"type":         bson.M{"$in": []string{types.EntityTypeService, types.EntityTypeComponent}},
			"soft_deleted": bson.M{"$exists": false},
		}).
		Decode(&service)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	now := datetime.NewCpsTime()
	match := bson.M{"soft_deleted": bson.M{"$exists": false}}
	switch service.Type {
	case types.EntityTypeService:
		match["services"] = service.ID
		match["_id"] = bson.M{"$ne": service.ID}
	case types.EntityTypeComponent:
		if service.StateInfo == nil {
			return &ContextGraphAggregationResult{Data: make([]ContextGraphEntity, 0)}, nil
		}

		match["component"] = service.ID
		match["_id"] = bson.M{"$ne": service.ID}
	}

	if r.DefineState {
		ec := entitycounters.EntityCounters{}
		err = s.entityCounters.FindOne(
			ctx, bson.M{"_id": service.ID}, options.FindOne().SetProjection(bson.M{"rule": 1})).Decode(&ec)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, err
		}
		if ec.Rule != nil {
			var entityPattern *pattern.Entity
			switch ec.Rule.Method {
			case statesetting.MethodInherited, statesetting.MethodDependencies:
				entityPattern = ec.Rule.InheritedEntityPattern
			}
			if entityPattern != nil {
				var patternMongoQuery bson.M
				patternMongoQuery, err = db.EntityPatternToMongoQuery(*entityPattern, "")
				if err != nil {
					return nil, err
				}

				match = bson.M{
					"$and": []bson.M{
						match,
						patternMongoQuery,
					},
				}
			}
		}
	}

	pipeline := s.getQueryBuilder().CreateTreeOfDepsAggregationPipeline(match, r.Query, r.SortRequest, r.Category, r.Search,
		r.WithFlags, r.DefineState, now)
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

	err = s.fillLinks(ctx, result, userID)
	if err != nil {
		s.logger.Err(err).Msg("cannot fetch links")
	}

	if r.WithFlags && r.DefineState {
		var defaultStateSetting StateSettingResponse
		for i := range result.Data {
			if result.Data[i].StateSetting == nil && result.Data[i].Type == types.EntityTypeService {
				if defaultStateSetting.ID == "" {
					err = s.stateSettingDbCollection.FindOne(ctx, bson.M{"_id": statesetting.ServiceID}).Decode(&defaultStateSetting)
					if err != nil {
						if errors.Is(err, mongodriver.ErrNoDocuments) {
							return result, nil
						}

						return nil, err
					}
				}

				result.Data[i].StateSetting = &defaultStateSetting
			}
		}
	}

	return result, nil
}

func (s *store) GetImpacts(ctx context.Context, r ContextGraphRequest, userID string) (*ContextGraphAggregationResult, error) {
	e := types.Entity{}
	err := s.dbCollection.FindOne(ctx,
		bson.M{"_id": r.ID, "soft_deleted": bson.M{"$exists": false}},
		options.FindOne().SetProjection(bson.M{"services": 1, "component": 1, "type": 1}),
	).Decode(&e)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	result := &ContextGraphAggregationResult{}
	match := make([]bson.M, 0)
	if e.Type == types.EntityTypeResource {
		match = append(match, bson.M{
			"_id":            e.Component,
			"type":           types.EntityTypeComponent,
			"state_info._id": bson.M{"$nin": bson.A{nil, ""}},
			"soft_deleted":   bson.M{"$exists": false},
		})
	}

	if len(e.Services) > 0 {
		match = append(match, bson.M{
			"_id":          bson.M{"$in": e.Services},
			"type":         types.EntityTypeService,
			"soft_deleted": bson.M{"$exists": false},
		})
	}

	if len(match) == 0 {
		result.Data = make([]ContextGraphEntity, 0)

		return result, nil
	}

	now := datetime.NewCpsTime()
	pipeline := s.getQueryBuilder().CreateTreeOfDepsAggregationPipeline(bson.M{"$or": match}, r.Query, r.SortRequest, r.Category, r.Search,
		r.WithFlags, false, now)
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

	err = s.fillLinks(ctx, result, userID)
	if err != nil {
		s.logger.Err(err).Msg("cannot fetch links")
	}

	return result, nil
}

func (s *store) Create(ctx context.Context, request CreateRequest) (*Response, error) {
	now := datetime.NewCpsTime()

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
			Author:        request.Author,
			EnableHistory: []datetime.CpsTime{},
			Enabled:       enabled,
			Infos:         transformInfos(request.EditRequest),
			Type:          types.EntityTypeService,
			Services:      []string{},
			Category:      request.Category,
			ImpactLevel:   request.ImpactLevel,
			SliAvailState: sliAvailState,
			Created:       now,
			Updated:       &now,
		},
		EntityPatternFields: request.EntityPatternFieldsRequest.ToModelWithoutFields(common.GetForbiddenFieldsInEntityPattern(mongo.EntityMongoCollection)),
		OutputTemplate:      request.OutputTemplate,
	}
	if request.Coordinates != nil {
		service.Coordinates = *request.Coordinates
	}

	service.ID = cmp.Or(request.ID, utils.NewID())
	var response *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		if !s.enableSameServiceNames {
			err := s.dbCollection.FindOne(ctx, bson.M{
				"name":         service.Name,
				"type":         types.EntityTypeService,
				"soft_deleted": nil,
			}).Err()
			if err == nil {
				return common.NewValidationError("name", "Name already exists.")
			}

			if !errors.Is(err, mongodriver.ErrNoDocuments) {
				return err
			}
		}

		_, err := s.dbCollection.InsertOne(ctx, service)
		if err != nil {
			return err
		}

		_, err = s.entityCounters.InsertOne(ctx, entitycounters.EntityCounters{
			ID:             service.ID,
			OutputTemplate: service.OutputTemplate,
		})
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

		"author":  request.Author,
		"updated": datetime.NewCpsTime(),
	}
	unset := bson.M{}

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
		if !s.enableSameServiceNames {
			err := s.dbCollection.FindOne(ctx, bson.M{
				"_id":          bson.M{"$ne": request.ID},
				"name":         request.Name,
				"type":         types.EntityTypeService,
				"soft_deleted": nil,
			}).Err()
			if err == nil {
				return common.NewValidationError("name", "Name already exists.")
			}

			if !errors.Is(err, mongodriver.ErrNoDocuments) {
				return err
			}
		}

		err := s.dbCollection.FindOneAndUpdate(
			ctx,
			bson.M{
				"_id":  request.ID,
				"type": types.EntityTypeService,
			},
			update,
			options.FindOneAndUpdate().
				SetProjection(bson.M{"enabled": 1, "entity_pattern": 1}).
				SetReturnDocument(options.Before),
		).Decode(oldValues)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}
			return err
		}
		serviceChanges.IsToggled = request.Enabled != nil && oldValues.Enabled != *request.Enabled
		serviceChanges.IsPatternChanged = !reflect.DeepEqual(oldValues.EntityPattern, request.EntityPattern)

		_, err = s.entityCounters.UpdateOne(
			ctx,
			bson.M{"_id": request.ID}, bson.M{
				"$set": bson.M{
					"output_template": request.OutputTemplate,
				},
			})
		if err != nil {
			return err
		}

		service, err = s.GetOneBy(ctx, request.ID)
		return err
	})

	return service, serviceChanges, err
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	updateRes, err := s.dbCollection.UpdateOne(ctx, bson.M{
		"_id":          id,
		"type":         types.EntityTypeService,
		"soft_deleted": nil,
	}, bson.M{"$set": bson.M{
		"enabled":      false,
		"author":       userID,
		"soft_deleted": datetime.NewCpsTime(),
	}})
	if err != nil || updateRes.MatchedCount == 0 {
		return false, err
	}

	return true, nil
}

func (s *store) fillLinks(ctx context.Context, response *ContextGraphAggregationResult, userID string) error {
	if response == nil || len(response.Data) == 0 {
		return nil
	}

	user, err := s.findUser(ctx, userID)
	if err != nil {
		return err
	}

	ids := make([]string, len(response.Data))
	for i, v := range response.Data {
		ids[i] = v.ID
	}

	linksByEntityId, err := s.linkGenerator.GenerateForEntities(ctx, ids, user)
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
	return entity.NewMongoQueryBuilder(s.dbClient, s.authorProvider)
}

func (s *store) findUser(ctx context.Context, id string) (link.User, error) {
	user := link.User{}
	cursor, err := s.userDbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$addFields": bson.M{"username": "$name"}},
	})
	if err != nil {
		return user, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&user)
		return user, err
	}

	return user, errors.New("user not found")
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
