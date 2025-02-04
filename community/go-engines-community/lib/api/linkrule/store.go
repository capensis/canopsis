package linkrule

import (
	"cmp"
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apiexternaldata "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/externaldatatable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const defaultCategoriesLimit = 100

type Store interface {
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	GetByID(ctx context.Context, id string) (*Response, error)
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id, userID string) (bool, error)
	GetCategories(ctx context.Context, r CategoriesRequest) (*CategoryResponse, error)
}

type store struct {
	client                mongo.DbClient
	collection            mongo.DbCollection
	exdataTableCollection mongo.DbCollection
	authorProvider        author.Provider

	defaultSearchByFields []string
	defaultSortBy         string
}

func NewStore(dbClient mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		client:                dbClient,
		collection:            dbClient.Collection(mongo.LinkRuleMongoCollection),
		exdataTableCollection: dbClient.Collection(mongo.ExternalDataTableCollection),
		authorProvider:        authorProvider,

		defaultSearchByFields: []string{"_id", "author.name", "name"},
		defaultSortBy:         "created",
	}
}

func (s *store) Insert(ctx context.Context, request EditRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	model, err := s.transformRequestToModel(ctx, request)
	if err != nil {
		return nil, err
	}

	model.ID = utils.NewID()
	model.Created = now
	model.Updated = now

	var response *Response
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		_, err = s.collection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		response, err = s.GetByID(ctx, model.ID)

		return err
	})

	return response, err
}

func (s *store) GetByID(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
	pipeline = append(pipeline, apiexternaldata.GetRefParametersLookups()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		response := Response{}
		err = cursor.Decode(&response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *store) Find(ctx context.Context, request ListRequest) (*AggregationResult, error) {
	pipeline := s.authorProvider.Pipeline()
	filter := common.GetSearchQuery(request.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sort := common.GetSortQuery(cmp.Or(request.SortBy, s.defaultSortBy), request.Sort)
	project := apiexternaldata.GetRefParametersLookups()
	project = append(project, sort)
	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		pipeline,
		sort,
		project,
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

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) Update(ctx context.Context, request EditRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	model, err := s.transformRequestToModel(ctx, request)
	if err != nil {
		return nil, err
	}

	model.ID = request.ID
	model.Updated = now
	var response *Response
	err = s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		res, err := s.collection.UpdateOne(
			ctx,
			bson.M{"_id": request.ID},
			bson.M{"$set": model},
		)
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		response, err = s.GetByID(ctx, model.ID)

		return err
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id, userID string) (bool, error) {
	var deleted int64

	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		deleted = 0

		// required to get the author in action log listener.
		result, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || result.MatchedCount == 0 {
			return err
		}

		deleted, err = s.collection.DeleteOne(ctx, bson.M{"_id": id})
		return err
	})

	return deleted > 0, err
}

// GetCategories returns list of distinct categories
func (s *store) GetCategories(ctx context.Context, r CategoriesRequest) (*CategoryResponse, error) {
	pipeline := make([]bson.M, 0)
	if r.Type != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"type": r.Type}})
	}
	queryLimit := r.Limit
	if queryLimit == 0 {
		queryLimit = defaultCategoriesLimit
	}
	pipeline = append(pipeline,
		bson.M{"$unwind": "$links"},
	)
	if r.Search != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"links.category": bson.Regex{
				Pattern: fmt.Sprintf(".*%s.*", r.Search),
				Options: "i",
			},
		}})
	}
	pipeline = append(pipeline,
		bson.M{"$group": bson.M{
			"_id": "$links.category",
		}},
		bson.M{"$sort": bson.M{"_id": 1}},
		bson.M{"$limit": queryLimit},
		bson.M{"$group": bson.M{
			"_id": nil,
			"categories": bson.M{
				"$push": "$_id",
			},
		}},
		bson.M{"$project": bson.M{"_id": 0}},
	)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	resp := CategoryResponse{}
	if !cursor.Next(ctx) {
		return &resp, nil
	}

	if err = cursor.Decode(&resp); err != nil {
		return nil, err
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *store) transformRequestToModel(ctx context.Context, r EditRequest) (link.Rule, error) {
	externalData, err := apiexternaldata.TransformRefParameters(ctx, r.ExternalData, s.exdataTableCollection)
	if err != nil {
		return link.Rule{}, err
	}

	rule := link.Rule{
		Name:         r.Name,
		Type:         r.Type,
		Enabled:      *r.Enabled,
		Links:        r.Links,
		SourceCode:   r.SourceCode,
		ExternalData: externalData,
		Author:       r.Author,
		EntityPatternFields: r.EntityPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInEntityPattern(mongo.LinkRuleMongoCollection),
		),
	}
	if r.Type == link.TypeAlarm {
		rule.AlarmPatternFields = r.AlarmPatternFieldsRequest.ToModelWithoutFields(
			common.GetForbiddenFieldsInAlarmPattern(mongo.LinkRuleMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.LinkRuleMongoCollection),
		)
	}

	return rule, nil
}
