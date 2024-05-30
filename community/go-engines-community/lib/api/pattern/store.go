package pattern

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

type Store interface {
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	GetById(ctx context.Context, id, userId string) (*Response, error)
	Find(ctx context.Context, r ListRequest, userId string) (*AggregationResult, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, pattern Response, userId string) (bool, error)
	CountAlarms(ctx context.Context, r CountRequest, maxCount int64) (CountAlarmsResponse, error)
	CountEntities(ctx context.Context, r CountRequest, maxCount int64) (CountEntitiesResponse, error)
}

type store struct {
	client         mongo.DbClient
	collection     mongo.DbCollection
	authorProvider author.Provider

	linkedCollections []string

	defaultSearchByFields []string
	defaultSortBy         string

	pbhComputeChan chan<- rpc.PbehaviorRecomputeEvent

	serviceChangeListener chan<- entityservice.ChangeEntityMessage

	logger zerolog.Logger
}

func NewStore(
	dbClient mongo.DbClient,
	pbhComputeChan chan<- rpc.PbehaviorRecomputeEvent,
	serviceChangeListener chan<- entityservice.ChangeEntityMessage,
	authorProvider author.Provider,
	logger zerolog.Logger,
) Store {
	return &store{
		client:         dbClient,
		collection:     dbClient.Collection(mongo.PatternMongoCollection),
		authorProvider: authorProvider,

		defaultSearchByFields: []string{"_id", "author.name", "title"},
		defaultSortBy:         "created",

		linkedCollections: []string{
			mongo.WidgetFiltersMongoCollection,
			mongo.EventFilterRuleCollection,
			mongo.MetaAlarmRulesMongoCollection,
			mongo.InstructionMongoCollection,
			mongo.PbehaviorMongoCollection,
			mongo.EntityMongoCollection,
			mongo.ResolveRuleMongoCollection,
			mongo.IdleRuleMongoCollection,
			mongo.DynamicInfosRulesMongoCollection,
			mongo.FlappingRuleMongoCollection,
			mongo.KpiFilterMongoCollection,
			mongo.DeclareTicketRuleMongoCollection,
			mongo.LinkRuleMongoCollection,
			mongo.AlarmTagCollection,
		},

		pbhComputeChan: pbhComputeChan,

		serviceChangeListener: serviceChangeListener,

		logger: logger,
	}
}

func (s *store) Insert(ctx context.Context, request EditRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	model := transformRequestToModel(request)
	model.ID = utils.NewID()
	model.Created = now
	model.Updated = now

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.collection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		response, err = s.GetById(ctx, model.ID, model.Author)
		return err
	})

	return response, err
}

func (s *store) GetById(ctx context.Context, id, userId string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{
		"_id": id,
		"$or": []bson.M{
			{"author": userId},
			{"is_corporate": true},
		},
	}}}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		response := Response{}
		err := cursor.Decode(&response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, request ListRequest, userId string) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	match := make([]bson.M, 0)

	if request.Corporate == nil {
		match = append(match, bson.M{"$or": []bson.M{
			{"author": userId},
			{"is_corporate": true},
		}})
	} else if *request.Corporate {
		match = append(match, bson.M{"is_corporate": true})
	} else {
		match = append(match, bson.M{"author": userId, "is_corporate": false})
	}

	if request.Type != "" {
		match = append(match, bson.M{"type": request.Type})
	}

	if len(match) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": match}})
	}

	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
	filter := common.GetSearchQuery(request.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := s.defaultSortBy
	if request.SortBy != "" {
		sortBy = request.SortBy
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		pipeline,
		common.GetSortQuery(sortBy, request.Sort),
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var result AggregationResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (s *store) Update(ctx context.Context, request EditRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	model := transformRequestToModel(request)
	model.ID = request.ID
	model.Updated = now

	var response *Response
	var pbhIds, serviceIds []string
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		pbhIds = nil
		serviceIds = nil
		prevPattern := savedpattern.SavedPattern{}
		err := s.collection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": request.ID},
			bson.M{"$set": model},
			options.FindOneAndUpdate().SetReturnDocument(options.Before),
		).Decode(&prevPattern)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}
			return err
		}

		response, err = s.GetById(ctx, model.ID, model.Author)
		if err != nil || response == nil {
			return err
		}

		err = s.updateLinkedModels(ctx, *response, request.Author)
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(response.EntityPattern, prevPattern.EntityPattern) {
			pbhIds, err = s.findPbehaviors(ctx, *response)
			if err != nil {
				return err
			}

			serviceIds, err = s.findEntityServices(ctx, *response)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if len(pbhIds) > 0 {
		s.pbhComputeChan <- rpc.PbehaviorRecomputeEvent{Ids: pbhIds}
	}

	if len(serviceIds) > 0 {
		for _, serviceId := range serviceIds {
			s.serviceChangeListener <- entityservice.ChangeEntityMessage{
				ID:                      serviceId,
				EntityType:              types.EntityTypeService,
				IsServicePatternChanged: true,
			}
		}
	}

	return response, err
}

