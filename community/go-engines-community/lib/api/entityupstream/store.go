package entityupstream

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store interface {
	GetDownstreams(ctx context.Context, r DownstreamsRequest) (*AggregationResult, error)
	GetUpstream(ctx context.Context, id string) (res *Response, entityExists bool, err error)
}

type store struct {
	dbClient       mongo.DbClient
	dbCollection   mongo.DbCollection
	authorProvider author.Provider
}

func NewStore(
	dbClient mongo.DbClient,
	authorProvider author.Provider,
) Store {
	return &store{
		dbClient:       dbClient,
		dbCollection:   dbClient.Collection(mongo.EntityMongoCollection),
		authorProvider: authorProvider,
	}
}

func (s *store) GetDownstreams(ctx context.Context, r DownstreamsRequest) (*AggregationResult, error) {
	upstream := types.Entity{}
	err := s.dbCollection.
		FindOne(ctx, bson.M{
			"_id":     r.ID,
			"type":    bson.M{"$in": []string{types.EntityTypeResource, types.EntityTypeComponent}},
			"enabled": true,
		}).
		Decode(&upstream)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	now := datetime.NewCpsTime()
	match := bson.M{
		"upstream": upstream.ID,
		"type":     bson.M{"$in": []string{types.EntityTypeResource, types.EntityTypeComponent}},
		"enabled":  true,
	}

	pipeline := s.getQueryBuilder().CreateDownstreamAggregationPipeline(match, r.Query, r.SortRequest, now)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	result := &AggregationResult{}
	if cursor.Next(ctx) {
		err = cursor.Decode(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *store) GetUpstream(ctx context.Context, id string) (*Response, bool, error) {
	e := types.Entity{}
	err := s.dbCollection.FindOne(ctx,
		bson.M{
			"_id":     id,
			"type":    bson.M{"$in": []string{types.EntityTypeResource, types.EntityTypeComponent}},
			"enabled": true,
		},
		options.FindOne().SetProjection(bson.M{"upstream": 1}),
	).Decode(&e)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, false, nil
		}

		return nil, false, err
	}

	match := bson.M{
		"_id":     e.Upstream,
		"type":    bson.M{"$in": []string{types.EntityTypeResource, types.EntityTypeComponent}},
		"enabled": true,
	}
	result := &Response{}
	now := datetime.NewCpsTime()
	pipeline := s.getQueryBuilder().CreateUpstreamPipeline(match, now)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, false, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		err = cursor.Decode(result)
		if err != nil {
			return nil, false, err
		}

		return result, true, nil
	}

	if err = cursor.Err(); err != nil {
		return nil, false, err
	}

	return nil, true, nil
}

func (s *store) getQueryBuilder() *entity.MongoQueryBuilder {
	return entity.NewMongoQueryBuilder(s.dbClient, s.authorProvider)
}
