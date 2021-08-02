package token

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Store interface {
	Save(ctx context.Context, token string, expiredAt types.CpsTime) error
	Exists(ctx context.Context, token string) (bool, error)
	Count(ctx context.Context) (int64, error)
	Delete(ctx context.Context, token string) (bool, error)
	DeleteExpired(ctx context.Context, interval time.Duration)
}

func NewMongoStore(client mongo.DbClient, logger zerolog.Logger) Store {
	return &mongoStore{
		collection: client.Collection(mongo.TokenMongoCollection),
		logger:     logger,
	}
}

type mongoStore struct {
	collection mongo.DbCollection
	logger     zerolog.Logger
}

func (s *mongoStore) Save(ctx context.Context, token string, expiredAt types.CpsTime) error {
	_, err := s.collection.InsertOne(ctx, bson.M{
		"_id":        token,
		"expired_at": expiredAt,
	})
	return err
}

func (s *mongoStore) Exists(ctx context.Context, token string) (bool, error) {
	err := s.collection.FindOne(ctx, bson.M{"_id": token}).Err()
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *mongoStore) Count(ctx context.Context) (int64, error) {
	return s.collection.CountDocuments(ctx, bson.M{
		"expired_at": bson.M{"$gt": types.CpsTime{Time: time.Now()}},
	})
}

func (s *mongoStore) Delete(ctx context.Context, token string) (bool, error) {
	deleted, err := s.collection.DeleteOne(ctx, bson.M{"_id": token})
	return deleted > 0, err
}

func (s *mongoStore) DeleteExpired(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, err := s.collection.DeleteMany(ctx, bson.M{
				"expired_at": bson.M{"$lt": types.CpsTime{Time: time.Now()}},
			})
			if err != nil {
				s.logger.Err(err).Msg("cannot delete expired tokens")
			}
		}
	}
}
