package broadcastmessage

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(ctx context.Context, model *BroadcastMessage) error
	GetById(ctx context.Context, id string) (*BroadcastMessage, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, model *BroadcastMessage) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetActive(ctx context.Context) ([]*BroadcastMessage, error)
}

type store struct {
	dbCollection          mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

func (s store) Insert(ctx context.Context, model *BroadcastMessage) error {
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

func (s store) GetById(ctx context.Context, id string) (*BroadcastMessage, error) {
	bm := &BroadcastMessage{}
	d := s.dbCollection.FindOne(ctx, bson.M{"_id": id})
	if d.Err() != nil {
		return nil, d.Err()
	}
	if err := d.Decode(&bm); err != nil {
		return nil, err
	}
	return bm, nil
}

func (s store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
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

func (s store) Update(ctx context.Context, model *BroadcastMessage) (bool, error) {
	var data BroadcastMessage
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

func (s store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s store) GetActive(ctx context.Context) ([]*BroadcastMessage, error) {
	now := time.Now().Unix()
	cursor, err := s.dbCollection.Find(ctx, bson.M{"$and": bson.A{
		bson.M{
			"start": bson.M{
				"$lte": now,
			},
		},
		bson.M{
			"end": bson.M{
				"$gte": now,
			},
		},
	}})
	if err != nil {
		return nil, err
	}

	actives := make([]*BroadcastMessage, 0)
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var bm BroadcastMessage
		err = cursor.Decode(&bm)
		if err != nil {
			return nil, err
		}
		actives = append(actives, &bm)
	}
	return actives, nil
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		dbCollection:          dbClient.Collection(mongo.BroadcastMessageMongoCollection),
		defaultSortBy:         "_id",
		defaultSearchByFields: []string{"_id", "message"},
	}
}
