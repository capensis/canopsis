package pattern

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/sync/errgroup"
)

type Store interface {
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	GetByID(ctx context.Context, id, userID string) (*Response, error)
	Find(ctx context.Context, r ListRequest, userID string) (*AggregationResult, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, pattern Response, userID string) (bool, error)
	CountAlarms(ctx context.Context, r CountRequest, maxCount int64) (AlarmCountResponse, error)
	CountEntities(ctx context.Context, r CountRequest, maxCount int64) (EntityCountResponse, error)
	GetLiteralsFieldStats(ctx context.Context, allLiterals []string) (map[string][]LiteralFieldStats, error)
	GetEntityIDs(ctx context.Context, entityPattern pattern.Entity) ([]string, int64, error)
}

type store struct {
	client                        mongo.DbClient
	readClient                    mongo.DbClient
	collection                    mongo.DbCollection
	entityInfosPropertyCollection mongo.DbCollection
	entityCollection              mongo.DbCollection
	authorProvider                author.Provider

	linkedCollections []string

	defaultSearchByFields []string
	defaultSortBy         string

	pbhComputeChan chan<- rpc.PbehaviorRecomputeEvent

	serviceChangeChan chan<- entityservice.ChangeEntityMessage

	stateSettingsUpdatesChan chan statesetting.RuleUpdatedMessage

	transformer patternfields.Transformer

	logger zerolog.Logger
}

func NewStore(
	dbClient mongo.DbClient,
	readDbClient mongo.DbClient,
	pbhComputeChan chan<- rpc.PbehaviorRecomputeEvent,
	serviceChangeChan chan<- entityservice.ChangeEntityMessage,
	stateSettingsUpdatesChan chan statesetting.RuleUpdatedMessage,
	authorProvider author.Provider,
	transformer patternfields.Transformer,
	logger zerolog.Logger,
) Store {
	return &store{
		client:                        dbClient,
		collection:                    dbClient.Collection(mongo.PatternMongoCollection),
		readClient:                    readDbClient,
		entityInfosPropertyCollection: dbClient.Collection(mongo.EntityInfosPropertyCollection),
		entityCollection:              dbClient.Collection(mongo.EntityMongoCollection),
		authorProvider:                authorProvider,
		defaultSearchByFields:         []string{"_id", "author.name", "title"},
		defaultSortBy:                 "created",

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
			mongo.DeclareTicketRuleCollection,
			mongo.LinkRuleMongoCollection,
			mongo.AlarmTagCollection,
			mongo.StateSettingsMongoCollection,
			mongo.ScenarioCollection,
		},

		pbhComputeChan:           pbhComputeChan,
		serviceChangeChan:        serviceChangeChan,
		stateSettingsUpdatesChan: stateSettingsUpdatesChan,
		transformer:              transformer,
		logger:                   logger,
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

		err := s.transformEntityPatternToModel(ctx, request, &model)
		if err != nil {
			return err
		}

		_, err = s.collection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		response, err = s.GetByID(ctx, model.ID, model.Author)
		return err
	})

	return response, err
}

