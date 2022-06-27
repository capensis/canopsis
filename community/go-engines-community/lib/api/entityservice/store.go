package entityservice

import (
	"context"
	"errors"
	"reflect"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	GetOneBy(ctx context.Context, id string) (*Response, error)
	GetDependencies(ctx context.Context, id string, query pagination.Query) (*ContextGraphAggregationResult, error)
	GetImpacts(ctx context.Context, id string, query pagination.Query) (*ContextGraphAggregationResult, error)
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
}

func NewStore(db mongo.DbClient) Store {
	return &store{
		dbClient:                  db,
		dbCollection:              db.Collection(mongo.EntityMongoCollection),
		alarmDbCollection:         db.Collection(mongo.AlarmMongoCollection),
		resolvedAlarmDbCollection: db.Collection(mongo.ResolvedAlarmMongoCollection),
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

func (s *store) GetDependencies(ctx context.Context, id string, q pagination.Query) (*ContextGraphAggregationResult, error) {
	service := types.Entity{}
	err := s.dbCollection.
		FindOne(ctx, bson.M{"_id": id, "type": types.EntityTypeService}).
		Decode(&service)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": service.Depends}}},
		// Find alarms
		{"$project": bson.M{
			"entity": "$$ROOT",
		}},
		{"$lookup": bson.M{
			"from": mongo.AlarmMongoCollection,
			"let":  bson.M{"eid": "$entity._id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$and": []bson.M{
					{"$expr": bson.M{"$eq": bson.A{"$d", "$$eid"}}},
					// Get only open alarm.
					{"v.resolved": nil},
				}}},
				{"$limit": 1},
			},
			"as": "alarm",
		}},
		{"$unwind": bson.M{"path": "$alarm", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"impact_state": bson.M{"$cond": bson.M{"if": "$alarm.v.state.val", "else": 0,
				"then": bson.M{"$multiply": bson.A{"$alarm.v.state.val", "$entity.impact_level"}},
			}},
		}},
	}
	projectPipeline := []bson.M{
		// Find category
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "entity.category",
			"foreignField": "_id",
			"as":           "entity.category",
		}},
		{"$unwind": bson.M{"path": "$entity.category", "preserveNullAndEmptyArrays": true}},
		// Find dependencies
		{"$addFields": bson.M{
			"has_dependencies": bson.M{"$and": []bson.M{
				{"$eq": bson.A{"$entity.type", types.EntityTypeService}},
				{"$gt": bson.A{bson.M{"$size": "$entity.depends"}, 0}},
			}},
		}},
	}
	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		q,
		pipeline,
		bson.M{"$sort": bson.D{{Key: "impact_state", Value: -1}, {Key: "entity._id", Value: 1}}},
		projectPipeline,
	))
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

	return result, nil
}

func (s *store) GetImpacts(ctx context.Context, id string, q pagination.Query) (*ContextGraphAggregationResult, error) {
	entity := types.Entity{}
	err := s.dbCollection.
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&entity)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{
			"_id":  bson.M{"$in": entity.Impacts},
			"type": types.EntityTypeService,
		}},
		// Find alarms
		{"$project": bson.M{
			"entity": "$$ROOT",
		}},
		{"$lookup": bson.M{
			"from": mongo.AlarmMongoCollection,
			"let":  bson.M{"eid": "$entity._id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$and": []bson.M{
					{"$expr": bson.M{"$eq": bson.A{"$d", "$$eid"}}},
					// Get only open alarm.
					{"v.resolved": nil},
				}}},
				{"$limit": 1},
			},
			"as": "alarm",
		}},
		{"$unwind": bson.M{"path": "$alarm", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"impact_state": bson.M{"$cond": bson.M{"if": "$alarm.v.state.val", "else": 0,
				"then": bson.M{"$multiply": bson.A{"$alarm.v.state.val", "$entity.impact_level"}},
			}},
		}},
	}
	const entitiesListLimit = 100
	projectPipeline := []bson.M{
		// Find category
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "entity.category",
			"foreignField": "_id",
			"as":           "entity.category",
		}},
		{"$unwind": bson.M{"path": "$entity.category", "preserveNullAndEmptyArrays": true}},
		// Find impacts
		{"$graphLookup": bson.M{
			"from":                    mongo.EntityMongoCollection,
			"startWith":               "$entity.impact",
			"connectFromField":        "impact",
			"connectToField":          "_id",
			"as":                      "service_impacts",
			"restrictSearchWithMatch": bson.M{"type": types.EntityTypeService},
			"maxDepth":                0,
		}},
		{"$addFields": bson.M{
			"has_impacts": bson.M{"$gt": bson.A{bson.M{"$size": "$service_impacts"}, 0}},
			// entity.impact and entity.depends arrays of the output document are
			// limited by 100 items each to prevent BSONObjectTooLarge error
			"entity.impact":  bson.M{"$slice": bson.A{"$entity.impact", entitiesListLimit}},
			"entity.depends": bson.M{"$slice": bson.A{"$entity.depends", entitiesListLimit}},
		}},
		{"$project": bson.M{"service_impacts": 0}},
	}
	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		q,
		pipeline,
		bson.M{"$sort": bson.D{{Key: "impact_state", Value: -1}, {Key: "entity._id", Value: 1}}},
		projectPipeline,
	))
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
