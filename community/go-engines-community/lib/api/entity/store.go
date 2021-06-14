package entity

import (
	"context"
	"encoding/json"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/expression/parser"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(r ListRequestWithPagination) (*AggregationResult, error)
}

type store struct {
	db                    mongo.DbClient
	dbCollection          mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

func NewStore(db mongo.DbClient) Store {
	return &store{
		db:           db,
		dbCollection: db.Collection(mongo.EntityMongoCollection),
		defaultSearchByFields: []string{
			"name", "type",
		},
		defaultSortBy: "name",
	}
}

func (s *store) Find(r ListRequestWithPagination) (*AggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pipeline := make([]bson.M, 0)
	err := s.addFilter(r.ListRequest, &pipeline)
	if err != nil {
		return nil, err
	}
	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		s.getSort(r.ListRequest),
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

func (s *store) addFilter(r ListRequest, pipeline *[]bson.M) error {
	match := make([]bson.M, 0)
	err := s.addQueryFilter(r, &match)
	if err != nil {
		return err
	}

	s.addSearchFilter(r, &match)

	if len(match) > 0 {
		*pipeline = append(*pipeline, bson.M{"$match": bson.M{"$and": match}})
	}

	return nil
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

	sortDir := 1
	if r.Sort == common.SortDesc {
		sortDir = -1
	}

	return bson.M{"$sort": bson.M{sortBy: sortDir}}
}