func (s *store) GetByID(ctx context.Context, id, userID string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{
		"_id": id,
		"$or": []bson.M{
			{"author": userID},
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

func (s *store) Find(ctx context.Context, request ListRequest, userID string) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	match := make([]bson.M, 0)

	if request.Corporate == nil {
		match = append(match, bson.M{"$or": []bson.M{
			{"author": userID},
			{"is_corporate": true},
		}})
	} else if *request.Corporate {
		match = append(match, bson.M{"is_corporate": true})
	} else {
		match = append(match, bson.M{"author": userID, "is_corporate": false})
	}

	if request.Type != "" {
		match = append(match, bson.M{"type": request.Type})
	}

	if len(match) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": match}})
	}

	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
	filter := mongoquery.GetSearchQuery(request.Search, s.defaultSearchByFields)
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
		mongoquery.GetSortQuery(sortBy, request.Sort),
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
	var stateSettingMsgs []statesetting.RuleUpdatedMessage
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		pbhIds = nil
		serviceIds = nil
		prevPattern := savedpattern.SavedPattern{}

		err := s.transformEntityPatternToModel(ctx, request, &model)
		if err != nil {
			return err
		}

		err = s.collection.FindOneAndUpdate(
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

		response, err = s.GetByID(ctx, model.ID, model.Author)
		if err != nil || response == nil {
			return err
		}

		err = s.updateLinkedModels(ctx, *response, request.Author, prevPattern.Aliases, model.Aliases)
		if err != nil {
			return err
		}

		if response.Type == savedpattern.TypeEntity && !reflect.DeepEqual(response.EntityPattern, prevPattern.EntityPattern) {
			pbhIds, err = s.findPbehaviors(ctx, *response)
			if err != nil {
				return err
			}

			serviceIds, err = s.findEntityServices(ctx, *response)
			if err != nil {
				return err
			}

			stateSettingMsgs, err = s.findStateSettings(ctx, *response, prevPattern.EntityPattern)
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
			s.serviceChangeChan <- entityservice.ChangeEntityMessage{
				ID:                      serviceId,
				EntityType:              types.EntityTypeService,
				IsServicePatternChanged: true,
			}
		}
	}

	if len(stateSettingMsgs) > 0 {
		for _, msg := range stateSettingMsgs {
			s.stateSettingsUpdatesChan <- msg
		}
	}

	return response, err
}

func (s *store) Delete(ctx context.Context, pattern Response, userID string) (bool, error) {
	var deleted int64

	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		result, err := s.collection.UpdateOne(ctx, bson.M{"_id": pattern.ID}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || result.MatchedCount == 0 {
			return err
		}

		deleted, err = s.collection.DeleteOne(ctx, bson.M{"_id": pattern.ID})
		if err != nil || deleted == 0 {
			return err
		}

		return s.cleanLinkedModels(ctx, pattern, userID)
	})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s *store) getNewAliases(prev, cur []string) []string {
	if len(prev) == 0 || len(cur) == 0 {
		return cur
	}

	prevMap := make(map[string]bool, len(prev))
	for _, v := range prev {
		prevMap[v] = true
	}

	added := make([]string, 0)
	for _, v := range cur {
		if !prevMap[v] {
			added = append(added, v)
		}
	}

	return added
}

func (s *store) updateLinkedModels(ctx context.Context, pattern Response, author string, prevAliases, newAliases []string) error {
	if !pattern.IsCorporate {
		return nil
	}

	switch pattern.Type {
	case savedpattern.TypeAlarm:
		return s.updateLinkedModelsOnAlarmUpdate(ctx, pattern, author)
	case savedpattern.TypeEntity:
		return s.updateLinkedModelsOnEntityUpdate(ctx, pattern, author, prevAliases, newAliases)
	case savedpattern.TypePbehavior:
		return s.updateLinkedModelsOnPbehaviorUpdate(ctx, pattern, author)
	case savedpattern.TypeWeatherService:
		return s.updateLinkedModelsOnWeatherServiceUpdate(ctx, pattern, author)
	default:
		return fmt.Errorf("unknown pattern type id=%s: %q", pattern.ID, pattern.Type)
	}
}

