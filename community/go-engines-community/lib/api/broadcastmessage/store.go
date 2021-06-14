package broadcastmessage

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"go.mongodb.org/mongo-driver/bson"

	"git.canopsis.net/canopsis/go-engines/lib/mongo"
)

type Store interface {
	Insert(model *BroadcastMessage) error
	GetById(id string) (*BroadcastMessage, error)
	Find(query FilteredQuery) (*AggregationResult, error)
	Update(model *BroadcastMessage) (bool, error)
	Delete(id string) (bool, error)
	GetActive() ([]*BroadcastMessage, error)
}

type AggregationResult struct {
	Data       []*BroadcastMessage `bson:"data" json:"data"`
	TotalCount int64               `bson:"total_count" json:"total_count"`
}

type store struct {
	dbCollection mongo.DbCollection
}

func (s store) Insert(model *BroadcastMessage) error {
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

func (s store) GetById(id string) (*BroadcastMessage, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

func (s store) Find(query FilteredQuery) (*AggregationResult, error) {
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
				{"message": searchRegexp},
			},
		}
	} else {
		filter = bson.M{}
	}
	pipeline := []bson.M{
		{"$match": filter},
	}

	sortBy := "_id"
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

func (s store) Update(model *BroadcastMessage) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (s store) Delete(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s store) GetActive() ([]*BroadcastMessage, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
		dbCollection: dbClient.Collection(mongo.BroadcastMessageMongoCollection),
	}
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
