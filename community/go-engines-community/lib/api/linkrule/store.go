package linkrule

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	GetById(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	client     mongo.DbClient
	collection mongo.DbCollection

	defaultSearchByFields []string
	defaultSortBy         string
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		client:     dbClient,
		collection: dbClient.Collection(mongo.LinkRuleMongoCollection),

		defaultSearchByFields: []string{"_id", "author.name", "name"},
		defaultSortBy:         "created",
	}
}

func (s *store) Insert(ctx context.Context, request EditRequest) (*Response, error) {
	now := types.NewCpsTime()
	model := transformRequestToModel(request)
	model.ID = utils.NewID()
	model.Created = now
	model.Updated = now

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.collection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		response, err = s.GetById(ctx, model.ID)
		return err
	})

	return response, err
}

func (s *store) GetById(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, author.Pipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
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

func (s *store) Find(ctx context.Context, request ListRequest) (*AggregationResult, error) {
	pipeline := author.Pipeline()
	filter := common.GetSearchQuery(request.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := s.defaultSortBy
	if request.SortBy != "" {
		sortBy = request.SortBy
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		pipeline,
		common.GetSortQuery(sortBy, request.Sort),
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

func (s *store) Update(ctx context.Context, request EditRequest) (*Response, error) {
	now := types.NewCpsTime()
	model := transformRequestToModel(request)
	model.ID = request.ID
	model.Updated = now
	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		res, err := s.collection.UpdateOne(
			ctx,
			bson.M{"_id": request.ID},
			bson.M{"$set": model},
		)
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		response, err = s.GetById(ctx, model.ID)
		return err
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return deleted > 0, err
}

func transformRequestToModel(r EditRequest) link.Rule {
	return link.Rule{
		Name:         r.Name,
		Type:         r.Type,
		Enabled:      *r.Enabled,
		Links:        r.Links,
		SourceCode:   r.SourceCode,
		ExternalData: r.ExternalData,
		Author:       r.Author,
		EntityPatternFields: r.EntityPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInEntityPattern(mongo.LinkRuleMongoCollection),
		),
		AlarmPatternFields: r.AlarmPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInAlarmPattern(mongo.LinkRuleMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.LinkRuleMongoCollection),
		),
	}
}
