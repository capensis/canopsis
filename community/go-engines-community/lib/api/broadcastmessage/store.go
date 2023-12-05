package broadcastmessage

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Insert(ctx context.Context, model *BroadcastMessage) error
	GetById(ctx context.Context, id string) (*BroadcastMessage, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, model *BroadcastMessage) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetActive(ctx context.Context) ([]BroadcastMessage, error)
}

type store struct {
	dbCollection          mongo.DbCollection
	maintenanceAdapter    config.MaintenanceAdapter
	defaultSearchByFields []string
	defaultSortBy         string
}

func (s store) Insert(ctx context.Context, model *BroadcastMessage) error {
	if model.ID == "" {
		model.ID = utils.NewID()
	}

	now := datetime.NewCpsTime()
	model.Created = &now
	model.Updated = &now

	_, err := s.dbCollection.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	return err
}

func (s store) GetById(ctx context.Context, id string) (*BroadcastMessage, error) {
	bm := BroadcastMessage{}
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&bm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &bm, nil
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

	updated := datetime.NewCpsTime()
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

func (s store) GetActive(ctx context.Context) ([]BroadcastMessage, error) {
	now := time.Now().Unix()

	conf, err := s.maintenanceAdapter.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	cursor, err := s.dbCollection.Find(ctx, bson.M{
		"start": bson.M{
			"$lte": now,
		},
		"end": bson.M{
			"$gte": now,
		},
	}, options.Find().SetSort(bson.D{{Key: "start", Value: -1}, {Key: "_id", Value: 1}}))
	if err != nil {
		return nil, err
	}

	messages := make([]BroadcastMessage, 0)
	err = cursor.All(ctx, &messages)
	if err != nil {
		return nil, err
	}

	if conf.Enabled {
		for idx := range messages {
			if messages[idx].ID == conf.BroadcastID {
				messages[idx].Maintenance = true
				break
			}
		}
	}

	return messages, nil
}

func NewStore(
	dbClient mongo.DbClient,
	maintenanceAdapter config.MaintenanceAdapter,
) Store {
	return &store{
		dbCollection:          dbClient.Collection(mongo.BroadcastMessageMongoCollection),
		maintenanceAdapter:    maintenanceAdapter,
		defaultSortBy:         "_id",
		defaultSearchByFields: []string{"_id", "message"},
	}
}
