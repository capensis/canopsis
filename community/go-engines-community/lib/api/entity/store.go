package entity

import (
	"context"
	"encoding/json"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/expression/parser"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const bulkMaxSize = 10000

type Store interface {
	Find(ctx context.Context, r ListRequestWithPagination) (*AggregationResult, error)
	ArchiveDisabledEntities(ctx context.Context, archiveDeps bool) (int64, error)
	DeleteArchivedEntities(ctx context.Context) (int64, error)
}

type store struct {
	db                    mongo.DbClient
	mainCollection        mongo.DbCollection
	archivedCollection    mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

func NewStore(db mongo.DbClient) Store {
	return &store{
		db:                 db,
		mainCollection:     db.Collection(mongo.EntityMongoCollection),
		archivedCollection: db.Collection(mongo.ArchivedEntitiesMongoCollection),
		defaultSearchByFields: []string{
			"name", "type",
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

	project := []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
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

func (s *store) ArchiveDisabledEntities(ctx context.Context, archiveDeps bool) (int64, error) {
	var totalArchived int64

	archived, err := s.archiveEntitiesByType(ctx, types.EntityTypeConnector, archiveDeps)
	if err != nil {
		return 0, err
	}

	totalArchived += archived

	archived, err = s.archiveEntitiesByType(ctx, types.EntityTypeComponent, archiveDeps)
	if err != nil {
		return 0, err
	}

	totalArchived += archived

	archived, err = s.archiveEntitiesByType(ctx, types.EntityTypeResource, archiveDeps)
	if err != nil {
		return 0, err
	}

	totalArchived += archived

	return totalArchived, nil
}

func (s *store) archiveEntitiesByType(ctx context.Context, eType string, archiveDeps bool) (int64, error) {
	models := make([]mongodriver.WriteModel, 0, bulkMaxSize)
	ids := make([]string, 0, bulkMaxSize)

	cursor, err := s.mainCollection.Find(
		ctx,
		bson.M{
			"enabled": bson.M{"$in": bson.A{false, nil}},
			"type":    eType,
		},
	)
	if err != nil {
		return 0, err
	}

	totalArchived, err := s.processCursor(ctx, cursor, &models, &ids, archiveDeps)
	if err != nil {
		return 0, err
	}

	if len(models) != 0 {
		archived, err := s.bulkArchive(ctx, &models, &ids)
		if err != nil {
			return 0, err
		}

		totalArchived += archived
	}

	return totalArchived, cursor.Close(ctx)
}

func (s *store) processCursor(ctx context.Context, cursor mongo.Cursor, models *[]mongodriver.WriteModel, ids *[]string, archiveDeps bool) (int64, error) {
	var totalArchived int64

	for cursor.Next(ctx) {
		var entity types.Entity

		err := cursor.Decode(&entity)
		if err != nil {
			return 0, err
		}

		if archiveDeps {
			switch entity.Type {
			case types.EntityTypeConnector:
				archived, err := s.archiveDependencies(ctx, append(entity.Impacts, entity.Depends...), models, ids)
				if err != nil {
					return 0, err
				}

				totalArchived += archived
			case types.EntityTypeComponent:
				archived, err := s.archiveDependencies(ctx, entity.Depends, models, ids)
				if err != nil {
					return 0, err
				}

				totalArchived += archived
			case types.EntityTypeResource:
			default:
				continue
			}
		}

		*ids = append(*ids, entity.ID)
		*models = append(
			*models,
			mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": entity.ID}).
				SetUpdate(bson.M{"$set": entity}).
				SetUpsert(true),
		)

		if len(*models) == bulkMaxSize {
			archived, err := s.bulkArchive(ctx, models, ids)
			if err != nil {
				return 0, err
			}

			totalArchived += archived
		}
	}

	return totalArchived, nil
}

func (s *store) archiveDependencies(ctx context.Context, depIds []string, models *[]mongodriver.WriteModel, ids *[]string) (int64, error) {
	cursor, err := s.mainCollection.Find(
		ctx,
		bson.M{"_id": bson.M{"$in": depIds}},
	)
	if err != nil {
		return 0, err
	}

	// do not archive dependencies' dependencies
	archived, err := s.processCursor(ctx, cursor, models, ids, false)
	if err != nil {
		return 0, err
	}

	return archived, cursor.Close(ctx)
}

func (s *store) bulkArchive(ctx context.Context, models *[]mongodriver.WriteModel, ids *[]string) (int64, error) {
	res, err := s.archivedCollection.BulkWrite(ctx, *models)
	if err != nil {
		return 0, err
	}

	_, err = s.mainCollection.DeleteMany(
		ctx,
		bson.M{"_id": bson.M{"$in": ids}},
	)
	if err != nil {
		return 0, err
	}

	*models = (*models)[:0]
	*ids = (*ids)[:0]

	return res.UpsertedCount, err
}

func (s *store) DeleteArchivedEntities(ctx context.Context) (int64, error) {
	deleted, err := s.archivedCollection.DeleteMany(ctx, bson.M{})

	return deleted, err
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
	if r.Search == "" {
		return
	}

	p := parser.NewParser()
	expr, err := p.Parse(r.Search)
	if err == nil {
		query := expr.Query()
		*match = append(*match, query)

		return
	}

	searchRegexp := primitive.Regex{
		Pattern: fmt.Sprintf(".*%s.*", r.Search),
		Options: "i",
	}

	fields := r.SearchBy
	if len(fields) == 0 {
		fields = s.defaultSearchByFields
	}

	searchMatch := make([]bson.M, len(fields))
	for i := range fields {
		searchMatch[i] = bson.M{fields[i]: searchRegexp}
	}

	*match = append(*match, bson.M{
		"$or": searchMatch,
	})
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
			"from":         mongo.AlarmMongoCollection,
			"localField":   "_id",
			"foreignField": "d",
			"as":           "alarms",
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