func (s *store) Delete(ctx context.Context, pattern Response, userId string) (bool, error) {
	var deleted int64

	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		result, err := s.collection.UpdateOne(ctx, bson.M{"_id": pattern.ID}, bson.M{"$set": bson.M{"author": userId}})
		if err != nil || result.MatchedCount == 0 {
			return err
		}

		deleted, err = s.collection.DeleteOne(ctx, bson.M{"_id": pattern.ID})
		if err != nil || deleted == 0 {
			return err
		}

		return s.cleanLinkedModels(ctx, pattern, userId)
	})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s *store) updateLinkedModels(ctx context.Context, pattern Response, author string) error {
	if !pattern.IsCorporate {
		return nil
	}

	var filter bson.M
	switch pattern.Type {
	case savedpattern.TypeAlarm:
		filter = bson.M{"corporate_alarm_pattern": pattern.ID}
	case savedpattern.TypeEntity:
		filter = bson.M{"corporate_entity_pattern": pattern.ID}
	case savedpattern.TypePbehavior:
		filter = bson.M{"corporate_pbehavior_pattern": pattern.ID}
	case savedpattern.TypeWeatherService:
		filter = bson.M{"corporate_weather_service_pattern": pattern.ID}
	default:
		return fmt.Errorf("unknown pattern type id=%s: %q", pattern.ID, pattern.Type)
	}

	for _, collection := range s.linkedCollections {
		set := bson.M{
			"updated": datetime.NewCpsTime(),
			"author":  author,
		}

		switch pattern.Type {
		case savedpattern.TypeAlarm:
			set["alarm_pattern"] = pattern.AlarmPattern.RemoveFields(
				common.GetForbiddenFieldsInAlarmPattern(collection),
				common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(collection),
			)
			set["corporate_alarm_pattern_title"] = pattern.Title
		case savedpattern.TypeEntity:
			set["entity_pattern"] = pattern.EntityPattern.RemoveFields(
				common.GetForbiddenFieldsInEntityPattern(collection),
			)
			set["corporate_entity_pattern_title"] = pattern.Title
		case savedpattern.TypePbehavior:
			set["pbehavior_pattern"] = pattern.PbehaviorPattern
			set["corporate_pbehavior_pattern_title"] = pattern.Title
		case savedpattern.TypeWeatherService:
			set["weather_service_pattern"] = pattern.WeatherServicePattern
			set["corporate_weather_service_pattern_title"] = pattern.Title
		default:
			return fmt.Errorf("unknown pattern type id=%s: %q", pattern.ID, pattern.Type)
		}

		_, err := s.client.Collection(collection).UpdateMany(ctx, filter, bson.M{
			"$set": set,
		})
		if err != nil {
			return err
		}
	}

	switch pattern.Type {
	case savedpattern.TypeEntity:
		metaAlarmRulesCollection := mongo.MetaAlarmRulesMongoCollection
		_, err := s.client.Collection(metaAlarmRulesCollection).UpdateMany(ctx, bson.M{"corporate_total_entity_pattern": pattern.ID}, bson.M{
			"$set": bson.M{
				"total_entity_pattern": pattern.EntityPattern.RemoveFields(
					common.GetForbiddenFieldsInEntityPattern(metaAlarmRulesCollection),
				),
				"corporate_total_entity_pattern_title": pattern.Title,
				"updated":                              datetime.NewCpsTime(),
				"author":                               author,
			},
		})
		if err != nil {
			return err
		}

		scenarioCollection := mongo.ScenarioMongoCollection
		_, err = s.client.Collection(scenarioCollection).UpdateMany(ctx,
			bson.M{"actions.corporate_entity_pattern": pattern.ID},
			bson.M{"$set": bson.M{
				"actions.$[action].entity_pattern": pattern.EntityPattern.RemoveFields(
					common.GetForbiddenFieldsInEntityPattern(scenarioCollection),
				),
				"actions.$[action].corporate_entity_pattern_title": pattern.Title,
				"updated": datetime.NewCpsTime(),
				"author":  author,
			}},
			options.Update().SetArrayFilters(options.ArrayFilters{
				Filters: []interface{}{bson.M{"action.corporate_entity_pattern": pattern.ID}},
			}),
		)
		if err != nil {
			return err
		}
	case savedpattern.TypeAlarm:
		scenarioCollection := mongo.ScenarioMongoCollection
		_, err := s.client.Collection(scenarioCollection).UpdateMany(ctx,
			bson.M{"actions.corporate_alarm_pattern": pattern.ID},
			bson.M{"$set": bson.M{
				"actions.$[action].alarm_pattern": pattern.AlarmPattern.RemoveFields(
					common.GetForbiddenFieldsInAlarmPattern(scenarioCollection),
					common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(scenarioCollection),
				),
				"actions.$[action].corporate_alarm_pattern_title": pattern.Title,
				"updated": datetime.NewCpsTime(),
				"author":  author,
			}},
			options.Update().SetArrayFilters(options.ArrayFilters{
				Filters: []interface{}{bson.M{"action.corporate_alarm_pattern": pattern.ID}},
			}),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) cleanLinkedModels(ctx context.Context, pattern Response, author string) error {
	if !pattern.IsCorporate {
		return nil
	}

	f := ""
	switch pattern.Type {
	case savedpattern.TypeAlarm:
		f = "corporate_alarm_pattern"
	case savedpattern.TypeEntity:
		f = "corporate_entity_pattern"
	case savedpattern.TypePbehavior:
		f = "corporate_pbehavior_pattern"
	case savedpattern.TypeWeatherService:
		f = "corporate_weather_service_pattern"
	default:
		return fmt.Errorf("unknown pattern type for deleted pattern id=%s: %q", pattern.ID, pattern.Type)
	}

	for _, collection := range s.linkedCollections {
		_, err := s.client.Collection(collection).UpdateMany(ctx, bson.M{f: pattern.ID}, bson.M{
			"$set": bson.M{
				"updated": datetime.NewCpsTime(),
				"author":  author,
			},
			"$unset": bson.M{
				f:            "",
				f + "_title": "",
			},
		})
		if err != nil {
			return err
		}
	}

	switch pattern.Type {
	case savedpattern.TypeEntity:
		_, err := s.client.Collection(mongo.MetaAlarmRulesMongoCollection).UpdateMany(ctx, bson.M{"corporate_total_entity_pattern": pattern.ID}, bson.M{
			"$set": bson.M{
				"updated": datetime.NewCpsTime(),
				"author":  author,
			},
			"$unset": bson.M{
				"corporate_total_entity_pattern":       "",
				"corporate_total_entity_pattern_title": "",
			},
		})
		if err != nil {
			return err
		}

		_, err = s.client.Collection(mongo.ScenarioMongoCollection).UpdateMany(ctx,
			bson.M{"actions.corporate_entity_pattern": pattern.ID},
			bson.M{
				"$set": bson.M{
					"updated": datetime.NewCpsTime(),
					"author":  author,
				},
				"$unset": bson.M{
					"actions.$[action].corporate_entity_pattern":       "",
					"actions.$[action].corporate_entity_pattern_title": "",
				},
			},
			options.Update().SetArrayFilters(options.ArrayFilters{
				Filters: []interface{}{bson.M{"action.corporate_entity_pattern": pattern.ID}},
			}),
		)
		if err != nil {
			return err
		}
	case savedpattern.TypeAlarm:
		_, err := s.client.Collection(mongo.ScenarioMongoCollection).UpdateMany(ctx,
			bson.M{"actions.corporate_alarm_pattern": pattern.ID},
			bson.M{
				"$set": bson.M{
					"updated": datetime.NewCpsTime(),
					"author":  author,
				},
				"$unset": bson.M{
					"actions.$[action].corporate_alarm_pattern":       "",
					"actions.$[action].corporate_alarm_pattern_title": "",
				},
			},
			options.Update().SetArrayFilters(options.ArrayFilters{
				Filters: []interface{}{bson.M{"action.corporate_alarm_pattern": pattern.ID}},
			}),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) CountAlarms(ctx context.Context, r CountRequest, maxCount int64) (CountAlarmsResponse, error) {
	res := CountAlarmsResponse{}
	g, ctx := errgroup.WithContext(ctx)
	var err error
	var alarmPatternQuery, entityPatternQuery, pbhPatternQuery bson.M
	alarmsPipeline := make([]bson.M, 0)
	entitiesPipeline := make([]bson.M, 0)
	var alarmPatternCount, entityPatternCount, pbhPatternCount, alarmsCount, entitiesCount Count
	if len(r.AlarmPattern) > 0 {
		alarmPatternQuery, err = db.AlarmPatternToMongoQuery(r.AlarmPattern, "")
		if err != nil {
			return res, err
		}

		if len(r.PbehaviorPattern) > 0 || len(r.EntityPattern) > 0 {
			alarmsPipeline = append(alarmsPipeline, bson.M{"$match": bson.M{"$and": []bson.M{
				{"v.resolved": nil},
				alarmPatternQuery,
			}}})
		}
	}
	if len(r.PbehaviorPattern) > 0 {
		pbhPatternQuery, err = db.PbehaviorInfoPatternToMongoQuery(r.PbehaviorPattern, "v")
		if err != nil {
			return res, err
		}

		if len(r.AlarmPattern) > 0 || len(r.EntityPattern) > 0 {
			alarmsPipeline = append(alarmsPipeline, bson.M{"$match": bson.M{"$and": []bson.M{
				{"v.resolved": nil},
				pbhPatternQuery,
			}}})
		}
	}
	if len(r.EntityPattern) > 0 {
		entityPatternQuery, err = db.EntityPatternToMongoQuery(r.EntityPattern, "entity")
		if err != nil {
			return res, err
		}

		if len(r.AlarmPattern) > 0 || len(r.PbehaviorPattern) > 0 {
			alarmsPipeline = append(alarmsPipeline,
				bson.M{"$match": bson.M{"v.resolved": nil}},
				bson.M{"$lookup": bson.M{
					"from":         mongo.EntityMongoCollection,
					"localField":   "d",
					"foreignField": "_id",
					"as":           "entity",
				}},
				bson.M{"$unwind": "$entity"},
				bson.M{"$match": entityPatternQuery},
			)
		}

		entityPatternQueryForEntities, err := db.EntityPatternToMongoQuery(r.EntityPattern, "")
		if err != nil {
			return res, err
		}

		entitiesPipeline = append(entitiesPipeline, bson.M{"$match": entityPatternQueryForEntities})
	}

	if len(alarmPatternQuery) > 0 {
		g.Go(func() error {
			alarmPatternCount.Count, err = s.fetchCount(ctx, s.client.Collection(mongo.AlarmMongoCollection),
				[]bson.M{{"$match": bson.M{"$and": []bson.M{
					{"v.resolved": nil},
					alarmPatternQuery,
				}}}})
			alarmPatternCount.OverLimit = alarmPatternCount.Count > maxCount

			return err
		})
	}
	if len(pbhPatternQuery) > 0 {
		g.Go(func() error {
			var err error
			pbhPatternCount.Count, err = s.fetchCount(ctx, s.client.Collection(mongo.AlarmMongoCollection),
				[]bson.M{{"$match": bson.M{"$and": []bson.M{
					{"v.resolved": nil},
					pbhPatternQuery,
				}}}})
			pbhPatternCount.OverLimit = pbhPatternCount.Count > maxCount

			return err
		})
	}
	if len(entityPatternQuery) > 0 {
		g.Go(func() error {
			var err error
			entityPatternCount.Count, err = s.fetchCount(ctx, s.client.Collection(mongo.AlarmMongoCollection),
				[]bson.M{
					{"$match": bson.M{"v.resolved": nil}},
					{"$lookup": bson.M{
						"from":         mongo.EntityMongoCollection,
						"localField":   "d",
						"foreignField": "_id",
						"as":           "entity",
					}},
					{"$unwind": "$entity"},
					{"$match": entityPatternQuery},
				})
			entityPatternCount.OverLimit = entityPatternCount.Count > maxCount

			return err
		})
	}
	fetchAlarmsCount := false
	if len(alarmsPipeline) > 0 {
		fetchAlarmsCount = true
		g.Go(func() error {
			var err error
			alarmsCount.Count, err = s.fetchCount(ctx, s.client.Collection(mongo.AlarmMongoCollection), alarmsPipeline)
			alarmsCount.OverLimit = alarmsCount.Count > maxCount

			return err
		})
	}
	if len(entitiesPipeline) > 0 {
		g.Go(func() error {
			var err error
			entitiesCount.Count, err = s.fetchCount(ctx, s.client.Collection(mongo.EntityMongoCollection), entitiesPipeline)
			entitiesCount.OverLimit = entitiesCount.Count > maxCount

			return err
		})
	}

	if err := g.Wait(); err != nil {
		return res, err
	}

	res.AlarmPattern = alarmPatternCount
	res.PbehaviorPattern = pbhPatternCount
	res.EntityPattern = entityPatternCount
	if fetchAlarmsCount {
		res.All = alarmsCount
	} else if len(r.AlarmPattern) > 0 {
		res.All = alarmPatternCount
	} else if len(r.PbehaviorPattern) > 0 {
		res.All = pbhPatternCount
	} else if len(r.EntityPattern) > 0 {
		res.All = entityPatternCount
	}
	res.Entities = entitiesCount

	return res, nil
}

func (s *store) CountEntities(ctx context.Context, r CountRequest, maxCount int64) (CountEntitiesResponse, error) {
	res := CountEntitiesResponse{}
	g, ctx := errgroup.WithContext(ctx)
	var err error
	var alarmPatternQuery, entityPatternQuery, pbhPatternQuery bson.M
	entitiesPipeline := make([]bson.M, 0)
	var alarmPatternCount, entityPatternCount, pbhPatternCount, entitiesCount Count
	if len(r.EntityPattern) > 0 {
		entityPatternQuery, err = db.EntityPatternToMongoQuery(r.EntityPattern, "")
		if err != nil {
			return res, err
		}

		if len(r.AlarmPattern) > 0 || len(r.PbehaviorPattern) > 0 {
			entitiesPipeline = append(entitiesPipeline, bson.M{"$match": entityPatternQuery})
		}
	}
	if len(r.PbehaviorPattern) > 0 {
		pbhPatternQuery, err = db.PbehaviorInfoPatternToMongoQuery(r.PbehaviorPattern, "")
		if err != nil {
			return res, err
		}

		if len(r.AlarmPattern) > 0 || len(r.EntityPattern) > 0 {
			entitiesPipeline = append(entitiesPipeline, bson.M{"$match": pbhPatternQuery})
		}
	}
	if len(r.AlarmPattern) > 0 {
		alarmPatternQuery, err = db.AlarmPatternToMongoQuery(r.AlarmPattern, "")
		if err != nil {
			return res, err
		}

		if len(r.PbehaviorPattern) > 0 || len(r.EntityPattern) > 0 {
			alarmPatternQueryForEntities, err := db.AlarmPatternToMongoQuery(r.AlarmPattern, "alarm")
			if err != nil {
				return res, err
			}

			entitiesPipeline = append(entitiesPipeline,
				bson.M{"$lookup": bson.M{
					"from": mongo.AlarmMongoCollection,
					"let":  bson.M{"id": "$_id"},
					"pipeline": []bson.M{
						{"$match": bson.M{"$and": []bson.M{
							{"$expr": bson.M{"$eq": bson.A{"$d", "$$id"}}},
							{"v.resolved": nil},
						}}},
						{"$limit": 1},
					},
					"as": "alarm",
				}},
				bson.M{"$unwind": bson.M{"path": "$alarm", "preserveNullAndEmptyArrays": true}},
				bson.M{"$match": alarmPatternQueryForEntities},
			)
		}
	}

	if len(entityPatternQuery) > 0 {
		g.Go(func() error {
			var err error
			entityPatternCount.Count, err = s.fetchCount(ctx, s.client.Collection(mongo.EntityMongoCollection),
				[]bson.M{{"$match": entityPatternQuery}})
			entityPatternCount.OverLimit = entityPatternCount.Count > maxCount

			return err
		})
	}
	if len(pbhPatternQuery) > 0 {
		g.Go(func() error {
			var err error
			pbhPatternCount.Count, err = s.fetchCount(ctx, s.client.Collection(mongo.EntityMongoCollection),
				[]bson.M{{"$match": pbhPatternQuery}})
			pbhPatternCount.OverLimit = pbhPatternCount.Count > maxCount

			return err
		})
	}
	if len(alarmPatternQuery) > 0 {
		g.Go(func() error {
			alarmPatternCount.Count, err = s.fetchCount(ctx, s.client.Collection(mongo.AlarmMongoCollection),
				[]bson.M{
					{"$match": bson.M{"v.resolved": nil}},
					{"$match": alarmPatternQuery},
				})
			alarmPatternCount.OverLimit = alarmPatternCount.Count > maxCount

			return err
		})
	}

	fetchEntitiesCount := false
	if len(entitiesPipeline) > 0 {
		fetchEntitiesCount = true
		g.Go(func() error {
			var err error
			entitiesCount.Count, err = s.fetchCount(ctx, s.client.Collection(mongo.EntityMongoCollection), entitiesPipeline)
			entitiesCount.OverLimit = entitiesCount.Count > maxCount

			return err
		})
	}

	if err := g.Wait(); err != nil {
		return res, err
	}

	res.AlarmPattern = alarmPatternCount
	res.PbehaviorPattern = pbhPatternCount
	res.EntityPattern = entityPatternCount
	if fetchEntitiesCount {
		res.All = entitiesCount
	} else if len(r.AlarmPattern) > 0 {
		res.All = alarmPatternCount
	} else if len(r.PbehaviorPattern) > 0 {
		res.All = pbhPatternCount
	} else if len(r.EntityPattern) > 0 {
		res.All = entityPatternCount
	}

	return res, nil
}

func (s *store) fetchCount(
	ctx context.Context,
	collection mongo.DbCollection,
	pipeline []bson.M,
) (int64, error) {
	pipeline = append(pipeline, bson.M{"$count": "total_count"})
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	res := AggregationResult{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
	}

	return res.GetTotal(), err
}

func (s *store) findPbehaviors(ctx context.Context, pattern Response) ([]string, error) {
	if pattern.Type != savedpattern.TypeEntity {
		return nil, nil
	}

	cursor, err := s.client.Collection(mongo.PbehaviorMongoCollection).Find(ctx, bson.M{
		"corporate_entity_pattern": pattern.ID,
	}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	ids := make([]string, 0)
	for cursor.Next(ctx) {
		pbh := pbehavior.PBehavior{}
		err := cursor.Decode(&pbh)
		if err != nil {
			s.logger.Err(err).Msg("cannot decode pbehavior")
			continue
		}
		ids = append(ids, pbh.ID)
	}

	return ids, nil
}

func (s *store) findEntityServices(ctx context.Context, pattern Response) ([]string, error) {
	if pattern.Type != savedpattern.TypeEntity {
		return nil, nil
	}

	cursor, err := s.client.Collection(mongo.EntityMongoCollection).Find(ctx, bson.M{
		"type":                     types.EntityTypeService,
		"corporate_entity_pattern": pattern.ID,
	}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	ids := make([]string, 0)
	for cursor.Next(ctx) {
		entity := types.Entity{}
		err := cursor.Decode(&entity)
		if err != nil {
			s.logger.Err(err).Msg("cannot decode entity service")
			continue
		}
		ids = append(ids, entity.ID)
	}

	return ids, nil
}

func transformRequestToModel(request EditRequest) savedpattern.SavedPattern {
	model := savedpattern.SavedPattern{
		Title:       request.Title,
		Type:        request.Type,
		IsCorporate: *request.IsCorporate,
		Author:      request.Author,
	}

	switch request.Type {
	case savedpattern.TypeAlarm:
		model.AlarmPattern = request.AlarmPattern
	case savedpattern.TypeEntity:
		model.EntityPattern = request.EntityPattern
	case savedpattern.TypePbehavior:
		model.PbehaviorPattern = request.PbehaviorPattern
	case savedpattern.TypeWeatherService:
		model.WeatherServicePattern = request.WeatherServicePattern
	}

	return model
}
