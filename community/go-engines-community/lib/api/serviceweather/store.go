package serviceweather

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	alarmapi "git.canopsis.net/canopsis/go-engines/lib/api/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/entityservice"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	pbehaviorlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	libtypes "git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"go.mongodb.org/mongo-driver/bson"
)

const linkFetchTimeout = 30 * time.Second

type Store interface {
	Find(context.Context, ListRequest) (*AggregationResult, error)
	FindEntities(ctx context.Context, id, apiKey string, query EntitiesListRequest) (*EntityAggregationResult, error)
}

func NewStore(
	dbClient mongo.DbClient,
	legacyURL fmt.Stringer,
	statsStore StatsStore,
	timezoneConfigProvider config.TimezoneConfigProvider,
	pbhStore redis.Store,
	pbhService pbehaviorlib.Service,
) Store {
	return &store{
		dbCollection:           dbClient.Collection(mongo.EntityMongoCollection),
		pbehaviorCollection:    dbClient.Collection(mongo.PbehaviorMongoCollection),
		statsStore:             statsStore,
		defaultSortBy:          "name",
		links:                  alarmapi.NewLinksFetcher(legacyURL, linkFetchTimeout),
		timezoneConfigProvider: timezoneConfigProvider,
		pbhStore:               pbhStore,
		pbhService:             pbhService,
	}
}

type store struct {
	dbCollection           mongo.DbCollection
	pbehaviorCollection    mongo.DbCollection
	links                  alarmapi.LinksFetcher
	statsStore             StatsStore
	defaultSortBy          string
	timezoneConfigProvider config.TimezoneConfigProvider
	pbhStore               redis.Store
	pbhService             pbehaviorlib.Service
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
			res.Data[i].Pbehaviors = []pbehavior.PBehavior{v}
		} else {
			res.Data[i].Pbehaviors = []pbehavior.PBehavior{}
		}
	}

	return &res, nil
}

func (s *store) FindEntities(ctx context.Context, id, apiKey string, query EntitiesListRequest) (*EntityAggregationResult, error) {
	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id}},
		// Find category
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
		{"$project": bson.M{
			"entity": "$$ROOT",
		}},
		{"$lookup": bson.M{
			"from": mongo.AlarmMongoCollection,
			"let":  bson.M{"eid": "$entity._id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$$eid", "$d"}}}},
				{"$match": bson.M{"$or": []bson.M{
					{"v.resolved": bson.M{"$in": bson.A{false, nil}}},
					{"v.resolved": bson.M{"$exists": false}},
				}}},
			},
			"as": "alarm",
		}},
		{"$unwind": bson.M{"path": "$alarm", "preserveNullAndEmptyArrays": true}},
	})
	if err != nil {
		return nil, err
	}

	var service entityservice.AlarmWithEntity
	if cursor.Next(ctx) {
		err := cursor.Decode(&service)
		if err != nil {
			cursor.Close(ctx)
			return nil, err
		}
	}
	cursor.Close(ctx)

	if service.Entity.ID == "" {
		return nil, nil
	}

	filter := bson.M{"$and": []bson.M{
		{"$expr": bson.M{"$in": bson.A{"$_id", service.Entity.Depends}}},
		{"$expr": bson.M{"$eq": bson.A{"$enabled", true}}},
	}}
	pipeline := []bson.M{
		{"$match": filter},
	}
	pipeline = append(pipeline, getFindEntitiesPipeline()...)
	cursor, err = s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
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
	entitiesWithoutPbh := make([]int, 0, len(res.Data))
	for idx, v := range res.Data {
		if v.PbehaviorInfo.ID != "" {
			pbhIDs = append(pbhIDs, v.PbehaviorInfo.ID)
		}

		if !v.PbehaviorInfo.IsActive() || service.Alarm != nil && !service.Alarm.Value.PbehaviorInfo.IsActive() {
			res.Data[idx].IsGrey = true
		}

		if v.PbehaviorInfo.ID == "" && !res.Data[idx].IsGrey {
			entitiesWithoutPbh = append(entitiesWithoutPbh, idx)
		}
	}

	pbhRestored, err := s.pbhStore.Restore(ctx, s.pbhService)
	if err == nil && pbhRestored {

		now := time.Now().In(s.timezoneConfigProvider.Get().Location)

		for _, idx := range entitiesWithoutPbh {
			infos := make(map[string]libtypes.Info, len(res.Data[idx].Infos))
			for k, v := range res.Data[idx].Infos {
				infos[k] = libtypes.Info{
					Description: v.Description,
					Value:       v.Value,
				}
			}
			entity := libtypes.NewEntity(res.Data[idx].ID, res.Data[idx].Name, res.Data[idx].Type, infos, nil, nil)
			result, err := s.pbhService.Resolve(ctx, &entity, now)
			if err != nil && !errors.Is(err, pbehaviorlib.ErrRecomputeNeed) {
				return nil, err
			}
			if result.ResolvedPbhID != "" {
				res.Data[idx].PbehaviorInfo = libtypes.PbehaviorInfo{
					ID:            result.ResolvedPbhID,
					Name:          result.ResolvedPbhName,
					Reason:        result.ResolvedPbhReason,
					TypeID:        result.ResolvedType.ID,
					TypeName:      result.ResolvedType.Name,
					CanonicalType: result.ResolvedType.Type,
				}
				res.Data[idx].IsGrey = true
				pbhIDs = append(pbhIDs, result.ResolvedPbhID)
			}
		}

	}

	pbhs, err := s.findPbehaviors(ctx, pbhIDs)
	if err != nil {
		return nil, err
	}

	for i := range res.Data {
		if v, ok := pbhs[res.Data[i].PbehaviorInfo.ID]; ok {
			res.Data[i].Pbehaviors = []pbehavior.PBehavior{v}
		} else {
			res.Data[i].Pbehaviors = []pbehavior.PBehavior{}
		}

		if s.statsStore != nil {
			res.Data[i].Stats, err = s.statsStore.GetStats(ctx, res.Data[i].ID, s.timezoneConfigProvider.Get().Location)
			if err != nil {
				return nil, err
			}
		}
	}

	err = s.fillLinks(ctx, apiKey, &res)
	if err != nil {
		log.Printf("fillLinks error %s", err)
	}

	return &res, nil
}

