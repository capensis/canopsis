package eventfilter

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(ctx context.Context, model *EventFilter) error
	GetById(ctx context.Context, id string) (*EventFilter, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, model *EventFilter) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type AggregationResult struct {
	Data       []*EventFilter `bson:"data" json:"data"`
	TotalCount int64          `bson:"total_count" json:"total_count"`
}

type store struct {
	dbCollection          mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		dbCollection:          dbClient.Collection(mongo.EventFilterRulesMongoCollection),
		defaultSearchByFields: []string{"_id", "author", "description", "type"},
		defaultSortBy:         "created",
	}
}

func (s *store) Insert(ctx context.Context, model *EventFilter) error {
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

func (s *store) GetById(ctx context.Context, id string) (*EventFilter, error) {
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

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := s.defaultSortBy
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

func (s *store) Update(ctx context.Context, model *EventFilter) (bool, error) {
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

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
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