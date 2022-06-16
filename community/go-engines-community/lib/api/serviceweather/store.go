package serviceweather

import (
	"context"
	"errors"
	"sort"
	"time"

	alarmapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	pbehaviorlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const linkFetchTimeout = 30 * time.Second

type Store interface {
	Find(context.Context, ListRequest) (*AggregationResult, error)
	FindEntities(ctx context.Context, id, apiKey string, query EntitiesListRequest) (*EntityAggregationResult, error)
}

func NewStore(
	dbClient mongo.DbClient,
	legacyURL string,
	alarmStore alarmapi.Store,
	timezoneConfigProvider config.TimezoneConfigProvider,
	logger zerolog.Logger,
) Store {
	return &store{
		dbCollection: dbClient.Collection(mongo.EntityMongoCollection),

		alarmStore: alarmStore,
		links:      alarmapi.NewLinksFetcher(legacyURL, linkFetchTimeout),

		timezoneConfigProvider: timezoneConfigProvider,

		queryBuilder: NewMongoQueryBuilder(dbClient),

		logger: logger,
	}
}

type store struct {
	dbCollection mongo.DbCollection

	links      alarmapi.LinksFetcher
	alarmStore alarmapi.Store

	timezoneConfigProvider config.TimezoneConfigProvider

	queryBuilder *MongoQueryBuilder

	logger zerolog.Logger
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline, err := s.queryBuilder.CreateListAggregationPipeline(ctx, r)
	if err != nil {
		return nil, err
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var res AggregationResult
	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *store) FindEntities(ctx context.Context, id, apiKey string, query EntitiesListRequest) (*EntityAggregationResult, error) {
	var service libtypes.Entity
	err := s.dbCollection.FindOne(ctx, bson.M{
		"_id":     id,
		"type":    libtypes.EntityTypeService,
		"enabled": true,
	}).Decode(&service)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	location := s.timezoneConfigProvider.Get().Location

	pipeline := []bson.M{
		{"$match": bson.M{
			"_id":     bson.M{"$in": service.Depends},
			"enabled": true,
		}},
	}
	pipeline = append(pipeline, getFindEntitiesPipeline(location)...)
	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		s.getSort(query.SortBy, query.Sort),
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var res EntityAggregationResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	if !service.PbehaviorInfo.IsActive() {
		for i := range res.Data {
			res.Data[i].IsGrey = true
		}
	}

	if query.WithInstructions {
		alarmIds := make([]string, 0, len(res.Data))
		for _, v := range res.Data {
			if v.AlarmID != "" {
				alarmIds = append(alarmIds, v.AlarmID)
			}
		}

		assignedInstructionsMap, err := s.alarmStore.GetAssignedInstructionsMap(ctx, alarmIds)
		if err != nil {
			return nil, err
		}

		statusesByAlarm, err := s.alarmStore.GetInstructionExecutionStatuses(ctx, alarmIds)
		if err != nil {
			return nil, err
		}

		for idx, v := range res.Data {
			sort.Slice(assignedInstructionsMap[v.AlarmID], func(i, j int) bool {
				return assignedInstructionsMap[v.AlarmID][i].Name < assignedInstructionsMap[v.AlarmID][j].Name
			})

			assignedInstructions := assignedInstructionsMap[v.AlarmID]
			if assignedInstructions == nil {
				assignedInstructions = make([]alarmapi.AssignedInstruction, 0)
			}
			res.Data[idx].AssignedInstructions = &assignedInstructions
			res.Data[idx].IsAutoInstructionRunning = statusesByAlarm[v.AlarmID].AutoRunning
			res.Data[idx].IsAllAutoInstructionsCompleted = statusesByAlarm[v.AlarmID].AutoAllCompleted
			res.Data[idx].IsAutoInstructionFailed = statusesByAlarm[v.AlarmID].AutoFailed
			res.Data[idx].IsManualInstructionRunning = statusesByAlarm[v.AlarmID].ManualRunning
			res.Data[idx].IsManualInstructionWaitingResult = statusesByAlarm[v.AlarmID].ManualWaitingResult
		}
	}

	err = s.fillLinks(ctx, apiKey, &res)
	if err != nil {
		s.logger.Err(err).Msg("cannot fill links")
	}

	return &res, nil
}

func (s *store) fillLinks(ctx context.Context, apiKey string, result *EntityAggregationResult) error {
	linksEntities := make([]alarmapi.AlarmEntity, 0, len(result.Data))
	entities := make(map[string][]int, len(result.Data))
	for i, entity := range result.Data {
		if _, ok := entities[entity.ID]; !ok {
			linksEntities = append(linksEntities, alarmapi.AlarmEntity{
				EntityID: entity.ID,
			})
			entities[entity.ID] = make([]int, 0, 1)
		}
		// map entity ID with record number in result.Data list
		entities[entity.ID] = append(entities[entity.ID], i)
	}
	res, err := s.links.Fetch(ctx, apiKey, linksEntities)
	if err != nil || res == nil {
		return err
	}

	for _, rec := range res.Data {
		if l, ok := entities[rec.EntityID]; ok {
			for _, i := range l {
				result.Data[i].Links = make([]WeatherLink, 0, len(rec.Links))
				for category, link := range rec.Links {
					result.Data[i].Links = append(result.Data[i].Links, WeatherLink{
						Category: category,
						Links:    link,
					})
				}
			}
		}
	}

	return nil
}

func (s *store) getSort(sortBy, sort string) bson.M {
	if sortBy == "" {
		sortBy = common.SortAsc
	}

	if sortBy == "state" {
		sortBy = "state.val"
	}

	sortDir := 1
	if sort == common.SortDesc {
		sortDir = -1
	}

	sortQuery := bson.D{{Key: sortBy, Value: sortDir}}
	if sortBy != "name" {
		sortQuery = append(sortQuery, bson.E{Key: "name", Value: 1})
	}

	return bson.M{"$sort": sortQuery}
}

func getFindEntitiesPipeline(location *time.Location) []bson.M {
	year, month, day := time.Now().In(location).Date()
	truncatedInLocation := time.Date(year, month, day, 0, 0, 0, 0, location).Unix()

	pipeline := []bson.M{
		// Find category
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
		// Pbehavior
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorMongoCollection,
			"localField":   "pbehavior_info.id",
			"foreignField": "_id",
			"as":           "pbehavior",
		}},
		{"$unwind": bson.M{"path": "$pbehavior", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"pbehavior.last_comment": bson.M{"$arrayElemAt": bson.A{"$pbehavior.comments", -1}},
		}},
		{"$project": bson.M{
			"pbehavior.comments": 0,
		}},
		{"$lookup": bson.M{
			"from":         mongo.RightsMongoCollection,
			"foreignField": "_id",
			"localField":   "pbehavior.author",
			"as":           "pbehavior.author",
		}},
		{"$unwind": bson.M{"path": "$pbehavior.author", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"foreignField": "_id",
			"localField":   "pbehavior.type_",
			"as":           "pbehavior.type",
		}},
		{"$unwind": bson.M{"path": "$pbehavior.type", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorReasonMongoCollection,
			"foreignField": "_id",
			"localField":   "pbehavior.reason",
			"as":           "pbehavior.reason",
		}},
		{"$unwind": bson.M{"path": "$pbehavior.reason", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			// todo keep array for backward compatibility
			"pbehaviors": bson.M{"$cond": bson.M{
				"if":   "$pbehavior._id",
				"then": bson.A{"$pbehavior"},
				"else": bson.A{},
			}},
		}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"foreignField": "_id",
			"localField":   "pbehavior_info.type",
			"as":           "pbehavior_info_type",
		}},
		{"$unwind": bson.M{"path": "$pbehavior_info_type", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"pbehavior_info": bson.M{"$cond": bson.M{
				"if": "$pbehavior_info",
				"then": bson.M{"$mergeObjects": bson.A{
					"$pbehavior_info",
					bson.M{"icon_name": "$pbehavior_info_type.icon_name"},
				}},
				"else": nil,
			}},
		}},
		{"$project": bson.M{"pbehavior_info_type": 0}},
		// Event statistics
		{"$lookup": bson.M{
			"from":         mongo.EventStatistics,
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "stats",
		}},
		{"$unwind": bson.M{"path": "$stats", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			// stats counters with "last_event" prior "truncatedInLocation" represent as 0
			"stats.ok": bson.M{"$cond": bson.M{
				"if":   bson.M{"$gt": bson.A{"$stats.last_event", truncatedInLocation}},
				"then": "$stats.ok",
				"else": 0,
			}},
			"stats.ko": bson.M{"$cond": bson.M{
				"if":   bson.M{"$gt": bson.A{"$stats.last_event", truncatedInLocation}},
				"then": "$stats.ko",
				"else": 0,
			}},
		}},
		// Find connected alarm.
		{"$lookup": bson.M{
			"from": alarm.AlarmCollectionName,
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
		// Add alarms fields to result
		{"$addFields": bson.M{
			"alarm_id":         "$alarm._id",
			"connector":        "$alarm.v.connector",
			"connector_name":   "$alarm.v.connector_name",
			"component":        "$alarm.v.component",
			"resource":         "$alarm.v.resource",
			"state":            "$alarm.v.state",
			"status":           "$alarm.v.status",
			"snooze":           "$alarm.v.snooze",
			"ack":              "$alarm.v.ack",
			"ticket":           "$alarm.v.ticket",
			"last_update_date": "$alarm.v.last_update_date",
			"creation_date":    "$alarm.v.creation_date",
			"display_name":     "$alarm.v.display_name",
			"pbehavior_info":   "$pbehavior_info",
			"is_grey": bson.M{"$and": []bson.M{
				{"$ifNull": bson.A{"$pbehavior_info", false}},
				{"$ne": bson.A{"$pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
			}},
			"impact_state": bson.M{"$multiply": bson.A{"$alarm.v.state.val", "$impact_level"}},
		}},
	}
	pipeline = append(pipeline, getFindEntitiesIconPipeline()...)

	return pipeline
}

func getFindEntitiesIconPipeline() []bson.M {
	defaultVal := libtypes.AlarmStateTitleOK
	stateVals := []bson.M{
		{
			"case": bson.M{"$eq": bson.A{libtypes.AlarmStateMinor, "$state.val"}},
			"then": libtypes.AlarmStateTitleMinor,
		},
		{
			"case": bson.M{"$eq": bson.A{libtypes.AlarmStateMajor, "$state.val"}},
			"then": libtypes.AlarmStateTitleMajor,
		},
		{
			"case": bson.M{"$eq": bson.A{libtypes.AlarmStateCritical, "$state.val"}},
			"then": libtypes.AlarmStateTitleCritical,
		},
	}

	return []bson.M{
		{"$addFields": bson.M{
			"icon": bson.M{"$switch": bson.M{
				"branches": append(
					// If service is not active return pbehavior type icon.
					[]bson.M{{
						"case": bson.M{"$and": []bson.M{
							{"$ifNull": bson.A{"$pbehavior_info", false}},
							{"$ne": bson.A{"$pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
						}},
						"then": "$pbehavior_info.canonical_type",
					}},
					// Else return state icon.
					stateVals...,
				),
				"default": defaultVal,
			}},
		}},
	}
}
