package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Save(ctx context.Context, token Token) error
	Exists(ctx context.Context, id string) (bool, error)
	Access(ctx context.Context, id string) error
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
	Accessed types.CpsTime `bson:"accessed"`

	Expired             *types.CpsTime         `bson:"expired"`
	ExpiredByInactivity *types.CpsTime         `bson:"expired_by_inactivity"`
	MaxInactiveInterval types.DurationWithUnit `bson:"max_inactive_interval"`
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
	if token.ID == "" || token.User == "" || token.Created.IsZero() {
		return fmt.Errorf("invalid token: %v", token)
	}
	token.Accessed = token.Created
	if token.MaxInactiveInterval.Value > 0 {
		expiredByInactivity := token.MaxInactiveInterval.AddTo(token.Accessed)
		token.ExpiredByInactivity = &expiredByInactivity
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

func (s *mongoStore) Access(ctx context.Context, id string) error {
	now := types.NewCpsTime()
	token := Token{}
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&token)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil
		}
		return err
	}

	update := bson.M{"accessed": now}
	if token.MaxInactiveInterval.Value > 0 {
		update["expired_by_inactivity"] = token.MaxInactiveInterval.AddTo(now)
	}

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (s *mongoStore) Count(ctx context.Context) (int64, error) {
	return s.collection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"expired": nil},
			{"expired": bson.M{"$gt": types.NewCpsTime()}},
		},
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
			err := s.deleteExpiredFromStorage(ctx)
			if err != nil {
				s.logger.Err(err).Msg("cannot delete expired tokens")
			}
		}
	}
}

func (s *mongoStore) deleteExpiredFromStorage(ctx context.Context) error {
	now := types.NewCpsTime()
	deleted, err := s.collection.DeleteMany(ctx, bson.M{
		"expired": bson.M{
			"$ne": nil,
			"$lt": now,
		},
	})
	if err != nil {
		return err
	}
	s.logger.Debug().Int64("deleted", deleted).Msg("deleted expired tokens")

	deleted, err = s.collection.DeleteMany(ctx, bson.M{
		"expired_by_inactivity": bson.M{
			"$ne": nil,
			"$lt": now,
		},
	})
	if err != nil {
		return err
	}

	s.logger.Debug().Int64("deleted", deleted).Msg("deleted inactive tokens")
	return nil
}
