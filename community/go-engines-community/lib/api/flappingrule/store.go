package flappingrule

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Insert(ctx context.Context, model *CreateRequest) (*RuleResponse, error)
	GetById(ctx context.Context, id string) (*RuleResponse, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, model *CreateRequest) (*RuleResponse, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.FlappingRuleMongoCollection),
	}
}

func (s *store) updateFollowingPriorities(ctx context.Context, id string, priority int) error {
	err := s.dbCollection.FindOne(ctx, bson.M{
		"_id":      bson.M{"$ne": id},
		"priority": priority,
	}).Err()
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil
		}
		return err
	}

	_, err = s.dbCollection.UpdateMany(
		ctx,
		bson.M{
			"_id":      bson.M{"$ne": id},
			"priority": bson.M{"$gte": priority},
		},
		bson.M{"$inc": bson.M{"priority": 1}},
	)

	return err
}

func (s *store) Insert(ctx context.Context, model *CreateRequest) (*RuleResponse, error) {
	if model.ID == "" {
		model.ID = utils.NewID()
	}
	now := types.NewCpsTime(time.Now().Unix())
	model.Created = &now
	model.Updated = &now

	_, err := s.dbCollection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	err = s.updateFollowingPriorities(ctx, model.ID, *model.Priority)
	if err != nil {
		return nil, err
	}

	return s.GetById(ctx, model.ID)
}

func (s *store) GetById(ctx context.Context, id string) (*RuleResponse, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorPipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		rr := &RuleResponse{}
		err := cursor.Decode(&rr)
		if err != nil {
			return nil, err
		}

		return rr, nil
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	var filter bson.M

	if query.Search != "" {
		searchRegexp := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", query.Search),
			Options: "i",
		}

		filter = bson.M{
			"$or": []bson.M{
				{"_id": searchRegexp},
				{"author": searchRegexp},
				{"description": searchRegexp},
			},
		}
	} else {
		filter = bson.M{}
	}
	pipeline := []bson.M{
		{"$match": filter},
	}

	sortBy := "created"
	if query.SortBy != "" {
		sortBy = query.SortBy
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		common.GetSortQuery(sortBy, query.Sort),
		s.authorPipeline(),
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var result AggregationResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}
	return &result, nil
}

func (s *store) Update(ctx context.Context, model *CreateRequest) (*RuleResponse, error) {
	newPriority := model.Priority
	updated := types.NewCpsTime(time.Now().Unix())
	model.Updated = &updated
	err := s.dbCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": model.ID},
		bson.M{"$set": model},
	).Decode(&model)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	if model.Priority != newPriority {
		err := s.updateFollowingPriorities(ctx, model.ID, *model.Priority)
		if err != nil {
			return nil, err
		}
	}

	rule, err := s.GetById(ctx, model.ID)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (s *store) authorPipeline() []bson.M {
	return []bson.M{
		// Author
		{"$lookup": bson.M{
			"from":         mongo.RightsMongoCollection,
			"localField":   "author",
			"foreignField": "_id",
			"as":           "author",
		}},
		{"$unwind": bson.M{"path": "$author", "preserveNullAndEmptyArrays": true}},
	}
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
