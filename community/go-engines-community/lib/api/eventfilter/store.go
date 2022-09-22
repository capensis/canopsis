package eventfilter

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Insert(ctx context.Context, request CreateRequest) (*Response, error)
	GetById(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, request UpdateRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type AggregationResult struct {
	Data       []*Response `bson:"data" json:"data"`
	TotalCount int64       `bson:"total_count" json:"total_count"`
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

func NewStore(
	dbClient mongo.DbClient,
) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.EventFilterRulesMongoCollection),
		defaultSearchByFields: []string{"_id", "author", "description", "type"},
		defaultSortBy:         "created",
	}
}

func (s *store) transformRequestToDocument(r EditRequest) eventfilter.Rule {
	exdates := make([]types.Exdate, len(r.Exdates))
	for i := range r.Exdates {
		exdates[i].Begin = r.Exdates[i].Begin
		exdates[i].End = r.Exdates[i].End
	}

	return eventfilter.Rule{
		Author:              r.Author,
		Description:         r.Description,
		Type:                r.Type,
		Priority:            r.Priority,
		Enabled:             r.Enabled,
		Config:              r.Config,
		ExternalData:        r.ExternalData,
		EventPattern:        r.EventPattern,
		EntityPatternFields: r.EntityPatternFieldsRequest.ToModel(),
		RRule:               r.RRule,
		Start:               r.Start,
		Stop:                r.Stop,
		ResolvedStart:       r.Start,
		ResolvedStop:        r.Stop,
		Exdates:             exdates,
		Exceptions:          r.Exceptions,
	}
}

func (s *store) Insert(ctx context.Context, request CreateRequest) (*Response, error) {
	model := s.transformRequestToDocument(request.EditRequest)
	model.ID = request.ID
	if model.ID == "" {
		model.ID = utils.NewID()
	}

	now := types.NewCpsTime(time.Now().Unix())
	model.Created = &now
	model.Updated = &now

	var response *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.dbCollection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		response, err = s.GetById(ctx, model.ID)
		return err
	})

	return response, err
}

func (s *store) GetById(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": id},
		},
		{
			"$lookup": bson.M{
				"from":         mongo.PbehaviorExceptionMongoCollection,
				"localField":   "exceptions",
				"foreignField": "_id",
				"as":           "exceptions",
			},
		},
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var rule Response
		err = cursor.Decode(&rule)
		if err != nil {
			return nil, err
		}

		return &rule, nil
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := s.defaultSortBy
	if query.SortBy != "" {
		sortBy = query.SortBy
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		common.GetSortQuery(sortBy, query.Sort),
		[]bson.M{
			{
				"$lookup": bson.M{
					"from":         mongo.PbehaviorExceptionMongoCollection,
					"localField":   "exceptions",
					"foreignField": "_id",
					"as":           "exceptions",
				},
			},
		},
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
	updated := types.NewCpsTime(time.Now().Unix())
	model := s.transformRequestToDocument(request.EditRequest)
	model.ID = request.ID
	model.Created = nil
	model.Updated = &updated

	update := bson.M{"$set": model}
	if request.CorporateEntityPattern != "" || len(request.EntityPattern) > 0 || len(request.EventPattern) > 0 {
		update["$unset"] = bson.M{"old_patterns": 1}
	}

	var response *Response
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err := s.dbCollection.UpdateOne(
			ctx,
			bson.M{"_id": model.ID},
			update,
		)
		if err != nil {
			return err
		}

		response, err = s.GetById(ctx, model.ID)
		return err
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return deleted > 0, nil
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
