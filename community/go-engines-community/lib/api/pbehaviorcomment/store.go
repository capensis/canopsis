package pbehaviorcomment

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(ctx context.Context, r Request) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
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

func (s *store) Insert(ctx context.Context, r Request) (*Response, error) {
	doc := pbehavior.Comment{
		ID:        utils.NewID(),
		Author:    r.Author,
		Timestamp: &types.CpsTime{Time: time.Now()},
		Message:   r.Message,
	}
	var response *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		res, err := s.dbCollection.UpdateOne(
			ctx,
			bson.M{"_id": r.Pbehavior},
			bson.M{"$push": bson.M{
				"comments": doc,
			}},
		)
		if err != nil || res.ModifiedCount == 0 {
			return err
		}

		pipeline := []bson.M{
			{"$match": bson.M{"_id": r.Pbehavior}},
			{"$unwind": "$comments"},
			{"$match": bson.M{"comments._id": doc.ID}},
			{"$replaceRoot": bson.M{"newRoot": "$comments"}},
		}
		pipeline = append(pipeline, author.Pipeline()...)
		cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)
		if cursor.Next(ctx) {
			response = &Response{}
			err = cursor.Decode(response)
			return err
		}
		return nil
	})

	return response, err
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