func (s *store) updateLinkedModelsOnEntityUpdate(ctx context.Context, pattern Response, author string, prevAliases, newAliases []string) error {
	addedAliases := s.getNewAliases(prevAliases, newAliases)
	for _, collection := range s.linkedCollections {
		switch collection {
		case mongo.ScenarioCollection:
			filter := bson.M{"actions.corporate_entity_pattern": pattern.ID}
			set := bson.M{
				"actions": bson.M{"$map": bson.M{
					"input": "$actions",
					"in": bson.M{"$cond": bson.M{
						"if": bson.M{"$eq": bson.A{"$$this.corporate_entity_pattern", pattern.ID}},
						"then": bson.M{"$mergeObjects": bson.A{
							"$$this",
							bson.M{
								"entity_pattern": pattern.EntityPattern.RemoveFields(
									patternfields.GetForbiddenFieldsInEntityPattern(collection),
								),
								"corporate_entity_pattern_title": pattern.Title,
							},
						}},
						"else": "$$this",
					}},
				}},
				"updated": datetime.NewCpsTime(),
				"author":  author,
			}
			// cannot clean removed aliases because they can be used in another entity pattern from the same document
			if len(addedAliases) > 0 {
				set["aliases"] = bson.M{"$setUnion": bson.A{
					bson.M{"$ifNull": bson.A{"$aliases", bson.A{}}},
					addedAliases,
				}}
			}

			_, err := s.client.Collection(collection).UpdateMany(ctx, filter, []bson.M{{"$set": set}})
			if err != nil {
				return fmt.Errorf("cannot update entity pattern: %w", err)
			}
		case mongo.MetaAlarmRulesMongoCollection,
			mongo.StateSettingsMongoCollection:
			var fields []string
			if collection == mongo.MetaAlarmRulesMongoCollection {
				fields = []string{"entity_pattern", "total_entity_pattern"}
			} else {
				fields = []string{"entity_pattern", "inherited_entity_pattern"}
			}

			for _, f := range fields {
				filter := bson.M{"corporate_" + f: pattern.ID}
				set := bson.M{
					f: pattern.EntityPattern.RemoveFields(
						patternfields.GetForbiddenFieldsInEntityPattern(collection),
					),
					"corporate_" + f + "_title": pattern.Title,
					"updated":                   datetime.NewCpsTime(),
					"author":                    author,
				}
				// cannot clean removed aliases because they can be used in another entity pattern from the same document
				if len(addedAliases) > 0 {
					set["aliases"] = bson.M{"$setUnion": bson.A{
						bson.M{"$ifNull": bson.A{"$aliases", bson.A{}}},
						addedAliases,
					}}
				}

				_, err := s.client.Collection(collection).UpdateMany(ctx, filter, []bson.M{{"$set": set}})
				if err != nil {
					return fmt.Errorf("cannot update entity pattern: %w", err)
				}
			}
		default:
			filter := bson.M{"corporate_entity_pattern": pattern.ID}
			update := bson.M{"$set": bson.M{
				"entity_pattern": pattern.EntityPattern.RemoveFields(
					patternfields.GetForbiddenFieldsInEntityPattern(collection),
				),
				"corporate_entity_pattern_title": pattern.Title,
				"aliases":                        newAliases, // can set newAliases because a document contains only one entity pattern
				"updated":                        datetime.NewCpsTime(),
				"author":                         author,
			}}
			_, err := s.client.Collection(collection).UpdateMany(ctx, filter, update)
			if err != nil {
				return fmt.Errorf("cannot update entity pattern: %w", err)
			}
		}
	}

	return nil
}

func (s *store) updateLinkedModelsOnAlarmUpdate(ctx context.Context, pattern Response, author string) error {
	for _, collection := range s.linkedCollections {
		var filter, update bson.M
		var opts *options.UpdateManyOptionsBuilder
		switch collection {
		case mongo.ScenarioCollection:
			filter = bson.M{
				"actions.corporate_alarm_pattern": pattern.ID,
			}
			update = bson.M{"$set": bson.M{
				"actions.$[action].alarm_pattern": pattern.AlarmPattern.RemoveFields(
					patternfields.GetForbiddenFieldsInAlarmPattern(collection),
					patternfields.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(collection),
				),
				"actions.$[action].corporate_alarm_pattern_title": pattern.Title,
				"updated": datetime.NewCpsTime(),
				"author":  author,
			}}
			opts = options.UpdateMany().SetArrayFilters([]any{bson.M{
				"action.corporate_alarm_pattern": pattern.ID,
			}})
		default:
			filter = bson.M{
				"corporate_alarm_pattern": pattern.ID,
			}
			update = bson.M{"$set": bson.M{
				"alarm_pattern": pattern.AlarmPattern.RemoveFields(
					patternfields.GetForbiddenFieldsInAlarmPattern(collection),
					patternfields.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(collection),
				),
				"corporate_alarm_pattern_title": pattern.Title,
				"updated":                       datetime.NewCpsTime(),
				"author":                        author,
			}}
		}

		_, err := s.client.Collection(collection).UpdateMany(ctx, filter, update, opts)
		if err != nil {
			return fmt.Errorf("cannot update alarm pattern: %w", err)
		}
	}

	return nil
}

