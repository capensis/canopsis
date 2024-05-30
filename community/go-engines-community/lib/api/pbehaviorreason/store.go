package pbehaviorreason

import (
	"cmp"
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Insert(ctx context.Context, model CreateRequest) (*Response, error)
	Find(ctx context.Context, query ListRequest) (*AggregationResult, error)
	GetById(ctx context.Context, id string) (*Response, error)
	Update(ctx context.Context, model UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id, userId string) (bool, error)
	IsLinkedToPbehavior(ctx context.Context, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.PbehaviorReasonMongoCollection),
		authorProvider:        authorProvider,
		defaultSearchByFields: []string{"_id", "name", "description"},
		defaultSortBy:         "created",
	}
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	authorProvider        author.Provider
	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, query ListRequest) (*AggregationResult, error) {
	match := bson.M{}
	if !query.WithHidden {
		match["hidden"] = bson.M{"$in": bson.A{false, nil}}
	}

	pipeline := make([]bson.M, 0)
	if len(match) > 0 {
		pipeline = append(pipeline, bson.M{"$match": match})
	}

	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	var project []bson.M
	if query.WithFlags {
		project = getDeletablePipeline()
	}

	cursor, err := s.dbCollection.Aggregate(
		ctx,
		pagination.CreateAggregationPipeline(
			query.Query,
			pipeline,
			common.GetSortQuery(cmp.Or(query.SortBy, s.defaultSortBy), query.Sort),
			project,
		),
		options.Aggregate().SetCollation(&options.Collation{Locale: "en"}),
	)

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

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Response, error) {
	doc := transformModelToDoc(r.EditRequest)

	doc.ID = cmp.Or(r.ID, utils.NewID())
	doc.Created = datetime.NewCpsTime()
	doc.Updated = datetime.NewCpsTime()

	var res *Response

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = nil

		_, err := s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}

		res, err = s.GetById(ctx, doc.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *store) GetById(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if cursor.Next(ctx) {
		var res Response

		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, nil
}

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Response, error) {
	doc := transformModelToDoc(r.EditRequest)

	doc.ID = r.ID
	doc.Updated = datetime.NewCpsTime()

	var res *Response

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = nil

		result, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": doc.ID}, bson.M{"$set": doc})
		if err != nil || result.MatchedCount == 0 {
			return err
		}

		res, err = s.GetById(ctx, doc.ID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *store) Delete(ctx context.Context, id, userId string) (bool, error) {
	isLinkedToPbehavior, err := s.IsLinkedToPbehavior(ctx, id)
	if err != nil {
		return false, err
	}

	if isLinkedToPbehavior {
		return false, ErrLinkedReasonToPbehavior
	}

	isLinkedToAction, err := s.isLinkedToAction(ctx, id)
	if err != nil {
		return false, err
	}

	if isLinkedToAction {
		return false, ErrLinkedReasonToAction
	}

	var deleted int64

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userId}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
}

// IsLinkedToPbehavior checks if there is pbehavior with linked reason.
func (s *store) IsLinkedToPbehavior(ctx context.Context, id string) (bool, error) {
	res := s.dbClient.
		Collection(mongo.PbehaviorMongoCollection).
		FindOne(ctx, bson.M{"reason": id})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *store) isLinkedToAction(ctx context.Context, id string) (bool, error) {
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
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func transformModelToDoc(r EditRequest) *pbehavior.Reason {
	return &pbehavior.Reason{
		Name:        r.Name,
		Description: r.Description,
		Hidden:      r.Hidden,
		Author:      r.Author,
	}
}

func getDeletablePipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from": mongo.PbehaviorMongoCollection,
			"let":  bson.M{"id": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$$id", "$reason"}}}},
				{"$limit": 1},
				{"$project": bson.M{"_id": 1}},
			},
			"as": "pbhs",
		}},
		{"$lookup": bson.M{
			"from": mongo.ScenarioMongoCollection,
			"let":  bson.M{"id": "$_id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$expr": bson.M{"$in": bson.A{"$$id", "$actions.parameters.reason"}}}},
				{"$limit": 1},
				{"$project": bson.M{"_id": 1}},
			},
			"as": "actions",
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
