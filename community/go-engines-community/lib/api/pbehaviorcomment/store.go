package pbehaviorcomment

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Store interface {
	Insert(pbehaviorID string, model *pbehavior.Comment) (bool, error)
	Delete(id string) (bool, error)
}

type store struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(pbehavior.PBehaviorCollectionName),
	}
}

func (s *store) Insert(pbehaviorID string, model *pbehavior.Comment) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (s *store) Delete(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
