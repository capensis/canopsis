package token

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Store interface {
	Save(ctx context.Context, token Token) error
	Exists(ctx context.Context, id string) (bool, error)
	Count(ctx context.Context) (int64, error)
	Delete(ctx context.Context, id string) (bool, error)
	DeleteExpired(ctx context.Context, interval time.Duration)
	DeleteBy(ctx context.Context, user, provider string) error
}

type Token struct {
	ID       string        `bson:"_id"`
	User     string        `bson:"user"`
	Provider string        `bson:"provider,omitempty"`
	Created  types.CpsTime `bson:"created"`
	Expired  types.CpsTime `bson:"expired"`
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

func (s *mongoStore) Save(ctx context.Context, token Token) error {
	if token.ID == "" || token.User == "" || token.Created.IsZero() || token.Expired.IsZero() {
		return fmt.Errorf("invalid token: %v", token)
	}
	_, err := s.collection.InsertOne(ctx, token)
	return err
}

func (s *mongoStore) Exists(ctx context.Context, id string) (bool, error) {
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Err()
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
		"expired": bson.M{"$gt": types.CpsTime{Time: time.Now()}},
	})
}

func (s *mongoStore) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return deleted > 0, err
}

func (s *mongoStore) DeleteBy(ctx context.Context, user, provider string) error {
	_, err := s.collection.DeleteMany(ctx, bson.M{"user": user, "provider": provider})
	return err
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
				"expired": bson.M{"$lt": types.CpsTime{Time: time.Now()}},
			})
			if err != nil {
				s.logger.Err(err).Msg("cannot delete expired tokens")
			}
		}
	}
}
