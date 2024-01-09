package token

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Token struct {
	ID       string           `bson:"_id"`
	User     string           `bson:"user"`
	Provider string           `bson:"provider,omitempty"`
	Created  datetime.CpsTime `bson:"created"`
	Accessed datetime.CpsTime `bson:"accessed"`

	Expired             *datetime.CpsTime          `bson:"expired,omitempty"`
	ExpiredByInactivity *datetime.CpsTime          `bson:"expired_by_inactivity,omitempty"`
	MaxInactiveInterval *datetime.DurationWithUnit `bson:"max_inactive_interval,omitempty"`
}

func NewMongoStore(client mongo.DbClient, logger zerolog.Logger) *MongoStore {
	return &MongoStore{
		collection: client.Collection(mongo.TokenMongoCollection),
		logger:     logger,
	}
}

type MongoStore struct {
	collection mongo.DbCollection
	logger     zerolog.Logger
}

func (s *MongoStore) Save(ctx context.Context, token Token) error {
	if token.ID == "" || token.User == "" || token.Created.IsZero() {
		return fmt.Errorf("invalid token: %v", token)
	}
	token.Accessed = token.Created
	if token.MaxInactiveInterval != nil && token.MaxInactiveInterval.Value > 0 {
		expiredByInactivity := token.MaxInactiveInterval.AddTo(token.Accessed)
		token.ExpiredByInactivity = &expiredByInactivity
	}
	_, err := s.collection.InsertOne(ctx, token)
	return err
}

func (s *MongoStore) Exists(ctx context.Context, id string) (bool, error) {
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *MongoStore) Access(ctx context.Context, id string) error {
	now := datetime.NewCpsTime()
	token := Token{}
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&token)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil
		}
		return err
	}

	update := bson.M{"accessed": now}
	if token.MaxInactiveInterval != nil && token.MaxInactiveInterval.Value > 0 {
		update["expired_by_inactivity"] = token.MaxInactiveInterval.AddTo(now)
	}

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (s *MongoStore) Count(ctx context.Context) (int64, error) {
	return s.collection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"expired": nil},
			{"expired": bson.M{"$gt": datetime.NewCpsTime()}},
		},
	})
}

func (s *MongoStore) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return deleted > 0, err
}

func (s *MongoStore) DeleteBy(ctx context.Context, user, provider string) error {
	_, err := s.collection.DeleteMany(ctx, bson.M{"user": user, "provider": provider})
	return err
}

func (s *MongoStore) DeleteByUserIDs(ctx context.Context, ids []string) error {
	_, err := s.collection.DeleteMany(ctx, bson.M{"user": bson.M{"$in": ids}})
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

	deleted, err = s.collection.DeleteMany(ctx, bson.M{
		"expired_by_inactivity": bson.M{
			"$ne": nil,
			"$lt": now,
		},
	})
	if err != nil {
		return err
	}

	if deleted > 0 {
		s.logger.Debug().Int64("deleted", deleted).Msg("deleted inactive tokens")
	}
	return nil
}
