package watcherweather

import (
	"context"
	"fmt"
	"log"
	"time"

	alarmapi "git.canopsis.net/canopsis/go-engines/lib/api/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	libentity "git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	pbehaviorlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	libtypes "git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const colorPause = "pause"
const linkFetchTimeout = 30 * time.Second

type Store interface {
	Find(query ListRequest) (*AggregationResult, error)
	FindEntities(id, apiKey string, query ListRequest) (*EntityAggregationResult, error)
}

func NewStore(
	dbClient mongo.DbClient,
	legacyURL fmt.Stringer,
	statsStore StatsStore,
	location *time.Location,
) Store {
	s := &store{
		dbCollection:        dbClient.Collection(libentity.EntityCollectionName),
		pbehaviorCollection: dbClient.Collection(pbehaviorlib.PBehaviorCollectionName),
		statsStore:          statsStore,
		location:            location,
		defaultSortBy:       "name",
	}
	s.links = alarmapi.NewLinksFetcher(legacyURL, linkFetchTimeout)
	return s
}

type store struct {
	dbCollection        mongo.DbCollection
	pbehaviorCollection mongo.DbCollection
	links               alarmapi.LinksFetcher
	location            *time.Location
	statsStore          StatsStore
	defaultSortBy       string
}

func (s *store) Find(query ListRequest) (*AggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pipeline := []bson.M{
		{"$match": bson.M{"$and": []bson.M{
			{"type": libtypes.EntityTypeWatcher},
			{"$expr": bson.M{"$eq": bson.A{"$enabled", true}}},
		}}},
	}
	pipeline = append(pipeline, getFindPipeline()...)
	parsedFilter, err := ParseFilter(query.Filter)
	if err != nil {
		return nil, err
	}
	if parsedFilter != nil {
		pipeline = append(pipeline, bson.M{"$match": parsedFilter})
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		s.getSort(query.SortBy, query.Sort),
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

	pbhs, err := s.findPbehaviors(pbhIDs)
	if err != nil {
		return nil, err
	}

	for i := range res.Data {
		if v, ok := pbhs[res.Data[i].PbehaviorID]; ok {
			res.Data[i].Pbehaviors = []pbehavior.PBehavior{v}
		} else {
			res.Data[i].Pbehaviors = []pbehavior.PBehavior{}
		}
		// Keep for compatibility.
		res.Data[i].IsActionRequired = res.Data[i].HasOpenAlarm
		res.Data[i].WatcherPbehaviors = res.Data[i].Pbehaviors
	}

	return &res, nil
}

func (s *store) FindEntities(id, apiKey string, query ListRequest) (*EntityAggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcherRes := s.dbCollection.FindOne(ctx, bson.M{"_id": id})
	if err := watcherRes.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	var watcher libtypes.Entity
	err := watcherRes.Decode(&watcher)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"$and": []bson.M{
		{"$expr": bson.M{"$in": bson.A{"$_id", watcher.Depends}}},
		{"$expr": bson.M{"$eq": bson.A{"$enabled", true}}},
	}}
	pipeline := []bson.M{
		{"$match": filter},
	}
	pipeline = append(pipeline, getFindEntitiesPipeline()...)
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
	for _, v := range res.Data {
		if v.PbehaviorID != "" {
			pbhIDs = append(pbhIDs, v.PbehaviorID)
		}
	}

	pbhs, err := s.findPbehaviors(pbhIDs)
	if err != nil {
		return nil, err
	}

	for i := range res.Data {
		if v, ok := pbhs[res.Data[i].PbehaviorID]; ok {
			res.Data[i].Pbehaviors = []pbehavior.PBehavior{v}
		} else {
			res.Data[i].Pbehaviors = []pbehavior.PBehavior{}
		}
		// Keep for compatibility.
		res.Data[i].EntityPbehaviors = res.Data[i].Pbehaviors
		if !res.Data[i].IsInactive && s.statsStore != nil {
			res.Data[i].Stats, err = s.statsStore.GetStats(ctx, res.Data[i].ID)
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

func (s *store) findPbehaviors(ids []string) (map[string]pbehavior.PBehavior, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
		// Find watcher's alarm.
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
			// Keep for compatibility.
			"entity_id":    "$_id",
			"display_name": "$name",
		}},
	}
	pipeline = append(pipeline, getFindColorAndIconPipeline()...)

	return pipeline
}

func getFindColorAndIconPipeline() []bson.M {
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
			"color": bson.M{"$switch": bson.M{
				"branches": append(
					[]bson.M{
						// If entity is not active return pause color.
						{
							"case": bson.M{"$and": []bson.M{
								{"$ifNull": bson.A{"$alarm.v.pbehavior_info", false}},
								{"$ne": bson.A{"$alarm.v.pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
							}},
							"then": colorPause,
						},
						// If all watched alarms are not active return pause color.
						{
							"case": bson.M{"$and": []bson.M{
								{"$gt": bson.A{"$alarms_total", 0}},
								{"$eq": bson.A{"$watched_inactive_count", "$alarms_total"}},
							}},
							"then": colorPause,
						},
					},
					// Else return state color.
					stateVals...,
				),
				"default": defaultVal,
			}},
			"icon": bson.M{"$switch": bson.M{
				"branches": append(
					[]bson.M{
						// If watcher is not active return pbehavior type icon.
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
		// Keep for compatibility.
		{"$addFields": bson.M{
			"is_all_watched_inactive": bson.M{"$eq": bson.A{"$watched_inactive_count", "$alarms_total"}},
			"is_inactive": bson.M{"$and": []bson.M{
				{"$ifNull": bson.A{"$alarm.v.pbehavior_info", false}},
				{"$ne": bson.A{"$alarm.v.pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
			}},
			"tileIcon":           "$icon",
			"tileColor":          "$color",
			"tileSecondary_icon": "$secondary_icon",
		}},
	}
}

func getFindEntitiesPipeline() []bson.M {
	pipeline := []bson.M{
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
			"pbehavior_id":     "$alarm.v.pbehavior_info.id",
			// Keep for compatibility.
			"entity_id": "$_id",
			"org":       "$infos.org",
			"is_inactive": bson.M{"$and": []bson.M{
				{"$ifNull": bson.A{"$alarm.v.pbehavior_info", false}},
				{"$ne": bson.A{"$alarm.v.pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
			}},
			"last_pbhleave_date": bson.M{"$max": bson.M{
				"$map": bson.M{
					"input": bson.M{"$filter": bson.M{
						"input": "$alarm.v.steps",
						"cond":  bson.M{"$eq": bson.A{"$$this._t", libtypes.AlarmStepPbhLeave}},
					}},
					"as": "each",
					"in": "$$each.t",
				},
			}},
		}},
	}
	pipeline = append(pipeline, getFindEntitiesColorAndIconPipeline()...)

	return pipeline
}

func getFindEntitiesColorAndIconPipeline() []bson.M {
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
			"color": bson.M{"$switch": bson.M{
				"branches": append(
					// If entity is not active return pause color.
					[]bson.M{{
						"case": bson.M{"$and": []bson.M{
							{"$ifNull": bson.A{"$alarm.v.pbehavior_info", false}},
							{"$ne": bson.A{"$alarm.v.pbehavior_info.canonical_type", pbehaviorlib.TypeActive}},
						}},
						"then": colorPause,
					}},
					// Else return state color.
					stateVals...,
				),
				"default": defaultVal,
			}},
			"icon": bson.M{"$switch": bson.M{
				"branches": append(
					// If watcher is not active return pbehavior type icon.
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
