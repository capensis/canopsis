package eventfilter

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/api/common"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"

	"git.canopsis.net/canopsis/go-engines/lib/utils"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(model *EventFilter) error
	GetById(id string) (*EventFilter, error)
	Find(query FilteredQuery) (*AggregationResult, error)
	Update(model *EventFilter) (bool, error)
	Delete(id string) (bool, error)
}

type AggregationResult struct {
	Data       []*EventFilter `bson:"data" json:"data"`
	TotalCount int64          `bson:"total_count" json:"total_count"`
}

type store struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.EventFilterRulesMongoCollection),
	}
}

func (s *store) Insert(model *EventFilter) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if model.ID == "" {
		model.ID = utils.NewID()
	}
	now := types.NewCpsTime(time.Now().Unix())
	model.Created = &now
	model.Updated = &now

	_, err := s.dbCollection.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	return err
}

func (s *store) GetById(id string) (*EventFilter, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ef := &EventFilter{}
	d := s.dbCollection.FindOne(ctx, bson.M{"_id": id})
	if d.Err() != nil {
		return nil, d.Err()
	}
	if err := d.Decode(&ef); err != nil {
		return nil, err
	}
	return ef, nil
}

func (s *store) Find(query FilteredQuery) (*AggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var filter bson.M

	if query.Search != "" {
		searchRegexp := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", query.Search),
			Options: "i",
		}

		filter = bson.M{
			"$or": []bson.M{
				{"_id": searchRegexp},
				{"author": searchRegexp},
				{"description": searchRegexp},
				{"type": searchRegexp},
			},
		}
	} else {
		filter = bson.M{}
	}
	pipeline := []bson.M{
		{"$match": filter},
	}

	sortBy := "created"
	if query.SortBy != "" {
		sortBy = query.SortBy
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		common.GetSortQuery(sortBy, query.Sort),
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

func (s *store) Update(model *EventFilter) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var data EventFilter
	updated := types.NewCpsTime(time.Now().Unix())
	model.Created = nil
	model.Updated = &updated
	err := s.dbCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": model.ID},
		bson.M{"$set": model},
	).Decode(&data)
	model.Created = data.Created
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) Delete(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
