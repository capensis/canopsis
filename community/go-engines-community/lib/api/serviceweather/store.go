package serviceweather

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	alarmapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	pbehaviorlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const linkFetchTimeout = 30 * time.Second

type Store interface {
	Find(context.Context, ListRequest) (*AggregationResult, error)
	FindEntities(ctx context.Context, id, apiKey string, query EntitiesListRequest) (*EntityAggregationResult, error)
}

func NewStore(
	dbClient mongo.DbClient,
	legacyURL fmt.Stringer,
	alarmStore alarmapi.Store,
	timezoneConfigProvider config.TimezoneConfigProvider,
) Store {
	return &store{
		dbCollection:           dbClient.Collection(mongo.EntityMongoCollection),
		pbehaviorCollection:    dbClient.Collection(mongo.PbehaviorMongoCollection),
		alarmStore:             alarmStore,
		defaultSortBy:          "name",
		links:                  alarmapi.NewLinksFetcher(legacyURL, linkFetchTimeout),
		timezoneConfigProvider: timezoneConfigProvider,
	}
}

type store struct {
	dbCollection           mongo.DbCollection
	pbehaviorCollection    mongo.DbCollection
	links                  alarmapi.LinksFetcher
	alarmStore             alarmapi.Store
	defaultSortBy          string
	timezoneConfigProvider config.TimezoneConfigProvider
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	match := []bson.M{
		{"type": libtypes.EntityTypeService},
		{"$expr": bson.M{"$eq": bson.A{"$enabled", true}}},
	}
	if r.Category != "" {
		match = append(match, bson.M{"category": r.Category})
	}
	pipeline := []bson.M{
		{"$match": bson.M{"$and": match}},
	}
	pipeline = append(pipeline, getFindPipeline()...)
	parsedFilter, err := ParseFilter(r.Filter)
	if err != nil {
		return nil, err
	}
	if parsedFilter != nil {
		pipeline = append(pipeline, bson.M{"$match": parsedFilter})
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		s.getSort(r.SortBy, r.Sort),
	))

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

	pbhIDs := make([]string, 0)
	for _, v := range res.Data {
		if v.PbehaviorID != "" {
			pbhIDs = append(pbhIDs, v.PbehaviorID)
		}
	}

	pbhs, err := s.findPbehaviors(ctx, pbhIDs)
	if err != nil {
		return nil, err
	}

	for i := range res.Data {
		if v, ok := pbhs[res.Data[i].PbehaviorID]; ok {
			res.Data[i].Pbehaviors = []pbehavior.Response{v}
		} else {
			res.Data[i].Pbehaviors = []pbehavior.Response{}
		}
	}

	return &res, nil
}

func (s *store) FindEntities(ctx context.Context, id, apiKey string, query EntitiesListRequest) (*EntityAggregationResult, error) {
	var service libtypes.Entity
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": id, "enabled": true}).Decode(&service)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	location := s.timezoneConfigProvider.Get().Location

	pipeline := []bson.M{
		{"$match": bson.M{"$and": []bson.M{
			{"$expr": bson.M{"$in": bson.A{"$_id", service.Depends}}},
			{"$expr": bson.M{"$eq": bson.A{"$enabled", true}}},
		}}},
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

	pbhIDs := make([]string, 0)
	alarmIds := make([]string, 0, len(res.Data))
	for _, v := range res.Data {
		if v.PbehaviorInfo != nil && v.PbehaviorInfo.ID != "" {
			pbhIDs = append(pbhIDs, v.PbehaviorInfo.ID)
		}

		if v.AlarmID != "" {
			alarmIds = append(alarmIds, v.AlarmID)
		}
	}

	if query.WithInstructions {
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
				assignedInstructions = make([]alarmapi.InstructionWithAlarms, 0)
			}
			res.Data[idx].AssignedInstructions = assignedInstructions
			res.Data[idx].IsAutoInstructionRunning = statusesByAlarm[v.AlarmID].AutoRunning
			res.Data[idx].IsAllAutoInstructionsCompleted = statusesByAlarm[v.AlarmID].AutoAllCompleted
			res.Data[idx].IsAutoInstructionFailed = statusesByAlarm[v.AlarmID].AutoFailed
			res.Data[idx].IsManualInstructionWaitingResult = statusesByAlarm[v.AlarmID].ManualRunning
		}
	}

	pbhs, err := s.findPbehaviors(ctx, pbhIDs)
	if err != nil {
		return nil, err
	}

	for i := range res.Data {
		if res.Data[i].PbehaviorInfo != nil && !res.Data[i].PbehaviorInfo.IsActive() || !service.PbehaviorInfo.IsActive() {
			res.Data[i].IsGrey = true
		}

		res.Data[i].Pbehaviors = make([]pbehavior.Response, 0)

		if res.Data[i].PbehaviorInfo != nil {
			if v, ok := pbhs[res.Data[i].PbehaviorInfo.ID]; ok {
				res.Data[i].Pbehaviors = append(res.Data[i].Pbehaviors, v)
			}
		}
	}

	err = s.fillLinks(ctx, apiKey, &res)
	if err != nil {
		log.Printf("fillLinks error %s", err)
	}

	return &res, nil
}