func (s *store) updateLinkedModelsOnPbehaviorUpdate(ctx context.Context, pattern Response, author string) error {
	filter := bson.M{"corporate_pbehavior_pattern": pattern.ID}
	update := bson.M{"$set": bson.M{
		"pbehavior_pattern":                 pattern.PbehaviorPattern,
		"corporate_pbehavior_pattern_title": pattern.Title,
		"updated":                           datetime.NewCpsTime(),
		"author":                            author,
	}}
	for _, collection := range s.linkedCollections {
		_, err := s.client.Collection(collection).UpdateMany(ctx, filter, update)
		if err != nil {
			return fmt.Errorf("cannot update pbehavior pattern: %w", err)
		}
	}

	return nil
}

func (s *store) updateLinkedModelsOnWeatherServiceUpdate(ctx context.Context, pattern Response, author string) error {
	filter := bson.M{"corporate_weather_service_pattern": pattern.ID}
	update := bson.M{"$set": bson.M{
		"weather_service_pattern":                 pattern.WeatherServicePattern,
		"corporate_weather_service_pattern_title": pattern.Title,
		"updated": datetime.NewCpsTime(),
		"author":  author,
	}}
	for _, collection := range s.linkedCollections {
		_, err := s.client.Collection(collection).UpdateMany(ctx, filter, update)
		if err != nil {
			return fmt.Errorf("cannot update weather service pattern: %w", err)
		}
	}

	return nil
}

func (s *store) cleanLinkedModels(ctx context.Context, pattern Response, author string) error {
	if !pattern.IsCorporate {
		return nil
	}

	field := ""
	switch pattern.Type {
	case savedpattern.TypeAlarm:
		field = "corporate_alarm_pattern"
	case savedpattern.TypeEntity:
		field = "corporate_entity_pattern"
	case savedpattern.TypePbehavior:
		field = "corporate_pbehavior_pattern"
	case savedpattern.TypeWeatherService:
		field = "corporate_weather_service_pattern"
	default:
		return fmt.Errorf("unknown pattern type for deleted pattern id=%s: %q", pattern.ID, pattern.Type)
	}

	for _, collection := range s.linkedCollections {
		switch collection {
		case mongo.ScenarioCollection:
			filter := bson.M{"actions." + field: pattern.ID}
			update := bson.M{
				"$set": bson.M{
					"updated": datetime.NewCpsTime(),
					"author":  author,
				},
				"$unset": bson.M{
					"actions.$[action]." + field:            "",
					"actions.$[action]." + field + "_title": "",
				},
			}
			opts := options.UpdateMany().SetArrayFilters([]any{bson.M{"action." + field: pattern.ID}})
			_, err := s.client.Collection(collection).UpdateMany(ctx, filter, update, opts)
			if err != nil {
				return fmt.Errorf("cannot clean linked models: %w", err)
			}
		case mongo.MetaAlarmRulesMongoCollection,
			mongo.StateSettingsMongoCollection:
			fields := []string{field}
			if pattern.Type == savedpattern.TypeEntity {
				if collection == mongo.MetaAlarmRulesMongoCollection {
					fields = append(fields, "corporate_total_entity_pattern")
				} else {
					fields = append(fields, "corporate_inherited_entity_pattern")
				}
			}

			for _, field := range fields {
				filter := bson.M{field: pattern.ID}
				update := bson.M{
					"$set": bson.M{
						"updated": datetime.NewCpsTime(),
						"author":  author,
					},
					"$unset": bson.M{
						field:            "",
						field + "_title": "",
					},
				}
				_, err := s.client.Collection(collection).UpdateMany(ctx, filter, update)
				if err != nil {
					return fmt.Errorf("cannot clean linked models: %w", err)
				}
			}
		default:
			filter := bson.M{field: pattern.ID}
			update := bson.M{
				"$set": bson.M{
					"updated": datetime.NewCpsTime(),
					"author":  author,
				},
				"$unset": bson.M{
					field:            "",
					field + "_title": "",
				},
			}
			_, err := s.client.Collection(collection).UpdateMany(ctx, filter, update)
			if err != nil {
				return fmt.Errorf("cannot clean linked models: %w", err)
			}
		}
	}

	return nil
}

