package pbehaviorreason

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Store interface {
	Insert(model *Reason) error
	Find(query ListRequest) (*AggregationResult, error)
	GetOneBy(filter bson.M) (*Reason, error)
	Update(model *Reason) (bool, error)
	Delete(id string) (bool, error)
	IsLinkedToPbehavior(id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbClient:      dbClient,
		defaultSortBy: "created",
	}
}

type store struct {
	dbClient      mongo.DbClient
	defaultSortBy string
}

func (s *store) Find(query ListRequest) (*AggregationResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filter bson.M

	if query.Search != "" {
		searchRegexp := primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", query.Search),
			Options: "i",
		}

		filter = bson.M{
			"$or": []bson.M{
				{"name": searchRegexp},
			},
		}
	} else {
		filter = bson.M{}
	}

	sortBy := query.SortBy
	if sortBy == "" {
		sortBy = s.defaultSortBy
	}

	var project []bson.M
	if query.WithFlags {
		project = getDeletablePipeline()
	}

	collection := s.dbClient.Collection(pbehavior.ReasonCollectionName)
	pipeline := pagination.CreateAggregationPipeline(
		query.Query,
		[]bson.M{{"$match": filter}},
		common.GetSortQuery(sortBy, query.Sort),
		project,
	)
	cursor, err := collection.Aggregate(ctx, pipeline,
		options.Aggregate().SetCollation(&options.Collation{Locale: "en"}))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	cursor.Next(ctx)

	var result AggregationResult
	err = cursor.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) Insert(model *Reason) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	doc := transformModelToDoc(model)

	if model.ID == "" {
		model.ID = utils.NewID()
	}

	doc.ID = model.ID
	doc.Created = types.NewCpsTime(time.Now().Unix())

	_, err := s.dbClient.Collection(pbehavior.ReasonCollectionName).InsertOne(ctx, doc)

	if err != nil {
		return err
	}

	return nil
}

func (s *store) GetOneBy(filter bson.M) (*Reason, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var reason Reason

	err := s.dbClient.Collection(pbehavior.ReasonCollectionName).FindOne(ctx, filter).Decode(&reason)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &reason, nil
}

func (s *store) Update(model *Reason) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	doc := transformModelToDoc(model)
	result, err := s.dbClient.Collection(pbehavior.ReasonCollectionName).UpdateOne(
		ctx,
		bson.M{"_id": model.ID},
		bson.M{
			"$set": doc,
		},
	)

	if err != nil {
		return false, err
	}

	return result.MatchedCount > 0, nil
}

func (s *store) Delete(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	isLinkedToPbehavior, err := s.IsLinkedToPbehavior(id)
	if err != nil {
		return false, err
	}

	if isLinkedToPbehavior {
		return false, ErrLinkedReasonToPbehavior
	}

	isLinkedToAction, err := s.isLinkedToAction(id)
	if err != nil {
		return false, err
	}

	if isLinkedToAction {
		return false, ErrLinkedReasonToAction
	}

	deleted, err := s.dbClient.Collection(pbehavior.ReasonCollectionName).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

// IsLinked checks if there is pbehavior with linked reason.
func (s *store) IsLinkedToPbehavior(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res := s.dbClient.
		Collection(pbehavior.PBehaviorCollectionName).
		FindOne(ctx, bson.M{"reason": id})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *store) isLinkedToAction(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res := s.dbClient.
		Collection(mongo.ScenarioMongoCollection).
		FindOne(ctx, bson.M{
			"actions": bson.M{
				"$elemMatch": bson.M{
					"type":              types.ActionTypePbehavior,
					"parameters.reason": id,
				},
			}})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func transformModelToDoc(reason *Reason) *pbehavior.Reason {
	return &pbehavior.Reason{
		Name:        reason.Name,
		Description: reason.Description,
	}
}

func getDeletablePipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         pbehavior.PBehaviorCollectionName,
			"localField":   "_id",
			"foreignField": "reason",
			"as":           "pbhs",
		}},
		{"$lookup": bson.M{
			"from":         mongo.ScenarioMongoCollection,
			"localField":   "_id",
			"foreignField": "actions.parameters.reason",
			"as":           "actions",
		}},
		{"$addFields": bson.M{
			"deletable": bson.M{"$and": []bson.M{
				{"$eq": bson.A{bson.M{"$size": "$pbhs"}, 0}},
				{"$eq": bson.A{bson.M{"$size": "$actions"}, 0}},
			}},
		}},
		{"$project": bson.M{
			"pbhs":    0,
			"actions": 0,
		}},
	}
}