func (s *store) findPbehaviors(ctx context.Context, ids []string) (map[string]pbehavior.Response, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	pipeline := []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": ids}}},
	}
	pipeline = append(pipeline, pbehavior.GetNestedObjectsPipeline()...)
	cursor, err := s.pbehaviorCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var pbhs []pbehavior.Response
	err = cursor.All(ctx, &pbhs)
	if err != nil {
		return nil, err
	}

	pbhsByID := make(map[string]pbehavior.Response, len(pbhs))
	for _, pbh := range pbhs {
		pbhsByID[pbh.ID] = pbh
	}

	return pbhsByID, nil
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
	if err != nil {
		return err
	}
	if res == nil {
		return nil
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
		sortBy = s.defaultSortBy
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

func getFindPipeline() []bson.M {
	pipeline := []bson.M{
		// Find category
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
		// Find pbehavior types
		{"$addFields": bson.M{
			"alarm_counters": bson.M{"$map": bson.M{
				"input": bson.M{"$objectToArray": "$alarms_cumulative_data.watched_pbehavior_count"},
				"as":    "each",
				"in": bson.M{
					"type":  "$$each.k",
					"count": "$$each.v",
				},
			}},
		}},
		{"$unwind": bson.M{"path": "$alarm_counters", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         pbehaviorlib.TypeCollectionName,
			"localField":   "alarm_counters.type",
			"foreignField": "_id",
			"as":           "alarm_counters.type",
		}},
		{"$unwind": bson.M{"path": "$alarm_counters.type", "preserveNullAndEmptyArrays": true}},
		{"$group": bson.M{
			"_id":            "$_id",
			"data":           bson.M{"$first": "$$ROOT"},
			"alarm_counters": bson.M{"$push": "$alarm_counters"},
		}},
		{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{"alarm_counters": bson.M{"$filter": bson.M{
					"input": "$alarm_counters",
					"cond":  bson.M{"$gt": bson.A{"$$this.count", 0}},
				}}},
			}},
		}},
		// Find service's alarm.
		{"$lookup": bson.M{
			"from": alarm.AlarmCollectionName,
			"let":  bson.M{"eid": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$and": []bson.M{
					{"$expr": bson.M{"$eq": bson.A{"$d", "$$eid"}}},
					// Get only open alarm.
					{"v.resolved": bson.M{"$exists": false}},
				}}},
				{"$limit": 1},
			},
			"as": "alarm",
		}},
		{"$unwind": bson.M{"path": "$alarm", "preserveNullAndEmptyArrays": true}},
		// Add alarm fields and computed fields to result
		{"$addFields": bson.M{
			"connector":        "$alarm.v.connector",
			"connector_name":   "$alarm.v.connector_name",
			"component":        "$alarm.v.component",
			"resource":         "$alarm.v.resource",
			"output":           "$alarm.v.output",
			"last_update_date": "$alarm.v.last_update_date",
			"state":            "$alarm.v.state",
			"status":           "$alarm.v.status",
			"snooze":           "$alarm.v.snooze",
			"ack":              "$alarm.v.ack",
			"pbehavior_id":     "$pbehavior_info.id",
			"depends_total":    bson.M{"$size": "$depends"},
			"impact_state":     bson.M{"$multiply": bson.A{"$alarm.v.state.val", "$impact_level"}},
			"has_open_alarm": bson.M{"$cond": bson.M{
				"if":   bson.M{"$gt": bson.A{"$alarms_cumulative_data.watched_not_acked_count", 0}},
				"then": true,
				"else": false,
			}},
			"watched_inactive_count": bson.M{"$sum": bson.M{
				"$map": bson.M{
					"input": bson.M{
						"$filter": bson.M{
							"input": "$alarm_counters",
							"cond":  bson.M{"$in": bson.A{"$$this.type.type", bson.A{pbehaviorlib.TypeMaintenance, pbehaviorlib.TypePause, pbehaviorlib.TypeInactive}}},
						},
					},
					"as": "each",
					"in": "$$each.count",
				},
			}},
			"watched_pbehavior_type": bson.M{"$switch": bson.M{
				"branches": []bson.M{
					{
						"case": bson.M{"$ne": bson.A{bson.A{}, bson.M{
							"$filter": bson.M{
								"input": "$alarm_counters",
								"cond":  bson.M{"$eq": bson.A{"$$this.type.type", pbehaviorlib.TypeMaintenance}},
							},
						}}},
						"then": pbehaviorlib.TypeMaintenance,
					},
					{
						"case": bson.M{"$ne": bson.A{bson.A{}, bson.M{
							"$filter": bson.M{
								"input": "$alarm_counters",
								"cond":  bson.M{"$eq": bson.A{"$$this.type.type", pbehaviorlib.TypePause}},
							},
						}}},
						"then": pbehaviorlib.TypePause,
					},
					{
						"case": bson.M{"$ne": bson.A{bson.A{}, bson.M{
							"$filter": bson.M{
								"input": "$alarm_counters",
								"cond":  bson.M{"$eq": bson.A{"$$this.type.type", pbehaviorlib.TypeInactive}},
							},
						}}},
						"then": pbehaviorlib.TypeInactive,
					},
				},
				"default": "",
			}},
			"links": bson.M{"$map": bson.M{
				"input": bson.M{"$objectToArray": "$links"},
				"as":    "each",
				"in": bson.M{
					"cat_name": "$$each.k",
					"links":    "$$each.v",
				},
			}},
		}},
	}
	pipeline = append(pipeline, getFindIconPipeline()...)

	return pipeline
}