func (s *store) CountAlarms(ctx context.Context, r CountRequest, maxCount int64) (AlarmCountResponse, error) {
	res := AlarmCountResponse{}
	g, ctx := errgroup.WithContext(ctx)

	hasAlarm := len(r.AlarmPattern) > 0
	hasEntity := len(r.EntityPattern) > 0
	hasPbehavior := len(r.PbehaviorPattern) > 0

	combinedPipeline := make([]bson.M, 0)

	if hasAlarm {
		alarmQuery, err := db.AlarmPatternToMongoQuery(r.AlarmPattern, "")
		if err != nil {
			return res, err
		}

		alarmAddFields := r.AlarmPattern.GetMongoFields("")
		if hasEntity || hasPbehavior {
			combinedPipeline = append(combinedPipeline, alarmMatchStages(alarmAddFields, alarmQuery)...)
		}

		s.fetchCountAsync(ctx, g, mongo.AlarmMongoCollection,
			alarmMatchStages(alarmAddFields, alarmQuery), maxCount, &res.Alarms.AlarmPattern)
	}

	if hasPbehavior {
		pbhQuery, err := db.PbehaviorInfoPatternToMongoQuery(r.PbehaviorPattern, "v")
		if err != nil {
			return res, err
		}

		if hasAlarm || hasEntity {
			combinedPipeline = append(combinedPipeline, alarmMatchStages(nil, pbhQuery)...)
		}

		s.fetchCountAsync(ctx, g, mongo.AlarmMongoCollection,
			alarmMatchStages(nil, pbhQuery), maxCount, &res.Alarms.PbehaviorPattern)
	}

	if hasEntity {
		var err error
		r.EntityPattern, _, err = s.transformer.TransformAliases(ctx, r.EntityPattern, r)
		if err != nil {
			return res, err
		}

		entityQueryForAlarms, err := db.EntityPatternToMongoQuery(r.EntityPattern, "entity")
		if err != nil {
			return res, err
		}

		entityQuery, err := db.EntityPatternToMongoQuery(r.EntityPattern, "")
		if err != nil {
			return res, err
		}

		if hasAlarm || hasPbehavior {
			combinedPipeline = append(combinedPipeline, bson.M{"$match": bson.M{"v.resolved": nil}})
			combinedPipeline = append(combinedPipeline, alarmEntityLookupStages(entityQueryForAlarms)...)
		}

		s.fetchCountAsync(ctx, g, mongo.AlarmMongoCollection,
			append([]bson.M{{"$match": bson.M{"v.resolved": nil}}}, alarmEntityLookupStages(entityQueryForAlarms)...),
			maxCount, &res.Alarms.EntityPattern)
		s.fetchCountAsync(ctx, g, mongo.EntityMongoCollection,
			[]bson.M{{"$match": entityQuery}}, maxCount, &res.Entities.EntityPattern)
	}

	combinedFetched := len(combinedPipeline) > 0
	if combinedFetched {
		s.fetchCountAsync(ctx, g, mongo.AlarmMongoCollection, combinedPipeline, maxCount, &res.Alarms.Combined)
	}

	if err := g.Wait(); err != nil {
		return res, err
	}

	if !combinedFetched {
		switch {
		case hasAlarm:
			res.Alarms.Combined = res.Alarms.AlarmPattern
		case hasPbehavior:
			res.Alarms.Combined = res.Alarms.PbehaviorPattern
		case hasEntity:
			res.Alarms.Combined = res.Alarms.EntityPattern
		}
	}

	return res, nil
}

