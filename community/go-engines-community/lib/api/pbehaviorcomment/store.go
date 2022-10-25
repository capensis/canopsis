package pbehaviorcomment

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(ctx context.Context, pbehaviorID string, model *pbehavior.Comment) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbCollection mongo.DbCollection
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbCollection: dbClient.Collection(mongo.PbehaviorMongoCollection),
	}
}

func (s *store) Insert(ctx context.Context, pbehaviorID string, model *pbehavior.Comment) (bool, error) {
	model.ID = utils.NewID()
	model.Timestamp = &types.CpsTime{Time: time.Now()}
	res, err := s.dbCollection.UpdateOne(
		ctx,
		bson.M{"_id": pbehaviorID},
		bson.M{"$push": bson.M{
			"comments": model,
		}},
	)
	if err != nil {
		return false, err
	}

	return res.ModifiedCount > 0, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	res, err := s.dbCollection.UpdateOne(
		ctx,
		bson.M{"comments._id": id},
		bson.M{"$pull": bson.M{
			"comments": bson.M{"_id": id},
		}},
	)
	if err != nil {
		return false, err
	}

	return res.ModifiedCount > 0, nil
}