func getFindIconPipeline() []bson.M {
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
					[]bson.M{
						// If service is not active return pbehavior type icon.
						{
							"case": bson.M{"$and": []bson.M{
								{"$ifNull": bson.A{"$pbehavior_info", false}},
								{"$ne": bson.A{"$pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
							}},
							"then": "$pbehavior_info.canonical_type",
						},
						// If all watched alarms are not active return most priority pbehavior type of watched alarms.
						{
							"case": bson.M{"$and": []bson.M{
								{"$gt": bson.A{"$watched_inactive_count", 0}},
								{"$eq": bson.A{"$watched_inactive_count", "$depends_total"}},
							}},
							"then": "$watched_pbehavior_type",
						},
					},
					// Else return state icon.
					stateVals...,
				),
				"default": defaultVal,
			}},
			"is_grey": bson.M{"$switch": bson.M{
				"branches": []bson.M{
					{
						"case": bson.M{"$and": []bson.M{
							{"$ifNull": bson.A{"$pbehavior_info", false}},
							{"$ne": bson.A{"$pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
						}},
						"then": true,
					},
					{
						"case": bson.M{"$and": []bson.M{
							{"$gt": bson.A{"$watched_inactive_count", 0}},
							{"$eq": bson.A{"$watched_inactive_count", "$depends_total"}},
						}},
						"then": true,
					},
				},
				"default": false,
			}},
			"secondary_icon": bson.M{"$switch": bson.M{
				"branches": []bson.M{
					{
						// If only some watched alarms are not active return most priority pbehavior type of watched alarms.
						"case": bson.M{"$and": []bson.M{
							{"$gt": bson.A{"$watched_inactive_count", 0}},
							{"$lt": bson.A{"$watched_inactive_count", "$depends_total"}},
						}},
						"then": "$watched_pbehavior_type",
					},
				},
				"default": "",
			}},
		}},
	}
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
		// Event statistics
		{"$lookup": bson.M{
			"from":       mongo.EventStatistics,
			"localField": "_id", "foreignField": "_id",
			"as": "stats",
		}},
		{"$unwind": bson.M{"path": "$stats", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			// stats counters with "last_event" prior "truncatedInLocation" represent as 0
			"stats.ok": bson.M{"$cond": bson.M{
				"if":   bson.M{"$gt": bson.A{"$stats.last_event", truncatedInLocation}},
				"then": "$stats.ok", "else": 0}},
			"stats.ko": bson.M{"$cond": bson.M{
				"if":   bson.M{"$gt": bson.A{"$stats.last_event", truncatedInLocation}},
				"then": "$stats.ko", "else": 0}},
		}},
		// Find connected alarm.
		{"$lookup": bson.M{
			"from": alarm.AlarmCollectionName,
			"let":  bson.M{"eid": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$d", "$$eid"}}}},
				{"$match": bson.M{"$or": []bson.M{
					{"v.resolved": bson.M{"$in": bson.A{false, nil}}},
					{"v.resolved": bson.M{"$exists": false}},
				}}},
				{"$sort": bson.M{"v.creation_date": -1}},
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
			"impact_state":     bson.M{"$multiply": bson.A{"$alarm.v.state.val", "$impact_level"}},
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