func (s *store) CountEntities(ctx context.Context, r CountRequest, maxCount int64) (EntityCountResponse, error) {
	res := EntityCountResponse{}
	g, ctx := errgroup.WithContext(ctx)

	hasAlarm := len(r.AlarmPattern) > 0
	hasEntity := len(r.EntityPattern) > 0
	hasPbehavior := len(r.PbehaviorPattern) > 0

	combinedPipeline := make([]bson.M, 0)

	if hasEntity {
		var err error
		r.EntityPattern, _, err = s.transformer.TransformAliases(ctx, r.EntityPattern, r)
		if err != nil {
			return res, err
		}

		entityQuery, err := db.EntityPatternToMongoQuery(r.EntityPattern, "")
		if err != nil {
			return res, err
		}

		if hasAlarm || hasPbehavior {
			combinedPipeline = append(combinedPipeline, bson.M{"$match": entityQuery})
		}

		s.fetchCountAsync(ctx, g, mongo.EntityMongoCollection,
			[]bson.M{{"$match": entityQuery}}, maxCount, &res.EntityPattern)
	}

	if hasPbehavior {
		pbhQuery, err := db.PbehaviorInfoPatternToMongoQuery(r.PbehaviorPattern, "")
		if err != nil {
			return res, err
		}

		if hasAlarm || hasEntity {
			combinedPipeline = append(combinedPipeline, bson.M{"$match": pbhQuery})
		}

		s.fetchCountAsync(ctx, g, mongo.EntityMongoCollection,
			[]bson.M{{"$match": pbhQuery}}, maxCount, &res.PbehaviorPattern)
	}

	if hasAlarm {
		alarmQuery, err := db.AlarmPatternToMongoQuery(r.AlarmPattern, "")
		if err != nil {
			return res, err
		}

		alarmAddFields := r.AlarmPattern.GetMongoFields("")
		if hasEntity || hasPbehavior {
			alarmQueryForEntities, err := db.AlarmPatternToMongoQuery(r.AlarmPattern, "alarm")
			if err != nil {
				return res, err
			}

			combinedPipeline = append(combinedPipeline, entityAlarmLookupStages()...)
			if len(alarmAddFields) > 0 {
				combinedPipeline = append(combinedPipeline, bson.M{"$addFields": r.AlarmPattern.GetMongoFields("alarm")})
			}
			combinedPipeline = append(combinedPipeline, bson.M{"$match": alarmQueryForEntities})
		}

		s.fetchCountAsync(ctx, g, mongo.AlarmMongoCollection,
			alarmMatchStages(alarmAddFields, alarmQuery), maxCount, &res.AlarmPattern)
	}

	combinedFetched := len(combinedPipeline) > 0
	if combinedFetched {
		s.fetchCountAsync(ctx, g, mongo.EntityMongoCollection, combinedPipeline, maxCount, &res.Combined)
	}

	if err := g.Wait(); err != nil {
		return res, err
	}

	if !combinedFetched {
		switch {
		case hasAlarm:
			res.Combined = res.AlarmPattern
		case hasPbehavior:
			res.Combined = res.PbehaviorPattern
		case hasEntity:
			res.Combined = res.EntityPattern
		}
	}

	return res, nil
}

