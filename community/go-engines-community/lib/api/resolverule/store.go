package resolverule

import (
	"cmp"
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/resolverule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	GetByID(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
}

type store struct {
	dbClient       mongo.DbClient
	dbCollection   mongo.DbCollection
	authorProvider author.Provider
	transformer    patternfields.Transformer

	defaultSearchByFields []string
	dupErrorParser        validation.DuplicateErrorParser
}

func NewStore(
	dbClient mongo.DbClient,
	authorProvider author.Provider,
	transformer patternfields.Transformer,
) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.ResolveRuleMongoCollection),
		authorProvider:        authorProvider,
		transformer:           transformer,
		defaultSearchByFields: []string{"_id", "author.name", "name", "description"},
		dupErrorParser:        validation.NewDuplicateErrorParser(),
	}
}

func (s *store) Insert(ctx context.Context, request CreateRequest) (*Response, error) {
	now := datetime.NewCpsTime()

	model := s.transformRequestToDocument(request.EditRequest)
	model.ID = cmp.Or(request.ID, utils.NewID())
	model.Created = now
	model.Updated = now

	var res *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = nil

		err := s.transformPatternRequestsToModel(ctx, request.EditRequest, &model)
		if err != nil {
			return err
		}

		_, err = s.dbCollection.InsertOne(ctx, model)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Response{})
			}

			return err
		}

		err = priority.UpdateFollowing(ctx, s.dbCollection, model.ID, request.Priority)
		if err != nil {
			return err
		}

		res, err = s.GetByID(ctx, model.ID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *store) GetByID(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		response := Response{}
		err := cursor.Decode(&response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	pipeline := s.authorProvider.Pipeline()
	filter := mongoquery.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		mongoquery.GetSortQuery(cmp.Or(query.SortBy, "created"), query.Sort),
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

func (s *store) Update(ctx context.Context, request UpdateRequest) (*Response, error) {
	model := s.transformRequestToDocument(request.EditRequest)
	model.Updated = datetime.NewCpsTime()
	var res *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = nil

		err := s.transformPatternRequestsToModel(ctx, request.EditRequest, &model)
		if err != nil {
			return err
		}

		_, err = s.dbCollection.UpdateOne(
			ctx,
			bson.M{"_id": request.ID},
			bson.M{"$set": model},
		)
		if err != nil {
			if mongodriver.IsDuplicateKeyError(err) {
				return s.dupErrorParser.Parse(err, Response{})
			}

			return err
		}

		err = priority.UpdateFollowing(ctx, s.dbCollection, request.ID, request.Priority)
		if err != nil {
			return err
		}

		res, err = s.GetByID(ctx, request.ID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	if id == resolverule.DefaultRule {
		return false, httperror.NewForbiddenError("The default rule cannot be deleted.")
	}

	var deleted int64

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		deleted, err = s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
}

func (s *store) transformRequestToDocument(r EditRequest) resolverule.Rule {
	return resolverule.Rule{
		Name:        r.Name,
		Description: r.Description,
		Duration:    r.Duration,
		Priority:    r.Priority,
		Author:      r.Author,
	}
}

func (s *store) transformPatternRequestsToModel(ctx context.Context, r EditRequest, model *resolverule.Rule) (err error) {
	model.AlarmPatternFields, model.EntityPatternFields, model.Aliases, err = s.transformer.TransformAlarmAndEntityRequest(ctx, r.AlarmRequest, r.EntityRequest, r, s.dbCollection.Name())

	return err
}
