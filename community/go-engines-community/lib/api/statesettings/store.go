package statesettings

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(ctx context.Context, query pagination.Query) (AggregationResult, error)
	Update(ctx context.Context, request StateSettingRequest) (*StateSetting, error)
}

type store struct {
	dbCollection mongo.DbCollection
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		dbCollection: dbClient.Collection(mongo.StateSettingsMongoCollection),
	}
}

func (s *store) Find(ctx context.Context, query pagination.Query) (AggregationResult, error) {
	result := AggregationResult{
		Data:       make([]StateSetting, 0),
		TotalCount: 0,
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query,
		[]bson.M{},
		bson.M{},
	))

	if err != nil {
		return result, err
	}

	if cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			return result, err
		}
	}

	err = cursor.Close(ctx)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *store) Update(ctx context.Context, r StateSettingRequest) (*StateSetting, error) {
	res := s.dbCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": r.ID},
		bson.M{"$set": r},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if err := res.Err(); err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	stateSetting := &StateSetting{}
	err := res.Decode(stateSetting)
	if err != nil {
		return nil, err
	}

	return stateSetting, nil
}