func (s *store) fetchCount(
	ctx context.Context,
	collectionName string,
	pipeline []bson.M,
) (int64, int64, error) {
	collection := s.readClient.Collection(collectionName)
	pipeline = append(pipeline, bson.M{"$count": "total_count"})
	start := time.Now()
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, err
	}

	defer cursor.Close(ctx)

	res := AggregationResult{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return 0, 0, err
		}
	}

	if err = cursor.Err(); err != nil {
		return 0, 0, err
	}

	return res.GetTotal(), max(time.Since(start).Milliseconds(), 1), nil
}

func (s *store) fetchCountAsync(
	ctx context.Context,
	g *errgroup.Group,
	collection string,
	pipeline []bson.M,
	maxCount int64,
	dst *CountResponse,
) {
	g.Go(func() error {
		var err error
		dst.Count, dst.Millisecs, err = s.fetchCount(ctx, collection, pipeline)
		dst.OverLimit = dst.Count > maxCount

		return err
	})
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

func (s *store) findStateSettings(ctx context.Context, pattern Response, prevEntityPattern pattern.Entity) ([]statesetting.RuleUpdatedMessage, error) {
	if pattern.Type != savedpattern.TypeEntity {
		return nil, nil
	}

	cursor, err := s.client.Collection(mongo.StateSettingsMongoCollection).Find(ctx, bson.M{
		"method": bson.M{"$in": []string{statesetting.MethodInherited, statesetting.MethodDependencies}},
		"$or": []bson.M{
			{"corporate_entity_pattern": pattern.ID},
			{"corporate_inherited_entity_pattern": pattern.ID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("cannot find state settings: %w", err)
	}

	defer cursor.Close(ctx)

	now := datetime.NewCpsTime()
	msgs := make([]statesetting.RuleUpdatedMessage, 0)
	for cursor.Next(ctx) {
		m := statesetting.StateSetting{}
		err = cursor.Decode(&m)
		if err != nil {
			s.logger.Err(err).Msg("cannot decode state setting")
			continue
		}

		msg := statesetting.RuleUpdatedMessage{
			ID:         m.ID,
			NewPattern: m.EntityPattern,
			NewType:    m.Type,
			OldType:    m.Type,
			Updated:    now,
		}
		if m.CorporateEntityPattern == pattern.ID {
			msg.OldPattern = prevEntityPattern
		} else {
			msg.OldPattern = m.EntityPattern
		}

		msgs = append(msgs, msg)
	}

	if len(msgs) > 0 {
		_, err = s.client.Collection(mongo.EngineNotificationCollection).UpdateOne(
			ctx,
			bson.M{"_id": statesetting.StateSettingsNotificationID},
			bson.M{"$set": bson.M{"time": time.Now()}},
			options.UpdateOne().SetUpsert(true),
		)
		if err != nil {
			return nil, fmt.Errorf("cannot notify of state setting changes: %w", err)
		}
	}

	return msgs, nil
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
	case savedpattern.TypePbehavior:
		model.PbehaviorPattern = request.PbehaviorPattern
	case savedpattern.TypeWeatherService:
		model.WeatherServicePattern = request.WeatherServicePattern
	}

	return model
}

func (s *store) transformEntityPatternToModel(
	ctx context.Context,
	r EditRequest,
	model *savedpattern.SavedPattern,
) error {
	if r.Type == savedpattern.TypeEntity {
		var err error
		r.EntityPattern, model.Aliases, err = s.transformer.TransformAliases(ctx, r.EntityPattern, r)
		if err != nil {
			return err
		}

		model.EntityPattern = r.EntityPattern
	}

	return nil
}

func (s *store) GetLiteralsFieldStats(ctx context.Context, allLiterals []string) (map[string][]LiteralFieldStats, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"$text": bson.M{
					"$search": strings.Join(allLiterals, " "),
				},
			},
		},
		{
			"$set": bson.M{
				"infos": bson.M{
					"$objectToArray": "$infos",
				},
			},
		},
		{
			"$set": bson.M{
				"fields": bson.M{
					"$map": bson.M{
						"input": bson.M{
							"$filter": bson.M{
								"input": bson.M{
									"$concatArrays": bson.A{"$infos", bson.A{
										bson.M{
											"k": EntityFieldName,
											"v": bson.M{
												"value": "$name",
											},
										},
										bson.M{
											"k": EntityFieldComponent,
											"v": bson.M{
												"value": "$component",
											},
										},
									}},
								},
								"cond": bson.M{
									"$in": bson.A{"$$this.v.value", allLiterals},
								},
							},
						},
						"in": bson.M{
							"field": "$$m.k",
							"val":   "$$m.v.value",
						},
						"as": "m",
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":    0,
				"fields": 1,
			},
		},
		{
			"$unwind": "$fields",
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"val":   "$fields.val",
					"field": "$fields.field",
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": "$_id.val",
				"counts": bson.M{
					"$push": bson.M{
						"k": "$_id.field",
						"v": "$count",
					},
				},
			},
		},
		{
			"$set": bson.M{
				"counts": bson.M{
					"$sortArray": bson.M{
						"input":  "$counts",
						"sortBy": bson.M{"v": -1, "k": 1},
					},
				},
			},
		},
	}

	cursor, err := s.readClient.Collection(mongo.EntityMongoCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get literals field stats: %w", err)
	}

	var docs []struct {
		ID     string              `bson:"_id"`
		Counts []LiteralFieldStats `bson:"counts"`
	}

	if err := cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("failed to decode literals field stats: %w", err)
	}

	fieldStats := make(map[string][]LiteralFieldStats)
	for _, doc := range docs {
		fieldStats[doc.ID] = doc.Counts
	}

	return fieldStats, nil
}

