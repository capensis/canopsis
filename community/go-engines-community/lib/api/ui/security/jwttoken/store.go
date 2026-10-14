package jwttoken

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/security/token"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

func NewStore(client mongo.DbClient, logger zerolog.Logger) *Store {
	return &Store{
		collection: client.Collection(mongo.TokenMongoCollection),
		logger:     logger,
	}
}

type Store struct {
	collection mongo.DbCollection
	logger     zerolog.Logger
}

func (s *Store) Save(ctx context.Context, token token.Token) error {
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

func (s *Store) Exists(ctx context.Context, id string) (bool, error) {
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *Store) Access(ctx context.Context, id string) error {
	now := datetime.NewCpsTime()
	token := token.Token{}
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

func (s *Store) Count(ctx context.Context) (int64, error) {
	return s.collection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"expired": nil},
			{"expired": bson.M{"$gt": datetime.NewCpsTime()}},
		},
	})
}

func (s *Store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return deleted > 0, err
}

func (s *Store) DeleteBy(ctx context.Context, user, provider string) error {
	_, err := s.collection.DeleteMany(ctx, bson.M{"user": user, "provider": provider})
	return err
}

func (s *Store) DeleteByUserIDs(ctx context.Context, ids []string) error {
	_, err := s.collection.DeleteMany(ctx, bson.M{"user": bson.M{"$in": ids}})
	return err
}
