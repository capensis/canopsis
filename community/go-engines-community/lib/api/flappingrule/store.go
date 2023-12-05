package flappingrule

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/priority"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(ctx context.Context, r CreateRequest) (*Response, error)
	GetById(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, r UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbClient       mongo.DbClient
	dbCollection   mongo.DbCollection
	authorProvider author.Provider

	defaultSearchByFields []string
}

func NewStore(
	dbClient mongo.DbClient,
	authorProvider author.Provider,
) Store {
	return &store{
		dbClient:       dbClient,
		dbCollection:   dbClient.Collection(mongo.FlappingRuleMongoCollection),
		authorProvider: authorProvider,

		defaultSearchByFields: []string{"_id", "author.name", "name", "description"},
	}
}

func (s *store) Insert(ctx context.Context, r CreateRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	rule := transformRequestToModel(r.EditRequest)
	if r.ID == "" {
		r.ID = utils.NewID()
	}

	rule.ID = r.ID
	rule.Created = now
	rule.Updated = now

	var resp *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		resp = nil

		_, err := s.dbCollection.InsertOne(ctx, rule)
		if err != nil {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.dbCollection, rule.ID, rule.Priority)
		if err != nil {
			return err
		}

		resp, err = s.GetById(ctx, rule.ID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *store) GetById(ctx context.Context, id string) (*Response, error) {
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
	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := "created"
	if query.SortBy != "" {
		sortBy = query.SortBy
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		common.GetSortQuery(sortBy, query.Sort),
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

func (s *store) Update(ctx context.Context, r UpdateRequest) (*Response, error) {
	model := transformRequestToModel(r.EditRequest)
	model.ID = r.ID
	model.Updated = datetime.NewCpsTime()
	var resp *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		resp = nil

		_, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": model.ID}, bson.M{"$set": model})
		if err != nil {
			return err
		}

		err = priority.UpdateFollowing(ctx, s.dbCollection, model.ID, model.Priority)
		if err != nil {
			return err
		}

		resp, err = s.GetById(ctx, model.ID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func transformRequestToModel(r EditRequest) flappingrule.Rule {
	return flappingrule.Rule{
		Name:        r.Name,
		Description: r.Description,
		FreqLimit:   r.FreqLimit,
		Duration:    r.Duration,
		AlarmPatternFields: r.AlarmPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInAlarmPattern(mongo.FlappingRuleMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.FlappingRuleMongoCollection),
		),
		EntityPatternFields: r.EntityPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInEntityPattern(mongo.FlappingRuleMongoCollection),
		),
		Priority: r.Priority,
		Author:   r.Author,
	}
}