func (s *store) GetEntityIDs(ctx context.Context, entityPattern pattern.Entity) ([]string, int64, error) {
	entityPatternQuery, err := db.EntityPatternToMongoQuery(entityPattern, "")
	if err != nil {
		return nil, 0, fmt.Errorf("failed to transform entity pattern to mongo query: %w", err)
	}

	start := time.Now()

	cursor, err := s.readClient.Collection(mongo.EntityMongoCollection).Aggregate(ctx, []bson.M{
		{
			"$match": entityPatternQuery,
		},
		{
			"$project": bson.M{"_id": 1},
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get entity ids: %w", err)
	}

	end := time.Since(start)

	var docs []struct {
		ID string `bson:"_id"`
	}

	err = cursor.All(ctx, &docs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode entity ids: %w", err)
	}

	if len(docs) == 0 {
		return nil, 0, nil
	}

	entityIDs := make([]string, len(docs))
	for i := range docs {
		entityIDs[i] = docs[i].ID
	}

	return entityIDs, max(end.Milliseconds(), 1), nil
}

// alarmMatchStages filters non-resolved alarms by an alarm pattern query.
func alarmMatchStages(addFields, query bson.M) []bson.M {
	stages := make([]bson.M, 0, 2)
	if len(addFields) > 0 {
		stages = append(stages, bson.M{"$addFields": addFields})
	}

	return append(stages, bson.M{"$match": bson.M{"$and": []bson.M{
		{"v.resolved": nil},
		query,
	}}})
}

// alarmEntityLookupStages joins alarms with their entity and filters by an entity pattern query.
func alarmEntityLookupStages(entityQuery bson.M) []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$match": entityQuery},
	}
}

// entityAlarmLookupStages joins entities with their non-resolved alarm.
func entityAlarmLookupStages() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
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
		{"$unwind": bson.M{"path": "$alarm", "preserveNullAndEmptyArrays": true}},
	}
}