func (s *store) findPbehaviors(ctx context.Context, ids []string) (map[string]pbehavior.PBehavior, error) {
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

	var pbhs []pbehavior.PBehavior
	err = cursor.All(ctx, &pbhs)
	if err != nil {
		return nil, err
	}

	pbhsByID := make(map[string]pbehavior.PBehavior, len(pbhs))
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
			"pbehavior_id":     "$alarm.v.pbehavior_info.id",
			"alarms_total":     bson.M{"$sum": "$alarms_cumulative_data.watched_count"},
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
								{"$ifNull": bson.A{"$alarm.v.pbehavior_info", false}},
								{"$ne": bson.A{"$alarm.v.pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
							}},
							"then": "$alarm.v.pbehavior_info.canonical_type",
						},
						// If all watched alarms are not active return most priority pbehavior type of watched alarms.
						{
							"case": bson.M{"$and": []bson.M{
								{"$gt": bson.A{"$alarms_total", 0}},
								{"$eq": bson.A{"$watched_inactive_count", "$alarms_total"}},
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
							{"$ifNull": bson.A{"$alarm.v.pbehavior_info", false}},
							{"$ne": bson.A{"$alarm.v.pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
						}},
						"then": true,
					},
					{
						"case": bson.M{"$and": []bson.M{
							{"$gt": bson.A{"$alarms_total", 0}},
							{"$eq": bson.A{"$watched_inactive_count", "$alarms_total"}},
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
							{"$gt": bson.A{"$alarms_total", 0}},
							{"$lt": bson.A{"$watched_inactive_count", "$alarms_total"}},
						}},
						"then": "$watched_pbehavior_type",
					},
				},
				"default": "",
			}},
		}},
	}
}

func getFindEntitiesPipeline() []bson.M {
	pipeline := []bson.M{
		// Find category
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
		// Find connected alarm.
		{"$lookup": bson.M{
			"from": alarm.AlarmCollectionName,
			"let":  bson.M{"eid": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$d", "$$eid"}}}},
				{"$sort": bson.M{"v.creation_date": -1}},
				{"$limit": 1},
			},
			"as": "alarm",
		}},
		{"$unwind": bson.M{"path": "$alarm", "preserveNullAndEmptyArrays": true}},
		// Add alarms fields to result
		{"$addFields": bson.M{
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
			"pbehavior_info":   "$alarm.v.pbehavior_info",
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
							{"$ifNull": bson.A{"$alarm.v.pbehavior_info", false}},
							{"$ne": bson.A{"$alarm.v.pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
						}},
						"then": "$alarm.v.pbehavior_info.canonical_type",
					}},
					// Else return state icon.
					stateVals...,
				),
				"default": defaultVal,
			}},
		}},
	}
}
