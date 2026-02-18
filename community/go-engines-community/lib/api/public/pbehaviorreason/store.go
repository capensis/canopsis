package pbehaviorreason

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/dbvalidation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store crud.Store[ListRequest, EditRequest, EditRequest, Response]

func NewStore(dbClient mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.PbehaviorReasonMongoCollection),
		dbPbhCollection:       dbClient.Collection(mongo.PbehaviorMongoCollection),
		dbScenarioCollection:  dbClient.Collection(mongo.ScenarioCollection),
		authorProvider:        authorProvider,
		defaultSearchByFields: []string{"_id", "name", "description"},
		dupErrorParser:        validation.NewDuplicateErrorParser(),
	}
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	dbPbhCollection       mongo.DbCollection
	dbScenarioCollection  mongo.DbCollection
	authorProvider        author.Provider
	defaultSearchByFields []string
	dupErrorParser        validation.DuplicateErrorParser
}

func (s *store) Find(ctx context.Context, r ListRequest) (pagination.ListResponse[Response], error) {
	beforeLimit := []bson.M{
		{"$match": bson.M{"hidden": bson.M{"$in": bson.A{false, nil}}}},
	}
	filter := mongoquery.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		beforeLimit = append(beforeLimit, bson.M{"$match": filter})
	}

	afterLimit := s.authorProvider.Pipeline()
	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		mongoquery.GetSortQuery(r.SortBy, r.Sort),
		afterLimit,
	)
	opts := options.Aggregate().SetCollation(&options.Collation{Locale: "en"})
	var res pagination.ListResponse[Response]
	ar, err := db.AggregateOne[pagination.Data[Response]](ctx, s.dbCollection, pipeline, opts)
	if err != nil {
		return res, err
	}

	res = pagination.NewResponse(r.Query, ar)

	return res, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest) (Response, error) {
	doc := s.transformRequestToDocument(r)
	doc.ID = utils.NewID()
	doc.Created = datetime.NewCpsTime()
	doc.Updated = datetime.NewCpsTime()
	var res Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = Response{}

		_, err := s.dbCollection.InsertOne(ctx, doc)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, doc)
			}

			return err
		}

		res, err = s.GetByID(ctx, doc.ID)

		return err
	})

	return res, err
}

func (s *store) GetByID(ctx context.Context, id string) (Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	return db.AggregateOne[Response](ctx, s.dbCollection, pipeline)
}

func (s *store) Update(ctx context.Context, id string, r EditRequest) (Response, error) {
	doc := s.transformRequestToDocument(r)
	doc.ID = id
	doc.Updated = datetime.NewCpsTime()
	var res Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = Response{}

		uptRes, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": doc.ID}, bson.M{"$set": doc})
		if err != nil || uptRes.MatchedCount == 0 {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, doc)
			}

			return err
		}

		res, err = s.GetByID(ctx, doc.ID)

		return err
	})

	return res, err
}

func (s *store) Delete(ctx context.Context, id, userID string) error {
	return s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		err := s.validateDeleteRequest(ctx, id)
		if err != nil {
			return err
		}

		// required to get the author in action log listener.
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil {
			return err
		}

		if res.MatchedCount == 0 {
			return httperror.ErrNotFound
		}

		deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			return err
		}

		if deleted == 0 {
			return httperror.ErrNotFound
		}

		return nil
	})
}

func (s *store) validateDeleteRequest(ctx context.Context, id string) error {
	err := dbvalidation.ValidateLinkedReference(ctx, s.dbPbhCollection, bson.M{
		"reason": id,
	}, "reason", "a pbehavior")
	if err != nil {
		return err
	}

	return dbvalidation.ValidateLinkedReference(ctx, s.dbScenarioCollection, bson.M{
		"actions": bson.M{
			"$elemMatch": bson.M{
				"type":              types.ActionTypePbehavior,
				"parameters.reason": id,
			},
		},
	}, "reason", "an action scenario")
}

func (s *store) transformRequestToDocument(r EditRequest) pbehavior.Reason {
	return pbehavior.Reason{
		Name:        r.Name,
		Description: r.Description,
		Author:      r.Author,
	}
}
