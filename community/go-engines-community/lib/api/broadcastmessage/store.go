package broadcastmessage

import (
	"cmp"
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	GetById(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id, userId string) (bool, error)
	GetActive(ctx context.Context) ([]Response, error)
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	maintenanceAdapter    config.MaintenanceAdapter
	authorProvider        author.Provider
	defaultSearchByFields []string
	defaultSortBy         string
}

func (s store) Insert(ctx context.Context, r CreateRequest) (*Response, error) {
	now := datetime.NewCpsTime()

	r.ID = cmp.Or(r.ID, utils.NewID())
	r.Created = &now
	r.Updated = &now

	var resp *Response

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		resp = nil

		_, err := s.dbCollection.InsertOne(ctx, r)
		if err != nil {
			return err
		}

		resp, err = s.GetById(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s store) GetById(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		response := Response{}
		err := cursor.Decode(&response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	return nil, nil
}

func (s store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		common.GetSortQuery(cmp.Or(query.SortBy, s.defaultSortBy), query.Sort),
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

func (s store) Update(ctx context.Context, r UpdateRequest) (*Response, error) {
	now := datetime.NewCpsTime()

	r.Updated = &now

	var resp *Response

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		resp = nil

		_, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": r})
		if err != nil {
			return err
		}

		resp, err = s.GetById(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s store) Delete(ctx context.Context, id, userId string) (bool, error) {
	var deleted int64

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userId}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
}

func (s store) GetActive(ctx context.Context) ([]Response, error) {
	now := time.Now().Unix()

	conf, err := s.maintenanceAdapter.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"start": bson.M{
					"$lte": now,
				},
				"end": bson.M{
					"$gte": now,
				},
			},
		},
		{
			"$sort": bson.M{
				"start": -1,
				"_id":   1,
			},
		},
	}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	messages := make([]Response, 0)
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
	authorProvider author.Provider,
) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.BroadcastMessageMongoCollection),
		maintenanceAdapter:    maintenanceAdapter,
		authorProvider:        authorProvider,
		defaultSortBy:         "_id",
		defaultSearchByFields: []string{"_id", "message"},
	}
}
