package entity

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(ctx context.Context, r ListRequestWithPagination) (*AggregationResult, error)
	Toggle(ctx context.Context, id string, enabled bool) (bool, SimplifiedEntity, error)
	GetContextGraph(ctx context.Context, id string) (*ContextGraphResponse, error)
}

type store struct {
	db                     mongo.DbClient
	mainCollection         mongo.DbCollection
	archivedCollection     mongo.DbCollection
	timezoneConfigProvider config.TimezoneConfigProvider
}

func NewStore(db mongo.DbClient, timezoneConfigProvider config.TimezoneConfigProvider) Store {
	return &store{
		db:                     db,
		mainCollection:         db.Collection(mongo.EntityMongoCollection),
		archivedCollection:     db.Collection(mongo.ArchivedEntitiesMongoCollection),
		timezoneConfigProvider: timezoneConfigProvider,
	}
}

func (s *store) Find(ctx context.Context, r ListRequestWithPagination) (*AggregationResult, error) {
	location := s.timezoneConfigProvider.Get().Location
	now := types.CpsTime{Time: time.Now().In(location)}

	pipeline, err := s.getQueryBuilder().CreateListAggregationPipeline(ctx, r, now)
	if err != nil {
		return nil, err
	}

	cursor, err := s.mainCollection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))

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

	s.fillConnectorType(&res)

	return &res, nil
}

func (s *store) Toggle(ctx context.Context, id string, enabled bool) (bool, SimplifiedEntity, error) {
	var isToggled bool
	var oldSimplifiedEntity SimplifiedEntity

	err := s.db.WithTransaction(ctx, func(ctx context.Context) error {
		isToggled = false
		oldSimplifiedEntity = SimplifiedEntity{}

		err := s.mainCollection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": id},
			bson.M{"$set": bson.M{"enabled": enabled}},
			options.
				FindOneAndUpdate().
				SetProjection(bson.M{"_id": 1, "enabled": 1, "type": 1}).
				SetReturnDocument(options.Before),
		).Decode(&oldSimplifiedEntity)
		if err != nil {
			return err
		}

		isToggled = oldSimplifiedEntity.Enabled != enabled
		return nil
	})

	return isToggled, oldSimplifiedEntity, err
}

func (s *store) GetContextGraph(ctx context.Context, id string) (*ContextGraphResponse, error) {
	res := ContextGraphResponse{}
	err := s.mainCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &res, nil
}

func (s *store) fillConnectorType(result *AggregationResult) {
	if result == nil {
		return
	}
	for i := range result.Data {
		result.Data[i].fillConnectorType()
	}
}

func (s *store) getQueryBuilder() *MongoQueryBuilder {
	return NewMongoQueryBuilder(s.db)
}
