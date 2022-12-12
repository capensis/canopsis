package pbehaviorcomment

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Insert(ctx context.Context, pbehaviorID string, model *pbehavior.Comment) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.PbehaviorMongoCollection),
	}
}

func (s *store) Insert(ctx context.Context, pbehaviorID string, model *pbehavior.Comment) (bool, error) {
	model.ID = utils.NewID()
	model.Timestamp = &types.CpsTime{Time: time.Now()}
	updated := false
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updated = false

		err := s.dbCollection.FindOne(ctx, bson.M{
			"_id":    pbehaviorID,
			"origin": bson.M{"$ne": nil},
		}).Err()
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}
		if err == nil {
			return common.NewValidationError("_id", errors.New("Cannot update a pbehavior with origin."))
		}

		res, err := s.dbCollection.UpdateOne(
			ctx,
			bson.M{"_id": pbehaviorID},
			bson.M{"$push": bson.M{
				"comments": model,
			}},
		)
		if err != nil {
			return err
		}

		updated = res.ModifiedCount > 0
		return nil
	})

	return updated, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	updated := false
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updated = false

		err := s.dbCollection.FindOne(ctx, bson.M{
			"comments._id": id,
			"origin":       bson.M{"$ne": nil},
		}).Err()
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}
		if err == nil {
			return common.NewValidationError("_id", errors.New("Cannot update a pbehavior with origin."))
		}

		res, err := s.dbCollection.UpdateOne(
			ctx,
			bson.M{"comments._id": id},
			bson.M{"$pull": bson.M{
				"comments": bson.M{"_id": id},
			}},
		)
		if err != nil {
			return err
		}

		updated = res.ModifiedCount > 0
		return nil
	})

	return updated, err
}
