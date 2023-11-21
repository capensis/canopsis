package sharetoken

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func NewMongoStore(client mongo.DbClient, logger zerolog.Logger) *MongoStore {
	return &MongoStore{
		collection: client.Collection(mongo.ShareTokenMongoCollection),
		logger:     logger,
	}
}

type MongoStore struct {
	collection mongo.DbCollection
	logger     zerolog.Logger
}

func (s *MongoStore) Exists(ctx context.Context, v string) (bool, error) {
	err := s.collection.FindOne(ctx, bson.M{"value": v}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *MongoStore) Access(ctx context.Context, v string) error {
	now := datetime.NewCpsTime()
	_, err := s.collection.UpdateOne(ctx, bson.M{"value": v}, bson.M{"$set": bson.M{"accessed": now}})
	return err
}

func (s *MongoStore) DeleteExpired(ctx context.Context) error {
	now := datetime.NewCpsTime()
	deleted, err := s.collection.DeleteMany(ctx, bson.M{
		"expired": bson.M{
			"$ne": nil,
			"$lt": now,
		},
	})
	if err != nil {
		return err
	}
	if deleted > 0 {
		s.logger.Debug().Int64("deleted", deleted).Msg("deleted expired tokens")
	}
	return nil
}
