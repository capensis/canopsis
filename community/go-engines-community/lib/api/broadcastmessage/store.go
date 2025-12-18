package broadcastmessage

import (
	"cmp"
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	GetByID(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query ListRequest) (*AggregationResult, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
	GetActive(ctx context.Context, userID string) ([]Response, error)
	Read(ctx context.Context, id, userID string) (bool, error)
}

func NewStore(
	dbClient mongo.DbClient,
	maintenanceAdapter config.MaintenanceAdapter,
	authorProvider author.Provider,
) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.BroadcastMessageCollection),
		dbReadCollection:      dbClient.Collection(mongo.BroadcastMessageReadCollection),
		maintenanceAdapter:    maintenanceAdapter,
		authorProvider:        authorProvider,
		defaultSortBy:         "_id",
		defaultSearchByFields: []string{"_id", "message"},
		dupErrorParser:        validation.NewDuplicateErrorParser(),
	}
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	dbReadCollection      mongo.DbCollection
	maintenanceAdapter    config.MaintenanceAdapter
	authorProvider        author.Provider
	defaultSearchByFields []string
	defaultSortBy         string
	dupErrorParser        validation.DuplicateErrorParser
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
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, r)
			}

			return err
		}

		resp, err = s.GetByID(ctx, r.ID)

		return err
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s store) GetByID(ctx context.Context, id string) (*Response, error) {
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

func (s store) Find(ctx context.Context, query ListRequest) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := mongoquery.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		mongoquery.GetSortQuery(cmp.Or(query.SortBy, s.defaultSortBy), query.Sort),
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

		resp, err = s.GetByID(ctx, r.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s store) Delete(ctx context.Context, id, userID string) (bool, error) {
	var deleted int64
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})

		return err
	})
	if err != nil || deleted == 0 {
		return false, err
	}

	cursor, err := s.dbReadCollection.Find(ctx, bson.M{"message": id}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return false, err
	}

	defer cursor.Close(ctx)
	ids := make([]string, 0, canopsis.DefaultBulkSize)
	for cursor.Next(ctx) {
		status := struct {
			ID string `bson:"_id"`
		}{}
		err = cursor.Decode(&status)
		if err != nil {
			return false, err
		}

		ids = append(ids, status.ID)
		if len(ids) == canopsis.DefaultBulkSize {
			_, err = s.dbReadCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
			if err != nil {
				return false, err
			}

			ids = ids[:0]
		}
	}

	if err = cursor.Err(); err != nil {
		return false, err
	}

	if len(ids) > 0 {
		_, err = s.dbReadCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (s store) GetActive(ctx context.Context, userID string) ([]Response, error) {
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
	}
	if userID != "" {
		pipeline = append(pipeline,
			bson.M{
				"$lookup": bson.M{
					"from":         mongo.BroadcastMessageReadCollection,
					"localField":   "_id",
					"foreignField": "message",
					"as":           "status",
					"pipeline": []bson.M{
						{"$match": bson.M{"user": userID}},
						{"$limit": 1},
					},
				},
			},
			bson.M{
				"$unwind": bson.M{"path": "$status", "preserveNullAndEmptyArrays": true},
			},
			bson.M{
				"$match": bson.M{
					"status": nil,
				},
			},
		)
	}

	pipeline = append(pipeline, bson.M{
		"$sort": bson.D{
			{Key: "start", Value: -1},
			{Key: "_id", Value: 1},
		},
	})
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

func (s store) Read(ctx context.Context, id, userID string) (bool, error) {
	now := datetime.NewCpsTime()
	var ok bool
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		ok = false
		err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		_, err = s.dbReadCollection.UpdateOne(ctx,
			bson.M{
				"user":    userID,
				"message": id,
			},
			bson.M{
				"$setOnInsert": bson.M{
					"_id":     utils.NewID(),
					"user":    userID,
					"message": id,
					"t":       now,
				},
			},
			options.UpdateOne().SetUpsert(true),
		)
		if err != nil {
			return err
		}

		ok = true

		return nil
	})

	return ok, err
}
