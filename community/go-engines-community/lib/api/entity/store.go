package entity

import (
	"context"
	"encoding/json"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(ctx context.Context, r ListRequestWithPagination) (*AggregationResult, error)
	Toggle(ctx context.Context, id string, enabled bool) (bool, SimplifiedEntity, error)
}

type store struct {
	db                     mongo.DbClient
	mainCollection         mongo.DbCollection
	archivedCollection     mongo.DbCollection
	timezoneConfigProvider config.TimezoneConfigProvider
	defaultSearchByFields  []string
	defaultSortBy          string
}

func NewStore(db mongo.DbClient, timezoneConfigProvider config.TimezoneConfigProvider) Store {
	return &store{
		db:                     db,
		mainCollection:         db.Collection(mongo.EntityMongoCollection),
		archivedCollection:     db.Collection(mongo.ArchivedEntitiesMongoCollection),
		timezoneConfigProvider: timezoneConfigProvider,
		defaultSearchByFields: []string{
			"_id", "name", "type",
		},
		defaultSortBy: "_id",
	}
}

func (s *store) Find(ctx context.Context, r ListRequestWithPagination) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	err := s.addFilter(r.ListRequest, &pipeline)
	if err != nil {
		return nil, err
	}

	location := s.timezoneConfigProvider.Get().Location

	year, month, day := time.Now().In(location).Date()
	truncatedInLocation := time.Date(year, month, day, 0, 0, 0, 0, location)

	project := []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from": mongo.EventStatistics,
			"let":  bson.M{"id": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$and": []bson.M{
					{"$expr": bson.M{"$eq": bson.A{"$_id", "$$id"}}},
					{"last_event": bson.M{"$gt": truncatedInLocation.Unix()}},
				}}},
			},
			"as": "eventStatistics",
		}},
		{"$unwind": bson.M{"path": "$eventStatistics", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"ok_events": "$eventStatistics.ok",
			"ko_events": "$eventStatistics.ko",
		}},
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
		{"$addFields": bson.M{"state": "$alarm.v.state.val"}},
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"foreignField": "_id",
			"localField":   "pbehavior_info.type",
			"as":           "pbehavior_type",
		}},
		{"$unwind": bson.M{"path": "$pbehavior_type", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"pbehavior_info": bson.M{"$cond": bson.M{
				"if": "$pbehavior_info",
				"then": bson.M{"$mergeObjects": bson.A{
					"$pbehavior_info",
					bson.M{"icon_name": "$pbehavior_type.icon_name"},
				}},
				"else": nil,
			}},
		}},
	}
	if r.NoEvents {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"idle_since": bson.M{"$gt": 0}}})
	}
	if r.WithFlags {
		project = append(project, getDeletablePipeline()...)
	}

	cursor, err := s.mainCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		s.getSort(r.ListRequest),
		project,
	), options.Aggregate().SetAllowDiskUse(true))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	res := AggregationResult{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *store) Toggle(ctx context.Context, id string, enabled bool) (bool, SimplifiedEntity, error) {
	var isToggled bool
	var oldSimplifiedEntity SimplifiedEntity

	err := s.db.WithTransaction(ctx, func(ctx context.Context) error {
		isToggled = false
		oldSimplifiedEntity = SimplifiedEntity{}

		err := s.mainCollection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": id},
			bson.M{"$set": bson.M{"enabled": enabled}},
			options.
				FindOneAndUpdate().
				SetProjection(bson.M{"_id": 1, "enabled": 1, "type": 1}).
				SetReturnDocument(options.Before),
		).Decode(&oldSimplifiedEntity)
		if err != nil {
			return err
		}

		isToggled = oldSimplifiedEntity.Enabled != enabled
		return nil
	})

	return isToggled, oldSimplifiedEntity, err
}

func (s *store) addFilter(r ListRequest, pipeline *[]bson.M) error {
	match := make([]bson.M, 0)
	err := s.addQueryFilter(r, &match)
	if err != nil {
		return err
	}

	s.addCategoryFilter(r, &match)
	s.addSearchFilter(r, &match)

	if len(match) > 0 {
		*pipeline = append(*pipeline, bson.M{"$match": bson.M{"$and": match}})
	}

	return nil
}

func (s *store) addCategoryFilter(r ListRequest, match *[]bson.M) {
	if r.Category == "" {
		return
	}

	*match = append(*match, bson.M{"category": r.Category})
}

func (s *store) addQueryFilter(r ListRequest, match *[]bson.M) error {
	if r.Filter == "" {
		return nil
	}

	var queryFilter bson.M
	err := json.Unmarshal([]byte(r.Filter), &queryFilter)
	if err != nil {
		return err
	}

	*match = append(*match, queryFilter)
	return nil
}

func (s *store) addSearchFilter(r ListRequest, match *[]bson.M) {
	searchBy := r.SearchBy
	if len(searchBy) == 0 {
		searchBy = s.defaultSearchByFields
	}

	query := common.GetSearchQuery(r.Search, searchBy)
	if len(query) != 0 {
		*match = append(*match, query)
	}
}

func (s *store) getSort(r ListRequest) bson.M {
	sortBy := r.SortBy
	if sortBy == "" {
		sortBy = s.defaultSortBy
	}

	return common.GetSortQuery(sortBy, r.Sort)
}

func getDeletablePipeline() []bson.M {
	return []bson.M{
		// Entity can be deleted if entity is service or if there aren't any alarm which is related to entity.
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
			"as": "alarms",
		}},
		{"$addFields": bson.M{
			"deletable": bson.M{"$cond": bson.M{
				"if": bson.M{"$or": []bson.M{
					{"$eq": bson.A{"$type", types.EntityTypeService}},
					{"$eq": bson.A{"$alarms", bson.A{}}},
				}},
				"then": true,
				"else": false,
			}},
		}},
		{"$project": bson.M{"alarms": 0}},
	}
}
