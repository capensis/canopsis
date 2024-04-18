package entitycomment

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const commentsLimit = 100

type Store interface {
	Find(context.Context, ListRequest) (*AggregationResult, error)
	Insert(ctx context.Context, r Request, userID, username string) (*Response, error)
	Update(ctx context.Context, r UpdateRequest, userID, username string) (*Response, error)
}

type store struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
	logger       zerolog.Logger
}

func NewStore(dbClient mongo.DbClient, logger zerolog.Logger) Store {
	return &store{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.EntityMongoCollection),
		logger:       logger,
	}
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": r.Entity}},
		{"$project": bson.M{"comments": 1}},
		{"$unwind": "$comments"},
		{"$replaceRoot": bson.M{"newRoot": "$comments"}},
	}
	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		bson.M{"$sort": bson.M{"t": -1}},
	))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	res := AggregationResult{}
	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *store) Insert(ctx context.Context, r Request, userID, username string) (*Response, error) {
	sr := s.dbCollection.FindOne(ctx, bson.M{"_id": r.Entity}, options.FindOne().SetProjection(bson.M{"type": 1}))
	if sr.Err() != nil {
		if errors.Is(sr.Err(), mongodriver.ErrNoDocuments) {
			return nil, common.NewValidationError("entity", "Entity not found.")
		}
		return nil, sr.Err()
	}
	var entity struct {
		Type string `bson:"type"`
	}
	err := sr.Decode(&entity)
	if err != nil || entity.Type != types.EntityTypeService && entity.Type != types.EntityTypeResource {
		return nil, common.NewValidationError("entity", "Invalid entity type.")
	}
	doc := types.EntityComment{
		ID:        utils.NewID(),
		Timestamp: datetime.NewCpsTime(),
		Author:    &types.Author{ID: userID, DisplayName: username},
		Message:   r.Message,
	}
	filter := bson.M{"_id": r.Entity}
	// set update with insert as first item in comments array
	update := []bson.M{
		{"$set": bson.M{
			"last_comment": doc,
			"comments": bson.M{
				"$ifNull": []any{
					bson.M{"$concatArrays": []interface{}{[]types.EntityComment{doc}, "$comments"}},
					[]types.EntityComment{doc},
				},
			},
		}},
		{"$set": bson.M{
			"comments": bson.M{
				"$slice": []any{
					"$comments", commentsLimit,
				},
			},
		}},
	}

	if _, err := s.dbCollection.UpdateOne(ctx, filter, update); err != nil {
		return nil, err
	}
	return &Response{
		ID:     doc.ID,
		Entity: r.Entity,
		Comment: Comment{
			Timestamp: doc.Timestamp,
			Author:    doc.Author,
			Message:   doc.Message,
		},
	}, nil

}

func (s *store) Update(ctx context.Context, r UpdateRequest, userID, username string) (*Response, error) {
	doc := types.EntityComment{
		ID:        r.ID,
		Timestamp: datetime.NewCpsTime(),
		Author:    &types.Author{ID: userID, DisplayName: username},
		Message:   r.Message,
	}
	filter := bson.M{
		"_id":            r.Entity,
		"comments.0._id": r.ID,
	}
	res, err := s.dbCollection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"comments.0": doc, "last_comment": doc}})
	if err != nil {
		return nil, err
	}
	if res.ModifiedCount == 0 {
		return nil, common.NewValidationError("_id", "Comment cannot be edited.")
	}
	return &Response{
		ID:     doc.ID,
		Entity: r.Entity,
		Comment: Comment{
			Timestamp: doc.Timestamp,
			Author:    doc.Author,
			Message:   doc.Message,
		},
	}, nil
}
