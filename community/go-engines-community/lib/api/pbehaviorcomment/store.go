package pbehaviorcomment

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Insert(ctx context.Context, r Request) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbClient       mongo.DbClient
	dbCollection   mongo.DbCollection
	authorProvider author.Provider
}

func NewStore(dbClient mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		dbClient:       dbClient,
		dbCollection:   dbClient.Collection(mongo.PbehaviorMongoCollection),
		authorProvider: authorProvider,
	}
}

func (s *store) Insert(ctx context.Context, r Request) (*Response, error) {
	now := datetime.NewCpsTime()

	doc := pbehavior.Comment{
		ID:        utils.NewID(),
		Author:    r.Author,
		Timestamp: &now,
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
		pipeline = append(pipeline, s.authorProvider.Pipeline()...)
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
	updated := false
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updated = false

		pbh := pbehavior.PBehavior{}
		err := s.dbCollection.FindOne(ctx, bson.M{"comments._id": id}).Decode(&pbh)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		if pbh.Origin != "" && len(pbh.Comments) > 0 && pbh.Comments[0].ID == id {
			return common.NewValidationError("_id", "Cannot remove main comment.")
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
